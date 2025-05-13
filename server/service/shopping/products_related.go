package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type GetRelatedList struct{}

type tempRelatedProductInfo struct {
	ProductID           string  `json:"product_id"`             // 関連商品のID
	ProductCode         string  `json:"product_code,omitempty"` // 関連商品のコード
	ProductName         string  `json:"product_name"`           // 関連商品の名称 (省略後)
	PriceRangeFormatted string  `json:"price_range_formatted"`  // ★価格帯文字列 (例: "2,990～3,990円")
	IsOnSale            bool    `json:"is_on_sale"`
	AverageRating       float64 `json:"average_rating"` // 平均評価
	ReviewCount         int     `json:"review_count"`
	ThumbnailImageURL   *string `json:"thumbnail_image_url"`
}

func (g *GetRelatedList) GetRelatedListInfo(productCode string, limit int) ([]dto.RelatedProductInfo, error) {

	var categoryId string
	err := global.GVA_DB.Table("products").
		Where("product_code = ?", productCode).
		Select("category_id").
		Scan(&categoryId).Error
	if err != nil {
		return nil, err
	}

	var temp []tempRelatedProductInfo

	db := global.GVA_DB.Raw(`
		WITH price_priority AS (
			SELECT
				product_skus.product_id,
			CASE 
				WHEN MIN(CASE WHEN prices.price_type_id = 2 AND NOW() BETWEEN prices.start_date AND prices.end_date THEN prices.price END) IS NOT NULL
					THEN MIN(CASE WHEN prices.price_type_id = 2 AND NOW() BETWEEN prices.start_date AND prices.end_date THEN prices.price END)
					ELSE MIN(CASE WHEN prices.price_type_id = 1 THEN prices.price END)
				END AS min_price,
				CASE 
				WHEN MAX(CASE WHEN prices.price_type_id = 2 AND NOW() BETWEEN prices.start_date AND prices.end_date THEN prices.price END) IS NOT NULL
					THEN MAX(CASE WHEN prices.price_type_id = 2 AND NOW() BETWEEN prices.start_date AND prices.end_date THEN prices.price END)
					ELSE MAX(CASE WHEN prices.price_type_id = 1 THEN prices.price END)
				END AS max_price,
				CASE
					WHEN COUNT(CASE WHEN prices.price_type_id = 2 AND NOW() BETWEEN prices.start_date AND prices.end_date THEN 1 END) > 0
					THEN TRUE
					ELSE FALSE
				END AS is_on_sale
			FROM product_skus
			LEFT JOIN prices ON product_skus.id = prices.sku_id
			WHERE product_skus.product_id IN (SELECT id FROM products WHERE category_id = ?)
			GROUP BY product_skus.product_id
		),
		product_info AS (
			SELECT
				products.id AS product_id,
				products.product_code,
				CASE
					WHEN CHAR_LENGTH(products.name) > 10 THEN CONCAT(LEFT(products.name, 10), '...')
					ELSE products.name
				END AS product_name,
				(
					SELECT sku_images.thumbnail_url
					FROM sku_images
					WHERE sku_images.sku_id = product_skus.id
					LIMIT 1
				) AS thumbnail_image_url,
				review_summaries.average_rating,
				COUNT(DISTINCT CASE WHEN product_reviews.status = 'approved' THEN product_reviews.id END) AS review_count
			FROM products
			LEFT JOIN product_skus ON products.id = product_skus.product_id
			LEFT JOIN product_reviews ON products.id = product_reviews.product_id
			LEFT JOIN review_summaries ON products.id = review_summaries.product_id
			WHERE products.category_id = ?
			GROUP BY products.id, product_skus.id
		)
		SELECT
			product_info.product_id,
			product_info.product_code,
			product_info.product_name,
			CASE
				WHEN price_priority.min_price IS NULL THEN NULL
				WHEN price_priority.min_price = price_priority.max_price THEN CONCAT(FORMAT(price_priority.min_price, 0), '円')
				ELSE CONCAT(FORMAT(price_priority.min_price, 0), '～', FORMAT(price_priority.max_price, 0), '円')
			END AS price_range_formatted,
			price_priority.is_on_sale,
			product_info.thumbnail_image_url,
			product_info.average_rating,
			product_info.review_count
		FROM product_info
		JOIN price_priority ON product_info.product_id = price_priority.product_id
		ORDER BY 
		price_priority.is_on_sale DESC
	`, categoryId, categoryId).Scan(&temp)

	if db.Error != nil {
		return nil, db.Error
	}
	//查询结果为空的话，就返回一个nil
	if len(temp) == 0 {
		return nil, nil
	}
	// 排除当前商品
	var filtered []tempRelatedProductInfo
	for _, v := range temp {
		if v.ProductCode != productCode { // 过滤掉当前商品
			filtered = append(filtered, v)
		}
	}

	// 如果没有符合条件的相关商品，nil
	if len(filtered) == 0 {
		return nil, nil
	}

	//使用商品code去重
	uniqueRelated := make(map[string]tempRelatedProductInfo)
	for _, v1 := range filtered {
		if _, exists := uniqueRelated[v1.ProductCode]; !exists {
			uniqueRelated[v1.ProductCode] = v1
		}
	}

	var output []dto.RelatedProductInfo
	for _, v := range uniqueRelated {
		output = append(output, dto.RelatedProductInfo{
			ProductID:           v.ProductID,
			ProductCode:         v.ProductCode,
			ProductName:         v.ProductName,
			PriceRangeFormatted: v.PriceRangeFormatted,
			IsOnSale:            v.IsOnSale,
			ThumbnailImageURL:   v.ThumbnailImageURL,
			ReviewSummary: &dto.ReviewSummaryInfo{
				AverageRating: v.AverageRating,
				ReviewCount:   v.ReviewCount,
			},
		})
	}

	//手动分页
	start := 0
	end := start + limit
	if end > len(output) {
		end = len(output)
	}
	return output[start:end], nil
}

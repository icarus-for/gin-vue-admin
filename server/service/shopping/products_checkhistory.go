package service

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type tempHistory struct {
	SkuId               string  `json:"sku_id"`                 // SKU ID
	ProductId           string  `json:"product_id"`             // 商品ID
	ProductName         string  `json:"product_name"`           // 商品名 (省略後)
	ProductCode         string  `json:"product_code,omitempty"` // 商品コード
	PriceRangeFormatted string  `json:"price_range_formatted"`
	ImageId             *int    `json:"image_id"`
	ImageUrl            *string `json:"image_url"`
	ImageAltText        *string `json:"image_alt_text"`
	AverageRating       float64 `json:"average_rating"` // 平均評価
	ReviewCount         int     `json:"review_count"`
	ViewedAtFormatted   string  `json:"viewed_at_formatted"`
}

type GetHistoryList struct{}

func (g *GetHistoryList) AddHistoryList(userId string, skuId string) (string, error) {
	err := global.GVA_DB.Table("user_viewed_skus").
		Create(map[string]interface{}{
			"user_id":   userId,
			"sku_id":    skuId,
			"viewed_at": time.Now(),
		}).Error
	if err != nil {
		return "", err
	}
	return "添加成功", err
}

func (g *GetHistoryList) DELETECheckList(userId string, skuId string) (string, error) {
	var checkCount int64
	err := global.GVA_DB.Table("user_viewed_skus").
		Where("user_id = ? AND sku_id = ?", userId, skuId).
		Count(&checkCount).Error
	if err != nil {
		return "", err
	}
	if checkCount == 0 {
		return "", fmt.Errorf("浏览历史不存在，无法删除")
	}

	err = global.GVA_DB.Table("user_viewed_skus").
		Where("user_id=? AND sku_id=?", userId, skuId).
		Delete(nil).Error
	if err != nil {
		return "", err
	}
	return "删除浏览历史成功", err
}

func (g *GetHistoryList) GetCheckListInfo(userId string, page int, limit int) (*dto.ViewedSKUListResponse, error) {

	var flatCheck []tempHistory
	offset := (page - 1) * limit
	subQuery := global.GVA_DB.Table("product_skus").
		Distinct().
		Select(`
			product_skus.id AS sku_id,
			product_skus.product_id,
			(
			CASE
			WHEN CHAR_LENGTH(products.name)>10
			THEN CONCAT(LEFT(products.name,10),'...')
			ELSE products.name
			END
			)AS product_name,
			products.product_code AS product_code,
			user_viewed_skus.viewed_at,
			DATE_FORMAT(user_viewed_skus.viewed_at, '%Y-%m-%d %H:%i:%s') AS viewed_at_formatted,
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
			sku_images.id AS image_id,
			sku_images.main_image_url AS image_url,
			sku_images.alt_text AS image_alt_text,
			review_summaries.average_rating,
			COUNT(DISTINCT CASE WHEN product_reviews.status = 'approved' THEN product_reviews.id END) AS review_count
			`).
		Joins("LEFT JOIN products ON Product_skus.product_id =products.id").
		Joins("LEFT JOIN user_viewed_skus ON Product_skus.id =user_viewed_skus.sku_id").
		Joins("LEFT JOIN prices ON product_skus.id = prices.sku_id").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Joins("LEFT JOIN sku_images ON product_skus.id =sku_images.sku_id").
		Joins("LEFT JOIN product_reviews ON products.id = product_reviews.product_id").
		Joins("LEFT JOIN review_summaries ON products.id = review_summaries.product_id").
		Where("user_viewed_skus.user_id=?", userId).
		Group("product_skus.id, product_skus.product_id, products.name, products.product_code, user_viewed_skus.viewed_at, sku_images.id, sku_images.main_image_url, sku_images.alt_text, review_summaries.average_rating")
		//子查询
	query := global.GVA_DB.Table("(?) as sub", subQuery).
		Select(`
				sku_id,
				product_id,
       			product_name,
        		product_code,
        		viewed_at_formatted,
        		CASE
            		WHEN min_price IS NULL THEN NULL
            		WHEN min_price = max_price THEN CONCAT(FORMAT(min_price, 0), '円')
            		ELSE CONCAT(FORMAT(min_price, 0), '～', FORMAT(max_price, 0), '円')
				END AS price_range_formatted,
        		image_id,
        		image_url,
        		image_alt_text,
        		average_rating,
        		review_count
				`)

	//query = query.Offset(offset).Limit(limit)
	err := query.Scan(&flatCheck).Error
	if err != nil {
		return nil, err
	}

	//使用skuid去重
	uniqueCheckList := make(map[string]tempHistory)
	for _, v := range flatCheck {
		if _, exist := uniqueCheckList[v.SkuId]; !exist {
			uniqueCheckList[v.SkuId] = v
		}
	}
	var temphislist []dto.ViewedSKUInfo
	for _, v1 := range uniqueCheckList {
		var image dto.ImageInfo
		if v1.ImageId != nil {
			image = dto.ImageInfo{
				Id:      *v1.ImageId,
				Url:     *v1.ImageUrl,
				AltText: v1.ImageAltText,
			}
		}
		reviewsummary := &dto.ReviewSummaryInfo{
			AverageRating: v1.AverageRating,
			ReviewCount:   v1.ReviewCount,
		}
		temphislist = append(temphislist, dto.ViewedSKUInfo{
			SkuId:               v1.SkuId,
			ProductId:           v1.ProductId,
			ProductName:         v1.ProductName,
			ProductCode:         v1.ProductCode,
			PriceRangeFormatted: v1.PriceRangeFormatted,
			PrimaryImage:        &image,
			ReviewSummary:       reviewsummary,
			ViewedAtFormatted:   v1.ViewedAtFormatted,
		})
	}

	sort.Slice(temphislist, func(i, j int) bool {
		return temphislist[i].ViewedAtFormatted > temphislist[j].ViewedAtFormatted
	})
	//手动分页，效率低
	start := offset
	end := offset + limit
	if start > len(temphislist) {
		start = len(temphislist)
	}
	if end > len(temphislist) {
		end = len(temphislist)
	}
	temphislist = temphislist[start:end]
	//查询收藏的商品的属性
	var checkTotal int64
	err = global.GVA_DB.Table("user_viewed_skus").
		Where("user_id=?", userId).
		Count(&checkTotal).Error
	if err != nil {
		return nil, err
	}
	output := &dto.ViewedSKUListResponse{
		History: temphislist,
		Pagination: dto.PaginationInfo{
			CurrentPage: page,
			Limit:       limit,
			TotalCount:  int(checkTotal),
			TotalPages:  int(math.Ceil(float64(checkTotal) / float64(limit))),
		},
	}
	fmt.Println("最终output:", output)
	fmt.Println("最终temphislist:", temphislist)

	return output, nil
}

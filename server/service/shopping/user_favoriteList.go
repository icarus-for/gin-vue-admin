package service

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type FavoriteSkuInfo struct {
	SkuId            string          `json:"sku_id"`                 // SKU ID
	ProductId        string          `json:"product_id"`             // 商品ID
	ProductName      string          `json:"product_name"`           // 商品名 (省略後)
	ProductCode      string          `json:"product_code,omitempty"` // 商品コード
	Price            *PriceInfo      `json:"price"`                  // 現在の表示価格情報 (Nullable)
	PrimaryImage     *ImageInfo      `json:"primary_image"`          // 代表画像 (Nullable)
	Attributes       []AttributeInfo `json:"attributes"`             // 対象SKUの属性リスト
	AddedAtFormatted string          `json:"added_at_formatted"`     // お気に入り追加日時 (表示用)
}

type PriceInfo struct {
	Amount                  float64  `json:"amount"`
	FormattedAmount         string   `json:"formatted_amount"`
	Type                    string   `json:"type"`
	TypeName                string   `json:"type_name"`
	OriginalAmount          *float64 `json:"original_amount"`
	FormattedOriginalAmount *string  `json:"formatted_original_amount"`
}

type ImageInfo struct {
	Id      int     `json:"image_id"`
	Url     string  `json:"image_url"`
	AltText *string `json:"image_alt_text"`
}

type AttributeInfo struct {
	AttributeId   int     `json:"attribute_id"`
	AttributeName string  `json:"attribute_name"`
	OptionId      *int    `json:"option_id,omitempty"`
	OptionValue   *string `json:"option_value,omitempty"`
	ValueString   *string `json:"value_string,omitempty"`
}

type PaginationInfo struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalCount  int `json:"total_count"` // お気に入りの総件数
	TotalPages  int `json:"total_pages"`
}

type FavoriteSKUListResponse struct {
	Favorites  []FavoriteSkuInfo `json:"favorites"`  // お気に入りSKUリスト
	Pagination PaginationInfo    `json:"pagination"` // ページネーション情報
}

type Favorite struct {
	UserId    string    `json:"user_id"`
	SkuId     string    `json:"sku_id"`
	CreatedAt time.Time `json:"created_at"`
}

type TempFavList struct {
	SkuId                   string   `json:"sku_id"`       // SKU ID
	ProductId               string   `json:"product_id"`   // 商品ID
	ProductName             string   `json:"product_name"` // 商品名 (省略後)
	ProductCode             string   `json:"product_code"`
	Amount                  float64  `json:"amount"`
	FormattedAmount         string   `json:"formatted_amount"`
	Type                    string   `json:"type"`
	TypeName                string   `json:"type_name"`
	OriginalAmount          *float64 `json:"original_amount,omitempty"`
	FormattedOriginalAmount *string  `json:"formatted_original_amount,omitempty"`
	ImageId                 *int     `json:"image_id"`
	ImageUrl                *string  `json:"image_url"`
	ImageAltText            *string  `json:"image_alt_text"`
	AddedAtFormatted        string   `json:"added_at_formatted"`
}

type FavoriteList struct{}

func (g *FavoriteList) AddfavoriteList(userId string, skuId string) (string, error) {
	var addst int64
	err := global.GVA_DB.Table("product_skus").
		Where("id = ?", skuId).
		Count(&addst).Error
	if err != nil {
		return "", err
	}
	if addst == 0 {
		return "", fmt.Errorf("SKU不存在,无法收藏")
	}

	var favCount int64
	err = global.GVA_DB.Table("user_favorite_skus").
		Where("user_id = ? AND sku_id = ?", userId, skuId).
		Count(&favCount).Error
	if err != nil {
		return "", err
	}
	if favCount > 0 {
		return "", fmt.Errorf("已收藏过该SKU")
	}

	err = global.GVA_DB.Table("user_favorite_skus").
		Create(map[string]interface{}{
			"user_id":    userId,
			"sku_id":     skuId,
			"created_at": time.Now(),
		}).Error
	if err != nil {
		return "", err
	}
	return "添加收藏成功", err
}

func (g *FavoriteList) DELETEfavoriteList(userId string, skuId string) (string, error) {
	var favCount int64
	err := global.GVA_DB.Table("user_favorite_skus").
		Where("user_id = ? AND sku_id = ?", userId, skuId).
		Count(&favCount).Error
	if err != nil {
		return "", err
	}
	if favCount == 0 {
		return "", fmt.Errorf("收藏不存在，无法删除")
	}

	err = global.GVA_DB.Table("user_favorite_skus").
		Where("user_id=? AND sku_id=?", userId, skuId).
		Delete(nil).Error
	if err != nil {
		return "", err
	}
	return "删除收藏成功", err
}

func (g *FavoriteList) GetFavoriteListInfo(userId string, page int, limit int, sort3 string) (interface{}, error) {

	var flatFavorites []TempFavList
	offset := (page - 1) * limit
	query := global.GVA_DB.Table("product_skus").
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
			user_favorite_skus.created_at,
			DATE_FORMAT(user_favorite_skus.created_at, '%Y-%m-%d %H:%i:%s') AS added_at_formatted,
			prices.price AS amount,
			CONCAT(FORMAT(prices.price, 0), '円') AS formatted_amount,
			price_types.type_code AS type,
			price_types.name AS type_name,
		   COALESCE(
            (
                CASE
                    WHEN prices.price_type_id IN (2,3) THEN (
                        SELECT p.price
                        FROM prices p
                        WHERE p.sku_id=prices.sku_id
                        AND p.price_type_id = 1
                        AND (p.start_date IS NULL OR p.start_date <= NOW())
                        AND (p.end_date IS NULL OR p.end_date >= NOW())
                        LIMIT 1
                    )
                    ELSE NULL
                END
            ), 
            NULL
        ) AS original_amount,
        COALESCE(
            (
                CASE
                    WHEN prices.price_type_id IN (2,3) THEN (
                        SELECT CONCAT(FORMAT(prices.price/100, 2), '円')
                        FROM prices p
                        WHERE p.sku_id=prices.sku_id
                        AND p.price_type_id = 1
                        AND (p.start_date IS NULL OR p.start_date <= NOW())
                        AND (p.end_date IS NULL OR p.end_date >= NOW())
                        LIMIT 1
                    )
                    ELSE NULL
                END
            ), 
            NULL
        ) AS formatted_original_amount,
			sku_images.id AS image_id,
			sku_images.main_image_url AS image_url,
			sku_images.alt_text AS image_alt_text
			`).
		Joins("LEFT JOIN products ON Product_skus.product_id =products.id").
		Joins("LEFT JOIN user_favorite_skus ON Product_skus.id =user_favorite_skus.sku_id").
		Joins("LEFT JOIN prices ON product_skus.id = prices.sku_id").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Joins("LEFT JOIN sku_images ON product_skus.id =sku_images.sku_id").
		Where("user_favorite_skus.user_id=?", userId)

	//排序
	switch sort3 {
	case "newest":
		query = query.Order("user_favorite_skus.created_at DESC")
	case "oldest":
		query = query.Order("user_favorite_skus.created_at ASC")
	default:
		query = query.Order("user_favorite_skus.created_at DESC")
	}
	//分页
	//query = query.Offset(offset).Limit(limit)
	err := query.Scan(&flatFavorites).Error
	if err != nil {
		return nil, err
	}
	uniqueFavorites := make(map[string]TempFavList) // 使用 SKU ID 去重
	for _, v := range flatFavorites {
		if _, exists := uniqueFavorites[v.SkuId]; !exists {
			uniqueFavorites[v.SkuId] = v
		}
	}
	var tempfavlist []FavoriteSkuInfo
	for _, v1 := range uniqueFavorites {
		var image ImageInfo
		if v1.ImageId != nil {
			image = ImageInfo{
				Id:      *v1.ImageId,
				Url:     *v1.ImageUrl,
				AltText: v1.ImageAltText,
			}
		}
		price := &PriceInfo{
			Amount:                  v1.Amount,
			FormattedAmount:         v1.FormattedAmount,
			Type:                    v1.Type,
			TypeName:                v1.TypeName,
			OriginalAmount:          v1.OriginalAmount,
			FormattedOriginalAmount: v1.FormattedOriginalAmount,
		}

		tempfavlist = append(tempfavlist, FavoriteSkuInfo{
			SkuId:            v1.SkuId,
			ProductId:        v1.ProductId,
			ProductName:      v1.ProductName,
			ProductCode:      v1.ProductCode,
			Price:            price,
			PrimaryImage:     &image,
			AddedAtFormatted: v1.AddedAtFormatted,
		})
	}

	//sort.Slice 对于切片进行自定义排序
	sort.Slice(tempfavlist, func(i, j int) bool {
		if sort3 == "oldest" {
			return tempfavlist[i].AddedAtFormatted < tempfavlist[j].AddedAtFormatted
		}
		return tempfavlist[i].AddedAtFormatted > tempfavlist[j].AddedAtFormatted
	})
	//手动分页，效率低
	start := offset
	end := offset + limit
	if start > len(tempfavlist) {
		start = len(tempfavlist)
	}
	if end > len(tempfavlist) {
		end = len(tempfavlist)
	}
	tempfavlist = tempfavlist[start:end]
	//查询收藏的商品的属性
	for v := range tempfavlist {
		var attributes []AttributeInfo
		err := global.GVA_DB.Table("sku_values").
			Distinct().
			Select(`sku_values.attribute_id,
			attributes.name AS attribute_name,
			attribute_options.id AS option_id,
			COALESCE(attribute_options.value,sku_values.value_string,CAST(sku_values.value_number AS CHAR),
			CASE
			WHEN value_boolean=1 THEN TRUE
			ELSE FALSE
			END
			)AS option_value,
			sku_values.value_number
		`).
			Joins("LEFT JOIN attributes ON sku_values.attribute_id = attributes.id").
			Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
			Order("sku_values.attribute_id ASC").
			Where("sku_values.sku_id=?", tempfavlist[v].SkuId).
			Scan(&attributes).Error
		if err == nil {
			tempfavlist[v].Attributes = attributes
		}
	}
	var favTotal int64
	err = global.GVA_DB.Table("user_favorite_skus").
		Where("user_id=?", userId).
		Count(&favTotal).Error
	if err != nil {
		return nil, err
	}
	output := &FavoriteSKUListResponse{
		Favorites: tempfavlist,
		Pagination: PaginationInfo{
			CurrentPage: page,
			Limit:       limit,
			TotalCount:  int(favTotal),
			TotalPages:  int(math.Ceil(float64(favTotal) / float64(limit))),
		},
	}

	return output, nil

}

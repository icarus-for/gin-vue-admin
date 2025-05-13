package service

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
)

type tempCartResponse struct {
	SkuId                   string   `json:"sku_id"`       // SKU ID
	ProductId               string   `json:"product_id"`   // 商品ID
	ProductName             string   `json:"product_name"` // 商品名 (省略後)
	ProductCode             string   `json:"product_code"`
	Quantity                int      `json:"quantity"`
	Amount                  float64  `json:"amount"`
	FormattedAmount         string   `json:"formatted_amount"`
	Type                    string   `json:"type"`
	TypeName                string   `json:"type_name"`
	OriginalAmount          *float64 `json:"original_amount,omitempty"`
	FormattedOriginalAmount *string  `json:"formatted_original_amount,omitempty"`
	SubtotalFormatted       string   `json:"subtotal_formatted"`
	ImageId                 *int     `json:"image_id"`
	ImageUrl                *string  `json:"image_url"`
	ImageAltText            *string  `json:"image_alt_text"`
	TotalItemsCount         int      `json:"total_items_count"`
	TotalAmount             float64  `json:"total_amount"`           // 合計金額 (計算用数値)
	TotalAmountFormatted    string   `json:"total_amount_formatted"` // 合計金額 (表示用文字列 例: "55,880円")
	StockStatus             string   `json:"stock_status"`
}
type userCartItem struct {
	UserId   string `gorm:"column:user_id" json:"user_id"` // 用户ID
	SkuId    string `gorm:"column:sku_id" json:"sku_id"`   // 商品SKU ID
	Quantity int    `gorm:"column:quantity" json:"quantity"`
}

type inventoryItem struct {
	SkuId    string `gorm:"column:sku_id" json:"sku_id"` // 商品SKU ID
	Quantity int    `gorm:"column:quantity" json:"quantity"`
}
type GetCart struct{}

func (g *GetCart) AddCartListInfo(userId string, skuId string, quantity int) (string, error) {
	var addst int64
	err := global.GVA_DB.Table("product_skus").
		Select("id").
		Where("id=?", skuId).
		Count(&addst).Error
	if err != nil {
		return "", err
	}
	if addst == 0 {
		return "", fmt.Errorf("SKU不存在,无法添加")
	}

	var cartCount userCartItem
	err = global.GVA_DB.Table("user_cart_items").
		Where("user_id=? AND sku_id=?", userId, skuId).
		First(&cartCount).Error
	if err != nil && err != gorm.ErrRecordNotFound { //ErrRecordNotFound 是gorm的一个常量，表示在查询时未找到任何匹配的记录
		//不意味着查询失败，而只是表示没有找到符合条件的记录
		return "", err
	}
	//查询该商品的总库存
	var itemCount inventoryItem
	err = global.GVA_DB.Table("inventory").
		Select("quantity").
		Where("sku_id = ?", skuId).
		Scan(&itemCount).Error
	if err != nil {
		return "", err
	}

	if quantity < 0 {
		quantity = 1
	}
	var carCount userCartItem
	err = global.GVA_DB.Table("user_cart_items").
		Where("user_id=? AND sku_id=?", userId, skuId).
		First(&cartCount).Error

	if err == gorm.ErrRecordNotFound {
		// 没有记录，插入新购物车项
		if quantity > itemCount.Quantity {
			return "カートに入れる数は在庫より多い", nil
		}
		newCartItem := userCartItem{
			UserId:   userId,
			SkuId:    skuId,
			Quantity: quantity,
		}
		err = global.GVA_DB.Table("user_cart_items").
			Create(&newCartItem).Error
		if err != nil {
			return "", err
		}
		return "カートに追加しました", nil
	} else if err != nil {
		return "", err
	}

	// 已经有记录，更新数量
	newQuantity := carCount.Quantity + quantity
	if newQuantity > itemCount.Quantity {
		return "カートに入れる数は在庫より多い", nil
	}

	err = global.GVA_DB.Table("user_cart_items").
		Where("user_id=? AND sku_id=?", userId, skuId).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
	if err != nil {
		return "", err
	}

	return "カートに追加しました", nil
}

func (g *GetCart) DelCartListInfo(userId string, skuId string) (string, error) {
	var cartCount int64
	err := global.GVA_DB.Table("user_cart_items").
		Where("user_id=? AND sku_id=?", userId, skuId).
		Count(&cartCount).Error
	if err != nil {
		return "", err
	}
	if cartCount == 0 {
		return "", fmt.Errorf("sku不存在,无法删除")
	}
	err = global.GVA_DB.Table("user_cart_items").
		Where("user_id=? AND sku_id=?", userId, skuId).
		Delete(nil).Error
	if err != nil {
		return "", err
	}
	return "カートから削除しました。", err
}

func (g *GetCart) CartQuantityChangeInfo(userId string, skuId string, quantity int) (string, error) {
	var changest int64
	err := global.GVA_DB.Table("product_skus").
		Select("id").
		Where("id=?", skuId).
		Count(&changest).Error
	if err != nil {
		return "", err
	}
	if changest == 0 {
		return "", fmt.Errorf("sku不存在,无法删除")
	}
	//变更商品数量
	err = global.GVA_DB.Table("user_cart_items").
		Select("quantity").
		Where("user_id=? AND sku_id=?", userId, skuId).
		Update("quantity", quantity).Error
	if err != nil {
		return "", err
	}
	//还是先查出总库存
	var itemCount inventoryItem
	err = global.GVA_DB.Table("inventory").
		Select("quantity").
		Where("sku_id = ?", skuId).
		Scan(&itemCount).Error
	if err != nil {
		return "", err
	}
	// //然后把变更之后的库存查出来
	var addedCount userCartItem
	err = global.GVA_DB.Table("user_cart_items").
		Select("quantity").
		Where("user_id=? AND sku_id=?", userId, skuId).
		First(&addedCount).Error
	if err != nil {
		return "", err
	}
	if addedCount.Quantity > itemCount.Quantity {
		return "カートに入れる数は在庫より多い", err
	}
	return "数が変更しました.", err
}

func (g *GetCart) GetCartListInfo(userId string) (*dto.CartResponse, error) {
	var tempCart []tempCartResponse
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
			user_cart_items.quantity,
			prices.price AS amount,
			CONCAT(FORMAT(prices.price,0), '円') AS formatted_amount,
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
                        SELECT CONCAT(FORMAT(p.price, 0), '円')
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
			CONCAT(FORMAT(user_cart_items.quantity * prices.price , 0), '円') AS subtotal_formatted,
			sku_images.id AS image_id,
			sku_images.main_image_url AS image_url,
			sku_images.alt_text AS image_alt_text,
			(
			CASE
			WHEN inventory.quantity > 5 THEN 'available'
			WHEN inventory.quantity > 0 THEN 'low_stock'
			ELSE 'out_of_stock'
			END
			) AS stock_status,
			user_cart_items.added_at
			`).
		Joins("LEFT JOIN user_cart_items ON Product_skus.id =user_cart_items.sku_id").
		Joins("LEFT JOIN products ON Product_skus.product_id =products.id").
		Joins("LEFT JOIN prices ON product_skus.id = prices.sku_id").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Joins("LEFT JOIN sku_images ON product_skus.id =sku_images.sku_id").
		Joins("LEFT JOIN inventory ON product_skus.id = inventory.sku_id").
		Group("product_skus.id, product_skus.product_id, products.name, products.product_code, user_cart_items.quantity, prices.price, prices.price_type_id, price_types.type_code, price_types.name, sku_images.id, sku_images.main_image_url, sku_images.alt_text, inventory.quantity, user_cart_items.added_at").
		Where("user_cart_items.user_id=?", userId).
		Order("`user_cart_items`.`added_at` DESC")
	err := query.Scan(&tempCart).Error
	if err != nil {
		return nil, err
	}

	var TotalItemsCount int
	var TotalAmount float64

	uniqueCareList := make(map[string]tempCartResponse)
	for _, v := range tempCart {
		if _, exists := uniqueCareList[v.SkuId]; !exists {
			uniqueCareList[v.SkuId] = v
		}
	}
	var temCart []dto.CartItemInfo
	for _, v1 := range uniqueCareList {
		TotalItemsCount += v1.Quantity
		TotalAmount += v1.Amount * float64(v1.Quantity)
		var image dto.ImageInfo
		if v1.ImageId != nil {
			image = dto.ImageInfo{
				Id:      *v1.ImageId,
				Url:     *v1.ImageUrl,
				AltText: v1.ImageAltText,
			}
		}
		price := &dto.PriceInfo{
			Amount:                  v1.Amount,
			FormattedAmount:         v1.FormattedAmount,
			Type:                    v1.Type,
			TypeName:                v1.TypeName,
			OriginalAmount:          v1.OriginalAmount,
			FormattedOriginalAmount: v1.FormattedOriginalAmount,
		}
		temCart = append(temCart, dto.CartItemInfo{
			SkuId:             v1.SkuId,
			ProductId:         v1.ProductId,
			ProductName:       v1.ProductName,
			ProductCode:       v1.ProductCode,
			Quantity:          v1.Quantity,
			SubtotalFormatted: v1.SubtotalFormatted,
			Price:             price,
			PrimaryImage:      &image,
			StockStatus:       v1.StockStatus,
		})
	}
	for v := range temCart {
		var attributes []dto.AttributeInfo
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
			Where("sku_values.sku_id=?", temCart[v].SkuId).
			Scan(&attributes).Error
		if err == nil {
			temCart[v].Attributes = attributes
		}
	}
	Output := &dto.CartResponse{
		Items:                temCart,
		TotalItemsCount:      TotalItemsCount,
		TotalAmount:          TotalAmount,
		TotalAmountFormatted: fmt.Sprintf("%.2f円", TotalAmount),
	}
	return Output, err
}

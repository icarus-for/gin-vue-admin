package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ProductInfo struct {
	SkuId           string          `json:"sku_id"`
	ProductId       string          `json:"product_id"`
	Weight          string          `json:"weight"`
	Size            string          `json:"size"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	IsTaxable       string          `json:"is_taxable"`
	ProductCode     string          `json:"product_code"`
	Price           Price           `json:"price" gorm:"embedded"`
	Quantity        int             `json:"quantity"`
	ChannelCode     string          `json:"channel_code"`
	PointOfSale     string          `json:"point_of_sale"`
	CategoriesName  string          `json:"categories_name"`
	Id              int             `json:"id"`
	Url             string          `json:"url"`
	AltText         string          `json:"alt_text"`
	Attributes      []Attribute     `gorm:"-" json:"attributes"`
	VariantOptions  []VariantOption `gorm:"-" json:"variant_options"`
	RelatedCategory []string        `gorm:"-" json:"related_category"`
}

type Attribute struct {
	AttributeId   int    `json:"attribute_id"`
	AttributeName string `json:"attribute_name"`
	Value         string `json:"value"`
}

type Option struct {
	OptionId     int      `json:"option_id"`
	OptionValue  string   `json:"option_value"`
	OptionCode   string   `json:"option_code"`
	LinkedSkuIds []string `json:"linked_sku_ids" gorm:"type:varchar(255)"`
}
type VariantOption struct {
	AttributeId   int      `json:"attribute_id"`
	AttributeName string   `json:"attribute_name"`
	AttributeCode string   `json:"attribute_code"`
	Options       []Option `json:"options"`
}

type Price struct {
	Amount                  float64    `json:"amount" gorm:"column:amount"`
	FormattedAmount         string     `json:"formatted_amount" gorm:"column:formatted_amount"`
	Type                    string     `json:"type" gorm:"column:type"`
	TypeName                string     `json:"type_name" gorm:"column:type_name"`
	StartDate               *time.Time `json:"start_date" gorm:"column:start_date"`
	EndDate                 *time.Time `json:"end_date" gorm:"column:end_date"`
	OriginalAmount          float64    `json:"original_amount" gorm:"column:original_amount"`
	FormattedOriginalAmount string     `json:"formatted_original_amount" gorm:"column:formatted_original_amount"`
}

type Category struct {
	// Id   int    `json:"id"`
	Name string `json:"name"`
}

var neededSkuId string

func GetProductInfoList(sku_id string, productcode string) (interface{}, error) {
	if sku_id != "" {
		neededSkuId = sku_id
	} //当输入sku_id时直接根据sku_id搜索
	//当没有sku_id时，根据输入的商品code，然后将其默认的daufault_sku_id作为sku_id来使用
	if sku_id == "" && productcode != "" {
		var gpil string
		err := global.GVA_DB.Table("Products").Select("default_sku_id").
			Where("product_code=?", productcode).Scan(&gpil).Error
		if err != nil {
			return nil, err
		}
		neededSkuId = gpil
	}
	if neededSkuId == "" {
		return nil, fmt.Errorf("不正な商品識別子です")
	}

	var productInfo ProductInfo
	err := global.GVA_DB.Table("product_skus").
		Select(`
			Product_skus.id as sku_id,
			Product_skus.product_id,
			CONCAT(Product_skus.weight,'kg')AS weight,
			TRIM(BOTH 'x' FROM CONCAT_WS('x',
			IF(Product_skus.width IS NOT NULL AND Product_skus.width != '',CONCAT(Product_skus.width,'cm'),NULL),
			IF(Product_skus.height IS NOT NULL AND Product_skus.height != '',CONCAT(Product_skus.height,'cm'),NULL),
			IF(Product_skus.depth IS NOT NULL AND Product_skus.depth != '',CONCAT(Product_skus.depth,'cm'),NULL)
			))AS size,
			products.name,
			(
			CASE
			WHEN CHAR_LENGTH(products.description)>20
			THEN CONCAT(LEFT(products.description,20),'...')
			ELSE products.description
			END
			)AS description,
			(
			CASE
			WHEN products.is_taxable = 1 THEN 'TRUE'
			ELSE 'FALSE'
			END
			)AS is_taxable,		
			IF(member_price.price IS NOT NULL, member_price.price,
			IF(discount_price.price IS NOT NULL, discount_price.price,
			IF(regular_price.price IS NOT NULL, regular_price.price, 0)
			)
			) AS amount,
			IF(member_price.price IS NOT NULL, CONCAT(FORMAT(member_price.price, 0), '円'),
			IF(discount_price.price IS NOT NULL, CONCAT(FORMAT(discount_price.price, 0), '円'),
			IF(regular_price.price IS NOT NULL, CONCAT(FORMAT(regular_price.price, 0), '円'), '')
			)
			) AS formatted_amount,
			IF(member_price.price IS NOT NULL, 'member_special',
			IF(discount_price.price IS NOT NULL, 'sale',
			IF(regular_price.price IS NOT NULL, 'regular', '')
			)
			) AS type,
			IF(member_price.price IS NOT NULL, '会員特別価格',
			IF(discount_price.price IS NOT NULL, 'セール価格',
			IF(regular_price.price IS NOT NULL, '通常価格', '')
			)
			) AS type_name,
			IF(discount_price.price IS NOT NULL, discount_price.start_date, NULL) AS start_date,
        	IF(discount_price.price IS NOT NULL, discount_price.end_date, NULL) AS end_date,
			IF((member_price.price IS NOT NULL OR discount_price.price IS NOT NULL), regular_price.price, NULL) AS original_amount,
			IF((member_price.price IS NOT NULL OR discount_price.price IS NOT NULL), CONCAT(FORMAT(regular_price.price, 0), '円'), NULL) AS formatted_original_amount,
			products.product_code,
			inventory.quantity,
			sales_channels.channel_code,
			sales_channels.name as point_of_sale,
			categories.name as categories_name
		`).
		Joins("LEFT JOIN products ON Product_skus.product_id =products.id").
		Joins("LEFT JOIN prices ON Product_skus.id =prices.sku_id").
		Joins("LEFT JOIN prices AS member_price ON product_skus.id = member_price.sku_id AND member_price.price_type_id = 3").
		Joins("LEFT JOIN prices AS discount_price ON product_skus.id = discount_price.sku_id AND discount_price.price_type_id = 2").
		Joins("LEFT JOIN prices AS regular_price ON product_skus.id = regular_price.sku_id AND regular_price.price_type_id = 1").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Joins("LEFT JOIN inventory ON Product_skus.id =inventory.sku_id").
		Joins("RIGHT JOIN categories ON products.category_id =categories.id").
		Joins("LEFT JOIN sku_values ON Product_skus.id =sku_values.sku_id").
		Joins("LEFT JOIN sku_availability ON Product_skus.id = sku_availability.sku_id").
		Joins("RIGHT JOIN sales_channels ON sku_availability.sales_channel_id =sales_channels.id").
		Joins("LEFT JOIN attribute_options ON sku_values.attribute_id =attribute_options.attribute_id").
		Joins("LEFT JOIN sku_images ON product_skus.id =sku_images.sku_id").
		Where("Product_skus.id = ?", neededSkuId).
		Scan(&productInfo).Error

	if err != nil {
		return nil, err
	}
	if productInfo.SkuId == "" {
		return nil, fmt.Errorf("不正なSKU ID形式です")
	}

	var ProductId string
	err = global.GVA_DB.Table("product_skus").Select("product_id").Where("id=?", neededSkuId).
		Scan(&ProductId).Error
	if err != nil {
		return nil, err
	}
	var RelatedProduct []string
	err = global.GVA_DB.Table("product_skus").Select("id").Where("product_id=?", ProductId).
		Scan(&RelatedProduct).Error
	if err != nil {
		return nil, err
	}

	var attributes []struct {
		AttributeId   int
		AttributeName string
		AttributeCode string
		OptionId      int
		OptionCode    string
		OptionValue   string
		SkuId         string
	}
	err = global.GVA_DB.Table("sku_values").
		Distinct(). //sql语句 去重复
		Select(`
			sku_values.attribute_id,
			attributes.name AS attribute_name,
			attribute_options.option_code,
			COALESCE(attribute_options.value,sku_values.value_string,CAST(sku_values.value_number AS CHAR),
			CASE
			WHEN value_boolean=1 THEN TRUE
			ELSE FALSE
			END
			)AS option_value,
			attribute_options.id AS option_id,
			attribute_options.sort_order,
			sku_values.sku_id,
			sku_values.value_string,
			sku_values.value_number
		`).
		Joins("LEFT JOIN attributes ON sku_values.attribute_id = attributes.id").
		Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
		Where("sku_values.sku_id IN (?)", neededSkuId).
		Order("sku_values.attribute_id ASC").
		Scan(&attributes).Error

	if err != nil {
		return nil, err
	}

	var attrs []Attribute
	for _, v := range attributes {
		attrs = append(attrs, Attribute{
			AttributeId:   v.AttributeId,
			AttributeName: v.AttributeName,
			Value:         v.OptionValue,
		})
	}
	productInfo.Attributes = attrs

	err = global.GVA_DB.Table("product_skus").Select("id").Where("product_id=?", ProductId).Scan(&RelatedProduct).Error
	if err != nil {
		return nil, err
	}

	// 处理变体选项（查询与当前商品关联的其他 SKU 的属性）
	variantOptionsMap := make(map[int]*VariantOption)

	var relatedAttributes []struct {
		AttributeId   int
		AttributeName string
		AttributeCode string
		OptionId      int
		OptionCode    string
		OptionValue   string
		SkuId         string
	}
	err = global.GVA_DB.Table("sku_values").
		Distinct().
		Select(`
			sku_values.attribute_id,
			attributes.name AS attribute_name,
			attribute_options.option_code,
			COALESCE(attribute_options.value, sku_values.value_string, CAST(sku_values.value_number AS CHAR),
			CASE
				WHEN value_boolean = 1 THEN TRUE
				ELSE FALSE
			END) AS option_value,
			attribute_options.id AS option_id,
			attribute_options.sort_order,
			sku_values.sku_id,
			sku_values.value_string,
			sku_values.value_number
		`).
		Joins("LEFT JOIN attributes ON sku_values.attribute_id = attributes.id").
		Joins("LEFT JOIN attribute_options ON sku_values.option_id = attribute_options.id").
		Where("sku_values.sku_id IN (?)", RelatedProduct). // 查询关联产品的 SKU 属性
		Order("sku_values.attribute_id ASC").
		Scan(&relatedAttributes).Error

	if err != nil {
		return nil, err
	}

	for _, v := range relatedAttributes {
		// 检查是否已存在相同的属性 ID
		vo, exists := variantOptionsMap[v.AttributeId]
		if !exists {
			// 如果不存在，创建新的 VariantOption
			vo = &VariantOption{
				AttributeId:   v.AttributeId,
				AttributeName: v.AttributeName,
				AttributeCode: v.AttributeCode,
				Options:       []Option{},
			}
			variantOptionsMap[v.AttributeId] = vo
		}
		linkedSkus := []string{v.SkuId}

		// 检查当前的 Option 是否已存在
		optionExists := false
		for i := range vo.Options {
			if vo.Options[i].OptionCode == v.OptionCode {
				// 如果已经存在相同的 OptionCode，将当前的 SkuId 添加到该 Option 中
				vo.Options[i].LinkedSkuIds = append(vo.Options[i].LinkedSkuIds, linkedSkus...)
				optionExists = true
				break
			}
		}
		if !optionExists {
			option := Option{
				OptionId:     v.OptionId,
				OptionValue:  v.OptionValue,
				OptionCode:   v.OptionCode,
				LinkedSkuIds: linkedSkus,
			}
			vo.Options = append(vo.Options, option)
		}
	}
	var variantOptions []VariantOption
	for _, vo := range variantOptionsMap {
		// 按照 option 的 sort_order 排序
		variantOptions = append(variantOptions, *vo)
	}

	sort.SliceStable(variantOptions, func(i, j int) bool {
		return variantOptions[i].AttributeId < variantOptions[j].AttributeId
	})

	productInfo.VariantOptions = variantOptions

	var relatedCategory []Category
	var categoryId int

	err = global.GVA_DB.Table("categories").
		Select("parent_id").
		Where("name=?", productInfo.CategoriesName).
		Scan(&categoryId).Error
	if err != nil {
		return nil, err
	}

	err = global.GVA_DB.Table("categories").
		Select("name").
		Where("parent_id=?", categoryId).
		Scan(&relatedCategory).Error
	if err != nil {
		return nil, err
	}
	relatedCategoriesName := make([]string, len(relatedCategory))
	for i := range relatedCategory {
		relatedCategoriesName[i] = relatedCategory[i].Name
	}

	productInfo.RelatedCategory = relatedCategoriesName

	return productInfo, nil
}

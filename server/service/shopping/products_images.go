package service

import (
	"regexp"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SkuImageInfo struct {
	//SkuId        string  `json:"sku_id"`
	ID           int     `json:"id"`                 // 画像ペアID
	MainImageURL string  `json:"main_image_url"`     // メイン画像URL
	ThumbnailURL string  `json:"thumbnail_url"`      // サムネイル画像URL
	AltText      *string `json:"alt_text,omitempty"` // 代替テキスト (Nullable)
	SortOrder    int     `json:"sort_order"`         // 表示順
}

type GetImages struct{}

func (g *GetImages) GetProductsImagesInfo(sku_id string) (interface{}, error) {

	fakeCode := regexp.MustCompile(`^[0-9a-zA-Z-]+$`)
	if !fakeCode.MatchString(sku_id) {
		return response{
			Code: "INVALID_PARAMETER",
			Msg:  "不正な商品識別子です。",
			Data: nil,
		}, nil
	}

	var SkuId string
	err := global.GVA_DB.Table("sku_images").
		Where("sku_id = ?", sku_id).
		Select("sku_id").
		Scan(&SkuId).Error
	if err != nil || SkuId == "" {
		return response{
			Code: "NOT FOUND",
			Msg:  "商品が見つかりません。",
			Data: nil,
		}, nil
	}

	var skuImageInfo []SkuImageInfo
	err = global.GVA_DB.Table("sku_images").
		Select(`
			id,
			main_image_url,
			thumbnail_url,
			alt_text,
			sort_order
			`).
		Where("sku_id=?", SkuId).
		Scan(&skuImageInfo).Error
	if err != nil {
		return nil, err
	}
	return skuImageInfo, err
}

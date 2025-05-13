package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type GetCoordinateList struct{}

func (g *GetCoordinateList) GetCoordinateListInfo(productCode string, limit int) ([]dto.CoordinateSetTeaserInfo, error) {
	var productId string
	err := global.GVA_DB.Table("products").
		Select("id").
		Where("product_code=?", productCode).
		Scan(&productId).Error
	if err != nil {
		return nil, err
	}

	var coo []dto.CoordinateSetTeaserInfo
	db := global.GVA_DB.Table("coordinate_sets").
		Select(`
			coordinate_set_items.coordinate_set_id AS set_id,
			coordinate_sets.theme_image_url AS set_theme_image_url,
			coordinate_sets.contributor_nickname,
			coordinate_sets.contributor_avatar_url,
			coordinate_sets.contributor_store_name
			`).
		Joins("LEFT JOIN coordinate_set_items ON coordinate_sets.id = coordinate_set_items.coordinate_set_id").
		Order("coordinate_sets.posted_at DESC").
		Where("product_id=?", productId)

	if limit > 0 {
		db = db.Limit(limit)
	}
	err = db.Scan(&coo).Error
	if err != nil {
		return nil, err
	}
	return coo, nil
}

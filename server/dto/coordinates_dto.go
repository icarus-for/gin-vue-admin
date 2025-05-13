package dto

type CoordinateSetTeaserListResponse struct {
	Coordinates []CoordinateSetTeaserInfo `json:"coordinates"` // コーディネートセット概要のリスト
}

// CoordinateSetTeaserInfo 商品詳細ページに表示するコーディネートセット概要DTO
type CoordinateSetTeaserInfo struct {
	SetID            string  `json:"set_id"`                        // コーディネートセットID
	SetThemeImageURL *string `json:"set_theme_image_url,omitempty"` // コーディネートセットのテーマ画像URL
	//CoordinateSetItems   []CoordinateSetItem `json:"coordinate_set_items"`
	ContributorNickname  string  `json:"contributor_nickname"`             // 投稿者ニックネーム
	ContributorAvatarURL *string `json:"contributor_avatar_url,omitempty"` // 投稿者頭像URL
	ContributorStoreName *string `json:"contributor_store_name,omitempty"` // 投稿者所属店名
}

type CoordinateSetItem struct {
	ProductId   string `json:"productid"`
	DisplayText string `json:"display_text"`
}

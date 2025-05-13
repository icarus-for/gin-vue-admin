package dto

type RelatedProductInfo struct {
	ProductID           string             `json:"product_id"`                      // 関連商品のID
	ProductCode         string             `json:"product_code,omitempty"`          // 関連商品のコード
	ProductName         string             `json:"product_name"`                    // 関連商品の名称 (省略後)
	PriceRangeFormatted string             `json:"price_range_formatted,omitempty"` // ★価格帯文字列 (例: "2,990～3,990円")
	IsOnSale            bool               `json:"is_on_sale"`                      // ★値下げフラグ
	ReviewSummary       *ReviewSummaryInfo `json:"review_summary,omitempty"`        // ★レビュー集計情報 (Nullable)
	ThumbnailImageURL   *string            `json:"thumbnail_image_url,omitempty"`   // 代表SKUのサムネイル画像URL (Nullable)
}

// ReviewSummaryInfo レビュー集計情報のDTO (関連商品用)
type ReviewSummaryInfo struct {
	AverageRating float64 `json:"average_rating,omitempty"` // 平均評価
	ReviewCount   int     `json:"review_count,omitempty"`   // レビュー件数
}

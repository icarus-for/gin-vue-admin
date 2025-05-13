package products

type ReviewSummary struct {
	AverageRating float64 `json:"average_rating"` // 平均評価
	ReviewCount   int     `json:"review_count"`   // 承認済みレビュー総数
	Rating1Count  int     `json:"rating_1_count"` // 星1の数
	Rating2Count  int     `json:"rating_2_count"` // 星2の数
	Rating3Count  int     `json:"rating_3_count"` // 星3の数
	Rating4Count  int     `json:"rating_4_count"` // 星4の数
	Rating5Count  int     `json:"rating_5_count"` // 星5の数
}

// ReviewInfo 個々のレビュー情報 (修正: image_urls, helpful_count 追加)
type ReviewInfo struct {
	ID                 int64    `json:"id"`                   // レビューID
	Nickname           string   `json:"nickname"`             // ニックネーム
	Rating             int      `json:"rating"`               // 評価 (1-5)
	Title              *string  `json:"title,omitempty"`      // タイトル (Nullable)
	Comment            string   `json:"comment"`              // 本文
	CreatedAtFormatted string   `json:"created_at_formatted"` // 表示用投稿日時 (例: "2023年10月26日")
	ImageUrls          []string `json:"image_urls,omitempty"` // ★添付画像URLリスト (画像がない場合は空配列 or 省略)
	HelpfulCount       int      `json:"helpful_count"`        // ★参考になった数
	// IsHelpfulByUser   *bool    `json:"is_helpful_by_user,omitempty"` // ★(オプション) ログインユーザーが参考になったを押したか (Nullable)
}

// PaginationInfo ページネーション情報 (変更なし)
type PaginationInfo struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalCount  int `json:"total_count"`
	TotalPages  int `json:"total_pages"`
}

package service

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type Rating struct {
	AverageRating float64 `json:"average_rating" gorm:"acolumn:verage_rating"` // 平均評価
	ReviewCount   int     `json:"review_count" gorm:"column:review_count"`     // 承認済みレビュー総数
	Rating1Count  int     `json:"rating_1_count" gorm:"column:rating_1_count"` // 星1の数
	Rating2Count  int     `json:"rating_2_count" gorm:"column:rating_2_count"` // 星2の数
	Rating3Count  int     `json:"rating_3_count" gorm:"column:rating_3_count"` // 星3の数
	Rating4Count  int     `json:"rating_4_count" gorm:"column:rating_4_count"` // 星4の数
	Rating5Count  int     `json:"rating_5_count" gorm:"column:rating_5_count"` // 星5の数
}

type Review struct {
	ID                 int64     `json:"id" gorm:"id"`                                     // レビューID
	Nickname           string    `json:"nickname" gorm:"nickname"`                         // ニックネーム
	Rating             int       `json:"rating" gorm:"rating"`                             // 評価 (1-5)
	Title              *string   `json:"title,omitempty" gorm:"title"`                     // タイトル (Nullable)
	Comment            string    `json:"comment" gorm:"comment"`                           // 本文
	CreatedAtFormatted string    `json:"created_at_formatted" gorm:"created_at_formatted"` // 表示用投稿日時 (例: "2023年10月26日")
	ImageUrls          []string  `json:"image_urls,omitempty" gorm:"image_urls"`           // ★添付画像URLリスト (画像がない場合は空配列 or 省略)
	HelpfulCount       int       `json:"helpful_count" gorm:"helpful_count"`
	RawImageUrls       string    `json:"-" gorm:"column:image_urls"` // ★参考になった数
	CreatedAt          time.Time `json:"-" gorm:"column:created_at"`
}
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalCount  int64 `json:"total_count"`
	TotalPages  int   `json:"total_pages"`
}

type response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type Getreview struct{}

func (g *Getreview) GetProductReviewInfo(productcode string, page int, limit int, sort1 string, rating int) (interface{}, error) {
	var reviews []Review
	var summary Rating

	fakeCode := regexp.MustCompile(`^[0-9]+$`)
	if !fakeCode.MatchString(productcode) {
		return response{
			Code: "INVALID_PARAMETER",
			Msg:  "不正な商品識別子です。",
			Data: nil,
		}, nil
	}
	//查询评价
	var productId string
	err := global.GVA_DB.Table("products").
		Where("product_code = ?", productcode).
		Select("id").
		Scan(&productId).Error
	if err != nil || productId == "" {
		return response{
			Code: "NOT FOUND",
			Msg:  "商品が見つかりません。",
			Data: nil,
		}, nil
	}
	err = global.GVA_DB.Table("review_summaries").
		Select(`
		average_rating,
		review_count,
		rating_1_count,
		rating_2_count,
		rating_3_count,
		rating_4_count,
		rating_5_count
		`).
		Where("review_summaries.product_id = ?", productId).
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}
	//具体评论   COUNT(user_review_helpful_votes.review_id) AS helpful_count
	query := global.GVA_DB.Table("product_reviews").
		Select(`
		product_reviews.id,
		product_reviews.nickname,
		product_reviews.rating,
		product_reviews.comment,
		DATE_FORMAT(product_reviews.created_at,'%Y年%m月%d日') AS created_at_formatted,
		IFNULL(GROUP_CONCAT(DISTINCT review_images.image_url), '') AS image_urls,
		product_reviews.created_at
		`).
		Joins("LEFT JOIN review_summaries ON product_reviews.product_id=review_summaries.product_id").
		Joins("LEFT JOIN review_images ON product_reviews.id = review_images.review_id").
		//Joins("LEFT JOIN user_review_helpful_votes ON product_reviews.id = user_review_helpful_votes.review_id").
		Joins("LEFT JOIN products ON product_reviews.product_id = products.id").
		Group("product_reviews.id").
		Where("product_reviews.product_id=? AND product_reviews.status=?", productId, "approved")
		//筛选评分
	if rating > 0 {
		query = query.Where("product_reviews.rating=?", rating)
	}
	//给评论排序
	switch sort1 {
	case "newest":
		query = query.Order("product_reviews.created_at DESC")
	case "oldest":
		query = query.Order("product_reviews.created_at ASC")
	case "highest_rating":
		query = query.Order("product_reviews.Rating DESC").Order("product_reviews.created_at DESC")
	case "lowest_rating":
		query = query.Order("product_reviews.Rating ASC").Order("product_reviews.created_at DESC")
	}

	//分页
	offset := (page - 1) * limit
	//得到评论
	err = query.Offset(offset).
		Limit(limit).
		Scan(&reviews).Error
	if err != nil {
		return nil, err
	}
	//处理图片，变成数组，竖向输出,然后单独拿出helpful_count
	for i := range reviews {
		if reviews[i].RawImageUrls != "" {
			reviews[i].ImageUrls = strings.Split(reviews[i].RawImageUrls, ",")
		}
		var helpfulCount int64
		err = global.GVA_DB.Table("user_review_helpful_votes").
			Where("review_id", reviews[i].ID).
			Count(&helpfulCount).Error
		if err != nil {
			return nil, err
		}
		reviews[i].HelpfulCount = int(helpfulCount)
	}

	if sort1 == "most_helpful" {
		sort.Slice(reviews, func(i, j int) bool {
			if reviews[i].HelpfulCount == reviews[j].HelpfulCount {
				return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
			}
			return reviews[i].HelpfulCount > reviews[j].HelpfulCount
		})
	}

	//获取评论总数
	var totalCount int64
	tempCount := global.GVA_DB.Table("product_reviews").
		Where("product_id=? AND product_reviews.status=?", productId, "approved")

	// 如果有指定 rating，加入筛选条件
	if rating > 0 {
		tempCount = tempCount.Where("product_reviews.rating=?", rating)
	}

	err = tempCount.Count(&totalCount).Error
	if err != nil {
		return nil, err
	}
	//输出分页信息
	allPage := Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalCount:  totalCount,
		TotalPages:  int(math.Ceil(float64(totalCount) / float64(limit))),
	}
	//输出所有信息
	output := struct {
		Summary    Rating
		Reviews    []Review
		Pagination Pagination
	}{
		Summary:    summary,
		Reviews:    reviews,
		Pagination: allPage,
	}
	return output, nil
}

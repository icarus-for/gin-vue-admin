package service

import (
	"math"
	"regexp"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type QuestionInfo struct {
	UserId             int    `json:"user_id" gorm:"column:user_id"`
	ID                 int64  `json:"id" gorm:"column:id"`                              // 質問ID
	QuestionText       string `json:"question_text" gorm:"question_text"`               // 質問本文
	CreatedAtFormatted string `json:"created_at_formatted" gorm:"created_at_formatted"` // 表示用投稿日時
}

// AnswerInfo 回答情報 (answerer_type削除)
type AnswerInfo struct {
	ID                 int64  `json:"id" gorm:"column:id"`                              // 回答ID
	AnswererName       string `json:"answerer_name" gorm:"column:answerer_name"`        // 回答者表示名
	AnswerText         string `json:"answer_text" gorm:"column:answer_text"`            // 回答本文
	HelpfulCount       int    `json:"helpful_count" gorm:"column:helpful_count"`        // 参考になった数
	CreatedAtFormatted string `json:"created_at_formatted" gorm:"created_at_formatted"` // 表示用回答日時
}
type Pagination1 struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalCount  int64 `json:"total_count"`
	TotalPages  int   `json:"total_pages"`
}

type QuestionAndAnswer struct {
	Question QuestionInfo `json:"question" gorm:"column:question"`
	Answer   AnswerInfo   `json:"answer" gorm:"column:answer"`
}

type Getquestion struct{}

func (g *Getquestion) GetProductQueationInfo(productcode string, page int, limit int, sort2 string) (interface{}, error) {

	fakeCode := regexp.MustCompile(`^[0-9]+$`)
	if !fakeCode.MatchString(productcode) {
		return response{
			Code: "INVALID_PARAMETER",
			Msg:  "不正な商品識別子です。",
			Data: nil,
		}, nil
	}

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
	//分页
	offset := (page - 1) * limit

	var qaList []QuestionAndAnswer

	//查出问题和回答
	query := global.GVA_DB.Table("product_questions").
		Select(`
				product_questions.user_id AS user_id,
				product_questions.id AS question_id,
				product_questions.question_text,
				DATE_FORMAT(product_questions.created_at,'%Y年%m月%d日') AS question_created_at,
				question_answers.id AS answer_id,
				question_answers.answerer_name,
				question_answers.answer_text,
				DATE_FORMAT(question_answers.created_at,'%Y年%m月%d日') AS answer_created_at,
				Count(user_answer_helpful_votes.answer_id) AS helpful_count
				`).
		Joins("INNER JOIN question_answers ON question_answers.question_id = product_questions.id").
		Joins("LEFT JOIN user_answer_helpful_votes ON user_answer_helpful_votes.answer_id = question_answers.id").
		Where("product_questions.product_id=? AND product_questions.status=?", productId, "approved").
		Group("question_answers.id, product_questions.id, product_questions.question_text, question_answers.answerer_name, question_answers.answer_text, question_answers.created_at").
		Offset(offset).Limit(limit)

	switch sort2 {
	case "newest":
		query = query.Order("product_questions.created_at DESC")
	case "oldest":
		query = query.Order("product_questions.created_at ASC")
	case "most_helpful":
		query = query.Order("helpful_count DESC")
	default:
		query = query.Order("product_questions.created_at DESC")
	}

	type tempQAA struct {
		UserId            int   `gorm:"column:user_id"`
		QuestionId        int64 `gorm:"column:question_id"`
		QuestionText      string
		QuestionCreatedAt string `gorm:"column:question_created_at"`
		AnswerId          int64  `gorm:"column:answer_id"`
		AnswererName      string
		AnswerText        string
		AnswerCreatedAt   string `gorm:"column:answer_created_at"`
		HelpfulCount      int
	}
	var qaa []tempQAA
	err = query.Scan(&qaa).Error
	if err != nil {
		return nil, err
	}
	//把temp这个临时结构加到QuestionAndAnswer 实现一对一
	for _, v := range qaa {
		qaList = append(qaList, QuestionAndAnswer{
			Question: QuestionInfo{
				UserId:             v.UserId,
				ID:                 v.QuestionId,
				QuestionText:       v.QuestionText,
				CreatedAtFormatted: v.QuestionCreatedAt,
			},
			Answer: AnswerInfo{
				ID:                 v.AnswerId,
				AnswererName:       v.AnswererName,
				AnswerText:         v.AnswerText,
				HelpfulCount:       v.HelpfulCount,
				CreatedAtFormatted: v.AnswerCreatedAt,
			},
		})
	}

	//获得总问题数量
	var qTotal int64
	err = global.GVA_DB.Table("product_questions").
		Joins("INNER JOIN question_answers ON question_answers.question_id = product_questions.id").
		Where("product_questions.product_id=? AND product_questions.status=? AND question_answers.status=?", productId, "approved", "approved").
		Count(&qTotal).Error
	if err != nil {
		return nil, err
	}
	allPage := Pagination1{
		CurrentPage: page,
		Limit:       limit,
		TotalCount:  qTotal,
		TotalPages:  int(math.Ceil(float64(qTotal) / float64(limit))),
	}

	output := struct {
		QAlist     []QuestionAndAnswer
		Pagination Pagination1
	}{
		QAlist:     qaList,
		Pagination: allPage,
	}
	return output, nil
}

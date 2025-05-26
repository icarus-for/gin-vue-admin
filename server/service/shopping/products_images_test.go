package service

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 假设你这个结构体存在
// type response1 struct {
// 	Code string      `json:"code"`
// 	Msg  string      `json:"msg"`
// 	Data interface{} `json:"data"`
// }

func TestGetProductsImagesInfo_MockDB(t *testing.T) {
	// 初始化 mock 数据库
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 包装成 GORM DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// 注入全局变量
	global.GVA_DB = gormDB

	// 构造参数
	skuID := "sku-0001"

	// 1. 预期第一次查询：确认 sku_id 是否存在
	mock.ExpectQuery(regexp.QuoteMeta("SELECT sku_id FROM `sku_images` WHERE sku_id = ?")).
		WithArgs(skuID).
		WillReturnRows(sqlmock.NewRows([]string{"sku_id"}).AddRow(skuID))

	// 2. 预期第二次查询：查询图片信息
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			id,
			main_image_url,
			thumbnail_url,
			alt_text,
			sort_order
		FROM sku_images WHERE sku_id=?
	`)).
		WithArgs(skuID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "main_image_url", "thumbnail_url", "alt_text", "sort_order",
		}).AddRow(1, "http://main.jpg", "http://thumb.jpg", "alt text", 1))

	// 调用方法
	g := GetImages{}
	result, err := g.GetProductsImagesInfo(skuID)
	assert.NoError(t, err)

	// 验证结果是否为 []SkuImageInfo 类型
	data, ok := result.([]SkuImageInfo)
	assert.True(t, ok)
	assert.Len(t, data, 1)
	assert.Equal(t, "http://main.jpg", data[0].MainImageURL)
}

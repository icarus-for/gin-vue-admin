package service

import (
	"fmt"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	// ⚠️ 修改为你自己的数据库配置
	dsn := "root:123@tcp(172.23.240.1:3306)/gva?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to test DB: %v", err))
	}
	global.GVA_DB = db
}
func TestGetProductsImagesInfo(t *testing.T) {
	initDB()
	g := GetImages{}

	// ✅ 合法 SKU ID（你需要确认这个 ID 在数据库存在）
	skuID := "sku-0001"

	result, err := g.GetProductsImagesInfo(skuID)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Errorf("Expected result, got nil")
	}
}

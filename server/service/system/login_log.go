package system

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
)

// RecordLoginLog 异步记录用户登录信息
func RecordLoginLog(userID uint, username string) {
	log := system.LoginLog{
		UserID:   userID,
		Username: username,
		LoginAt:  time.Now(),
	}
	if err := global.GVA_DB.Create(&log).Error; err != nil {
		global.GVA_LOG.Error("记录登录日志失败", zap.Error(err))
	}
}

func GetUserLoginHistoryById(userid *uint, start *time.Time, end *time.Time) ([]system.LoginLog, error) {
	var htr []system.LoginLog
	db := global.GVA_DB.Model(&system.LoginLog{})
	if userid != nil {
		db = db.Where("user_id=?", *userid)
	}
	if start != nil {
		db = db.Where("login_at>=?", *start)
	}
	if end != nil {
		db = db.Where("login_at<=?", *end)
	}
	err := db.Order("login_at DESC").Find(&htr).Error
	return htr, err
}

// func DelUserLoginHistoryById(userid *uint, start *time.Time, end *time.Time) ([]system.LoginLog, error) {
// 	tx := global.GVA_DB.Model(&system.LoginLog{})
// 	if err := tx.Where("user_id = ?", *userid).Delete(&system.LoginLog{}).Error; err != nil {
// 		return []system.LoginLog,err
// 	}
// 	if err := tx.Where("login_at>=?", *start).Delete(&system.LoginLog{}).Error; err != nil {
// 		return err
// 	}
// 	if err := tx.Where("login_at>=?", *start).Delete(&system.LoginLog{}).Error; err != nil {
// 		return err
// 	}

// }

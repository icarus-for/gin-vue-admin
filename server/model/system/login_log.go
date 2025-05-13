package system

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type LoginLog struct {
	global.GVA_MODEL
	UserID   uint      `json:"user_id" form:userId` // 用户ID
	Username string    `json:"username"`            // 用户名
	LoginAt  time.Time `json:"login_at"`            // 登录时间
}

func (LoginLog) TableName() string {
	return "gva.login_logs"
}

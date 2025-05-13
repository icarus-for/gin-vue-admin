package system

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	modelSystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/gin-gonic/gin"
)

func RecordLoginLog(userID uint, username string) {
	log := modelSystem.LoginLog{
		UserID:   userID,
		Username: username,
		LoginAt:  time.Now(),
	}
	global.GVA_DB.Create(&log)
}

type LoginLogSearchRequest struct {
	UserID uint   `form:"userId"` // 用户 ID
	Start  string `form:"start"`  // 开始时间（格式：2006-01-02）
	End    string `form:"end"`
}

// @Tags LoginLog
// @Summary 条件查询用户登录历史
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param userId query int false "用户ID"
// @Param start query string false "开始时间（格式：2006-01-02）"
// @Param end query string false "结束时间（格式：2006-01-02）"
// @Success 200 {object} response.Response{data=[]modelSystem.LoginLog,msg=string} "查询成功"
// @Failure 400 {object} response.Response "请求失败"
// @Router /user/getUserLoginHistoryById [get]
func GetUserLoginHistorycondition(c *gin.Context) {
	var req LoginLogSearchRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	starturl := c.Query("start")
	endurl := c.Query("end")
	const timeLayOut = "2006-01-02"
	var startTime, endTime time.Time
	if starturl != "" {
		st, err := time.Parse(timeLayOut, starturl)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		startTime = st
	} else {
		startTime = time.Time{}
	}
	if endurl != "" {
		et, err := time.Parse(timeLayOut, endurl)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		endTime = et
	} else {
		endTime = time.Now()
	}
	var userIDPtr *uint = nil
	if req.UserID != 0 {
		userIDPtr = &req.UserID
	}
	result, err := system.GetUserLoginHistoryById(userIDPtr, &startTime, &endTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

package main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/pkg/redis"
	"github.com/gin-gonic/gin"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// 这部分 @Tag 设置用于排序, 需要排序的接口请按照下面的格式添加
// swag init 对 @Tag 只会从入口文件解析, 默认 main.go
// 也可通过 --generalInfo flag 指定其他文件
// @Tag.Name        Base
// @Tag.Name        SysUser
// @Tag.Description 用户

// @title                       Gin-Vue-Admin Swagger API接口文档
// @version                     v2.8.1
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	if global.GVA_DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}
	core.RunWindowsServer()
	redis.InitRedis()

	// 初始化 Gin 路由
	r := gin.Default()

	// 设置路由
	r.GET("/", func(c *gin.Context) {
		// 设置 Redis 缓存中的 key
		err := redis.RedisClient.Set(redis.Ctx, "name", "RedisDemo", 0).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": "Redis 设置失败"})
			return
		}

		// 获取 Redis 缓存中的 key
		val, err := redis.RedisClient.Get(redis.Ctx, "name").Result()
		if err != nil {
			c.JSON(500, gin.H{"error": "Redis 获取失败"})
			return
		}

		// 返回 Redis 中的值
		c.JSON(200, gin.H{"key": val})
	})

	// 启动 Gin 服务器
	r.Run(":8080") // 监听端口 8080
}

package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

// 初始化 Redis 连接
func InitRedis() {
	// 创建一个 Redis 客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址，默认是 localhost:6379
		Password: "",               // 如果没有设置密码，留空即可
		DB:       0,                // 使用 Redis 的第 0 个数据库
	})

	// 测试 Redis 是否连接成功
	err := RedisClient.Ping(Ctx).Err()
	if err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	} else {
		fmt.Println("Redis 连接成功！")
	}
}

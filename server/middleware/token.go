package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// token认证
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		// 如果 Authorization 头中带有 Bearer 前缀，去掉它
		fmt.Println("Received Token: ", token)

		if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
			token = token[7:] // 去掉前面的 Bearer 空格
		}
		if token == "" {
			// 默认错误处理
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未登录或Token缺失", // 默认错误信息
			})
			return
		}
		c.Next()
	}
}

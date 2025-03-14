package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/youngking/gin-blog/pkg/e"
	utills "github.com/youngking/gin-blog/pkg/util"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS

		// 获取请求中的token
		token := c.Query("token")
		// 判断是否合法
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			// 解析token
			claims, err := utills.ParseToken(token)
			// 判读是否合法
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt { // 判读是否过期
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		// 不合法
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": make(map[string]string),
			})
			// 终止后续处理
			c.Abort()
			return
		}
		// 继续后续路由处理
		c.Next()
	}
}

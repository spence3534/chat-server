package middleware

import (
	"chat-server/models/common/response"
	"chat-server/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-Token")
		fmt.Println(token)

		if token == "" {
			response.FailWithMessage("用户未登录", c)
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.FailWithMessage("用户信息已过期", c)
				c.Abort()
				return
			}
			response.FailWithMessage(err.Error(), c)
			fmt.Println(claims)
			c.Abort()
		}

		c.Next()
	}
}

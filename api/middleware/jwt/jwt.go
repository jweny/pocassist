package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/util"
	"net/http"
	"time"
)

// middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		var errStr string

		code = msg.SuccessCode
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			// 非登录状态
			code = msg.ErrCode
			errStr = "请登录后操作"
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
			//	token 校验不通过
				code = msg.ErrCode
				errStr = "身份验证失败，请重新登录"
			} else if time.Now().Unix() > claims.ExpiresAt {
			//	token 已过期
				code = msg.ErrCode
				errStr = "身份信息已过期，请重新登录"
			}
		}

		if code != msg.SuccessCode {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code" : code,
				"msg" : errStr,
				"data" : data,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
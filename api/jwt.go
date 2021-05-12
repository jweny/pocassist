package api

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var jwtSecret = []byte("pocassist")

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	// 过期时间
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims {
			ExpiresAt : expireTime.Unix(),
			Issuer : "pocassist",
		},
	}
	// sha256
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(authHeader string) (*Claims, error) {
	// token 前缀为 "JWT"
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "JWT") {
		return nil, errors.New("请求头中auth格式不正确")
	}
	token := parts[1]
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		var errStr string

		code = SuccessCode

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			// 非登录状态
			code = ErrCode
			errStr = "请登录后操作"
		} else {
			claims, err := ParseToken(token)
			if err != nil {
			//	token 校验不通过
				code = ErrCode
				errStr = "身份验证失败，请重新登录"
			} else if time.Now().Unix() > claims.ExpiresAt {
			//	token 已过期
				code = ErrCode
				errStr = "身份信息已过期，请重新登录"
			}
		}

		if code != SuccessCode {
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
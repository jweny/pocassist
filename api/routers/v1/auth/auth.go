package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/util"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type ResetPwd struct {
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

func Login(c *gin.Context) {
	login := auth{}
	err := c.BindJSON(&login)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}

	username := login.Username
	password := login.Password
	data := make(map[string]interface{})

	isExist, auth := db.CheckAuth(username, password)
	if isExist {
		token, err := util.GenerateToken(username, password)
		if err != nil {
			c.JSON(msg.ErrResp("token生成失败"))
			return
		} else {
			data["token"] = token
			data["nickname"] = username
			data["uid"] = auth.Id
			c.JSON(msg.SuccessResp(data))
			return
		}
	} else {
		c.JSON(msg.ErrResp("用户名/密码错误"))
		return
	}
}

func Reset(c *gin.Context) {
	resetPwd := ResetPwd{}
	err := c.BindJSON(&resetPwd)
	if err != nil {
		c.JSON(msg.ErrResp("参数校验不通过"))
		return
	}
	token := c.Request.Header.Get("Authorization")
	claims, err := util.ParseToken(token)
	if err != nil || claims == nil {
		c.JSON(msg.ErrResp("token校验失败"))
		return
	}
	isExist, auth := db.CheckAuth(claims.Username, resetPwd.Password)
	if isExist {
		// 修改密码
		db.ResetPassword(auth.Id,resetPwd.NewPassword)
		c.JSON(msg.SuccessResp("密码更新成功"))
	} else {
		c.JSON(msg.ErrResp("密码错误"))
	}
}

func Self(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	claims, err := util.ParseToken(token)
	if err != nil || claims == nil {
		c.JSON(msg.ErrResp("token校验失败"))
		return
	}
	userInfo := make(map[string]string)
	userInfo["name"] = claims.Username
	c.JSON(msg.SuccessResp(userInfo))
	return
}

func Logout(c *gin.Context) {
	// 后端伪登出 todo:优化jwt
	c.JSON(msg.SuccessResp("登出成功"))
	return
}
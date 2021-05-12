package api

import (
	"github.com/gin-gonic/gin"
	"pocassist/database"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type ResetPwd struct {
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

func GetAuth(c *gin.Context) {
	login := auth{}
	err := c.BindJSON(&login)
	if err != nil {
		c.JSON(ErrResp("参数校验不通过"))
		return
	}
	username := login.Username
	password := login.Password
	data := make(map[string]interface{})

	isExist, auth := database.CheckAuth(username, password)
	if isExist {
		token, err := GenerateToken(username, password)
		if err != nil {
			c.JSON(ErrResp("token生成失败"))
		} else {
			data["token"] = token
			data["nickname"] = username
			data["uid"] = auth.Id
			c.JSON(SuccessResp(data))
		}
	} else {
		c.JSON(ErrResp("用户名/密码错误"))
	}

}

func SelfResetPassword(c *gin.Context) {
	resetPwd := ResetPwd{}
	err := c.BindJSON(&resetPwd)
	if err != nil {
		c.JSON(ErrResp("参数校验不通过"))
		return
	}
	token := c.Request.Header.Get("Authorization")
	claims, err := ParseToken(token)
	if err != nil || claims == nil {
		c.JSON(ErrResp("token校验失败"))
		return
	}
	isExist, auth := database.CheckAuth(claims.Username, resetPwd.Password)
	if isExist {
		// 修改密码
		database.ResetPassword(auth.Id,resetPwd.NewPassword)
		c.JSON(SuccessResp("密码更新成功"))
	} else {
		c.JSON(ErrResp("密码错误"))
	}
}

func SelfGetInfo(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	claims, err := ParseToken(token)
	if err != nil || claims == nil {
		c.JSON(ErrResp("token校验失败"))
		return
	}
	userInfo := make(map[string]string)
	userInfo["name"] = claims.Username
	c.JSON(SuccessResp(userInfo))
	return
}

func SelfLogout(c *gin.Context) {
	// 后端伪登出 todo:优化jwt
	c.JSON(SuccessResp("登出成功"))
	return
}
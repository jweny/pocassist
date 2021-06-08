package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/msg"
	"github.com/jweny/pocassist/pkg/db"
	"github.com/jweny/pocassist/pkg/util"
)

type Auth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPwd struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"newpassword" binding:"required"`
}

// @Summary Login
// @Tags User
// @Description 登录
// @accept json
// @Produce  json
// @Param auth body Auth true "用户/密码"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	login := Auth{}
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(msg.ErrResp("用户名密码不可为空"))
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

// @Summary Reset Password
// @Tags User
// @Description 重置密码
// @accept json
// @Produce  json
// @Param resetpwd body ResetPwd true "旧/新密码"
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/self/resetpwd/ [post]
func Reset(c *gin.Context) {
	resetPwd := ResetPwd{}
	err := c.ShouldBindJSON(&resetPwd)
	if err != nil {
		c.JSON(msg.ErrResp("原密码、新密码不可为空"))
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

// @Summary Self
// @Tags User
// @Description 获取个人信息
// @Produce  json
// @Security token
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/user/info [get]
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

// @Summary Logout
// @Tags User
// @Description 登出
// @Produce  json
// @Security token
// @Success 200 {object} msg.Response
// @Failure 200 {object} msg.Response
// @Router /api/v1/user/logout [get]
func Logout(c *gin.Context) {
	// 后端伪登出 todo:优化jwt
	c.JSON(msg.SuccessResp("登出成功"))
	return
}
package db

import "gorm.io/gorm"

// auths 表

type Auth struct {
	gorm.Model
	Id       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CheckAuth checks if authentication information exists
func CheckAuth(username, password string) (bool, *Auth) {
	var auth Auth
	GlobalDB.Select("id").Where(Auth{Username : username, Password : password}).First(&auth)
	if auth.Id > 0 {
		return true, &auth
	}
	return false, nil
}

// ResetPassword
func ResetPassword(id int, newpassword string) {
	var auth Auth
	GlobalDB.Model(auth).Where(Auth{Id : id}).Update("password",newpassword)
}
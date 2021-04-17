package model

import (
	"errors"
	"go-gateway/inc"

	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Name   string `json:"name" gorm:"type:varchar(100);not null;default:'';comment:用户名"`
	Pwd    string `json:"pwd" gorm:"type:varchar(100);not null;default:'';comment:密码"`
	Email  string `json:"email" gorm:"type:varchar(100);not null;default:'';comment:邮箱"`
	Status int    `json:"status" gorm:"type:tinyint(2);not null;default:1;comment:状态: 1 正常,2 禁用"`
}

type AdminLoginInfo struct {
	Name string `json:"name" form:"name" binding:"required"`
	Pwd  string `json:"pwd" form:"pwd" binding:"required"`
}

func GetFirstAdmin(query interface{}, args ...interface{}) (admin Admin, row int64) {
	row = GetFirst(&admin, query, args...)
	return
}

//登陆操作
func DoLogin(a AdminLoginInfo) (admin Admin, err error) {

	row := GetFirst(&admin, "name =? and pwd = ?", a.Name, inc.MD5(a.Pwd))
	if row < 1 {
		err = errors.New("用户名密码错误")
		return
	}

	if admin.Status != 1 {
		err = errors.New("账号已被禁用")
		return
	}
	return
}

package build

import (
	"go-gateway/debug"
	"go-gateway/inc"
	"go-gateway/model"
)

func InitDb(pwd string) {

	if pwd == "" {
		panic("请输入admin账户初始化密码")
	}

	//创建admin表
	db := model.Db.Migrator()

	if ok := db.HasTable(&model.Admin{}); ok {
		panic("系统表已存在，无需初始化")
	}

	model.Db.Set("gorm:table_options", "ENGINE=InnoDB comment '管理账户表'").AutoMigrate(&model.Admin{})
	debug.DebugPrint("[success] 管理账户表初始化完成")

	//插入 e10adc3949ba59abbe56e057f20f883e
	admin := model.Admin{Name: "admin", Pwd: inc.MD5(pwd)}
	model.Db.Create(&admin)

}

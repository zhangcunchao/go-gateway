package model

import (
	"go-gateway/debug"
	"go-gateway/inc"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

//基础model
type Model struct {
	Id uint `gorm:"primaryKey" json:"id"`
}

func init() {
	defer func() {
		if r := recover(); r != nil {
			debug.DebugPrint("捕获到dao层错误：%s\n", r)
			os.Exit(2)
		}
	}()
	//初始化数据库配置
	dConf, _ := inc.Cfg.GetSection("db")

	_dsn := dConf["UserName"] + ":" + dConf["PassWord"] + "@tcp(" + dConf["Host"] + ":" + dConf["Port"] + ")/" + dConf["Database"] + "?charset=utf8mb4"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)
	var err error

	Db, err = gorm.Open(mysql.Open(_dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   dConf["TablePrefix"] + "_", // 表名前缀，`Article` 的表名应该是 `gg_articles`
			SingularTable: true,                       // 使用单数表名，启用该选项，此时，`Article` 的表名应该是 `gg_article`
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	sqlDb, err := Db.DB()

	if err != nil {
		panic(err)
	}
	//连接池设置
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxIdleConns(inc.Cfg.MustInt("db", "MaxIdleConns", 10))

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDb.SetMaxOpenConns(inc.Cfg.MustInt("db", "MaxOpenConns", 100))
}

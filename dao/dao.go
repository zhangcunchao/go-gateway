package dao

import (
	"database/sql"
	"go-gateway/debug"
	"go-gateway/inc"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DbWorker struct {
	Dsn string
	Db  *sql.DB
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

	dbw := DbWorker{Dsn: _dsn}

	dbtemp, err := sql.Open("mysql", dbw.Dsn)
	dbw.Db = dbtemp
	if err != nil {
		panic(err)
	}
	err = dbw.Db.Ping()
	if err != nil {
		panic(err)
	}
	defer dbw.Db.Close()
}

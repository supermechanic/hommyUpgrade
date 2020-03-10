package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//GormDb 数据库实例
var GormDb *gorm.DB

var (
	mSQLUser string
	mSQLPass string
	mSQLHost string
	mSQLDb   string
)

var err error
var lock sync.Mutex

func init() {
	mSQLUser = Config.Mysql.Base.Username
	mSQLPass = Config.Mysql.Base.Password
	mSQLHost = Config.Mysql.Base.Address +
		":" + Config.Mysql.Base.Port
	mSQLDb = Config.Mysql.DatabaseName
	newGorm()
}

func newGorm() (*gorm.DB, error) {
	args := []string{mSQLUser, ":", mSQLPass, "@", "tcp(", mSQLHost, ")/", mSQLDb, "?charset=utf8&parseTime=True&loc=Local"}
	argBuf := bytes.Buffer{}
	for _, arg := range args {
		argBuf.WriteString(arg)
	}
	argsStr := argBuf.String()
	if GormDb == nil {
		lock.Lock()
		defer lock.Unlock()
		if GormDb == nil {
			GormDb, err = gorm.Open("mysql", argsStr)
			if err != nil {
				GormDb = nil
				log.Println("mysql数据库连接失败:", err)
			}
			fmt.Println("gormdb get success")
		}
		GormDb.DB().SetMaxIdleConns(0)
		GormDb.DB().SetMaxOpenConns(100)
		GormDb.LogMode(true)
	}

	return GormDb, err
}

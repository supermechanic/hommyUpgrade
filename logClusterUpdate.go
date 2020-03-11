package main

import (
	"github.com/golang/glog"
)
/**
数据库检索所有的抓拍日志，



*/
type logInfo struct {
	ID        uint32
	ImagePath string
	ClusterID int32
	Time      int64
}

func logTableNames() error {
	rows, err := GormDb.Raw("select table_name from information_schema.tables where table_schema='xgface' and table_name like 'log_20';").rows()
	if err != nil {
		glog.Error.Println(err)
		return err
	}
	defer rows.Close()
}

func getAllImageWithTime() {

}

func sendToRedis() {

}

func updateLogClusterID() {

}

func main() {
	flag.Parse()
	defer glog.Flush()
}

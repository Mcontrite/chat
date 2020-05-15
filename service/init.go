package service

import (
	"chat/model"
	"errors"
	"fmt"
	"log"

	"github.com/go-xorm/xorm"
)

const (
	driveName = "mysql"
	dsName    = "root:123456@(127.0.0.1:3306)/chat?charset=utf8"
	showSQL   = true
	maxCon    = 10
	NONERROR  = "noerror" //没有错误
)

var DbEngin *xorm.Engine

//初始化数据库
func init() {
	err := errors.New(NONERROR)
	DbEngin, err = xorm.NewEngine(driveName, dsName)
	if nil != err && NONERROR != err.Error() {
		log.Fatal(err.Error())
	}
	DbEngin.ShowSQL(showSQL)                                                 //是否显示SQL语句
	DbEngin.SetMaxOpenConns(maxCon)                                          //最大打开的连接数
	DbEngin.Sync2(new(model.User), new(model.Contact), new(model.Community)) //自动建表
	fmt.Println("init data base ok")
}

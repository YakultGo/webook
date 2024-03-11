package ioc

import (
	"basic-go/webook/config"
	"basic-go/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 初始化过程中出现错误，直接退出
		panic(err)
	}
	// 建表
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

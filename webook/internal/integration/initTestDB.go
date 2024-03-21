package integration

import (
	"basic-go/webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initTestGORM() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"), nil)
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
func initRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
	return rdb
}

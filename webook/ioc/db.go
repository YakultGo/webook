package ioc

import (
	"basic-go/webook/internal/repository/dao"
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("mysql", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), nil)
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

func InitMongoDB() *mongo.Database {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("mongo", &cfg)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DSN))
	if err != nil {
		panic(err)
	}
	return client.Database("webook")
}

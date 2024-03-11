package ioc

import (
	"basic-go/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return rdb
}

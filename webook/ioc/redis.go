package ioc

import (
	"basic-go/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	return rdb
}

func NewRateLimiter(rdb redis.Cmdable) ratelimit.Limiter {
	return ratelimit.NewRedisSlidingWindowLimiter(rdb, time.Second, 100)
}

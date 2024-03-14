package ioc

import (
	"basic-go/webook/config"
	"basic-go/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return rdb
}

func NewRateLimiter(rdb redis.Cmdable) ratelimit.Limiter {
	return ratelimit.NewRedisSlidingWindowLimiter(rdb, time.Second, 100)
}

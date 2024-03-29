package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证码错误次数太多")
	ErrUnknown                = errors.New("未知错误")
)

// 编译器在编译时，会将lua/set_code.lua文件的内容嵌入到luaSetCode变量中
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

var _ CodeCache = (*RedisCodeCache)(nil)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(c redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: c,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context,
	biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		zap.L().Warn("短信发送太频繁",
			zap.String("biz", biz),
			zap.String("phone", phone),
		)
		return ErrCodeSendTooMany
	default:
		return errors.New("系统错误")
	}
}
func (c *RedisCodeCache) Verify(ctx context.Context,
	biz, phone, code string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 频繁出现可能是被恶意攻击
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	default:
		return false, ErrUnknown
	}
}
func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

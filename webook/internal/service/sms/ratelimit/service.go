package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/ratelimit"
	"context"
	"fmt"
)

var (
	errLimited = fmt.Errorf("触发限流")
)

type RateLimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func (s *RateLimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:alibaba")
	if err != nil {
		// 系统错误
		return fmt.Errorf("判断限流服务出现问题 %w", err)
	}
	if limited {
		return errLimited
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	return err
}

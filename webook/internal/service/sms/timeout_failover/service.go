package timeout_failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
)

type FailOverSMSService struct {
	svcs []sms.Service
}

func NewFailOverSMSService(svcs []sms.Service) sms.Service {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

func (f FailOverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		// 发送成功
		if err == nil {
			return nil
		}
	}
	return errors.New("所有服务都发送失败")
}

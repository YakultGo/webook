package retryable

import (
	"basic-go/webook/internal/service/sms"
	"context"
)

type Service struct {
	svc      sms.Service
	retryCnt int
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	if err != nil {
		for i := 0; i < s.retryCnt; i++ {
			err = s.svc.Send(ctx, tpl, args, numbers...)
			if err == nil {
				return nil
			}
		}
	}
	return err
}

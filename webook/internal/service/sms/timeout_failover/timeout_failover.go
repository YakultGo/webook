package timeout_failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"sync/atomic"
)

type TimeoutFailOverSMSService struct {
	svcs []sms.Service
	idx  int32
	// 连续超时的个数
	cnt int32
	// 阈值
	threshold int32
}

func NewTimeoutFailOverSMSService() sms.Service {
	return &TimeoutFailOverSMSService{}
}
func (t *TimeoutFailOverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		} else {
			idx = atomic.LoadInt32(&t.idx)
		}
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt32(&t.cnt, 1)
	case err == nil:
		atomic.StoreInt32(&t.cnt, 0)
	default:
		return err
	}
	return nil
}

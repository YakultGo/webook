package memory

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
)

var _ sms.Service = (*MemoryService)(nil)

type MemoryService struct {
}

func NewService() *MemoryService {
	return &MemoryService{}
}
func (s MemoryService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Printf("[%s : %s]", args[0], args[1])
	return nil
}

package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args[0], args[1])
	return nil
}

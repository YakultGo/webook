package service

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

const (
	codeTemplateId = "SMS_465341536"
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

// Send 发送验证码
/*
	biz 业务类型
	code 验证码
	phone 手机号
*/
func (svc *CodeService) Send(ctx context.Context,
	biz, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTemplateId, []string{code}, phone)
	return err
}

func (svc *CodeService) Verify(ctx context.Context,
	biz, phone, code string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, code)
}

func (svc *CodeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

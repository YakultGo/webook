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

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)
var _ CodeService = (*CodeServiceStruct)(nil)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
type CodeServiceStruct struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &CodeServiceStruct{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发送验证码
/*
	biz 业务类型
	code 验证码
	phone 手机号
*/
func (svc *CodeServiceStruct) Send(ctx context.Context,
	biz, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTemplateId, []string{"code", code}, phone)
	return err
}

func (svc *CodeServiceStruct) Verify(ctx context.Context,
	biz, phone, code string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, code)
}

func (svc *CodeServiceStruct) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

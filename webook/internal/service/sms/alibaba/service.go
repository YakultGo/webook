package alibaba

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

var _ sms.Service = (*AliService)(nil)

type AliService struct {
	signName string
	client   *dysmsapi20170525.Client
}

func NewService(client *dysmsapi20170525.Client, signName string) *AliService {
	return &AliService{
		signName: signName,
		client:   client,
	}
}

func (s AliService) Send(ctx context.Context, tpl string, args []string, number ...string) error {
	req := dysmsapi20170525.SendSmsRequest{}
	req.SignName = tea.String(s.signName)
	req.PhoneNumbers = tea.String(number[0])
	req.TemplateCode = tea.String(tpl)
	req.TemplateParam = tea.String(fmt.Sprintf(`{"%s":"%s"}`, args[0], args[1]))
	_, err := s.client.SendSms(&req)
	if err != nil {
		return err
	}
	return nil
}

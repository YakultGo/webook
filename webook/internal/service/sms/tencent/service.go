package tencent

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
	Tsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

var _ sms.Service = (*TenService)(nil)

type TenService struct {
	appId    *string
	signName *string
	client   *Tsms.Client
}

func NewService(client *Tsms.Client, appId, signName string) *TenService {
	return &TenService{
		appId:    &appId,
		signName: &signName,
		client:   client,
	}
}
func (s *TenService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	req := Tsms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = &tpl
	req.PhoneNumberSet = s.sliceToPtrStr(numbers)
	req.TemplateParamSet = s.sliceToPtrStr(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *status.Code != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *TenService) sliceToPtrStr(str []string) []*string {
	var res []*string
	for _, v := range str {
		res = append(res, &v)
	}
	return res
}

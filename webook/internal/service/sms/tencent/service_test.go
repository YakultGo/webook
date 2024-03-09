package tencent

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"testing"
)

func TestService_Send(t *testing.T) {
	secretId, ok := os.LookupEnv("Tencent_sms_secret_id")
	if !ok {
		t.Fatal()
	}
	secretKey, ok := os.LookupEnv("Tencent_sms_secret_key")
	if !ok {
		t.Fatal()
	}
	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-guangzhou", profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}
	s := NewService(c, "1400491234", "腾讯云")
	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:    "发送验证码",
			tplId:   "123456",
			params:  []string{"123456"},
			numbers: []string{"12345678901"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

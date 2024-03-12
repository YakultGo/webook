package integration

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				val := rdb.GetDel(ctx, "phone_code:login:10010")
				cancel()
				assert.True(t, len(val.Val()) == 6)
			},
			reqBody: `
			{
				"phone": "10010"
			}
`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:10010", "123456",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				val := rdb.GetDel(ctx, "phone_code:login:10010")
				cancel()
				assert.Equal(t, "123456", val.Val())
			},
			reqBody: `
			{
				"phone": "10010"
			}
`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送太频繁，请稍后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:10010", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				val := rdb.GetDel(ctx, "phone_code:login:10010")
				cancel()
				assert.True(t, len(val.Val()) == 6)
			},
			reqBody: `
			{
				"phone": "10010"
			}
`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手机号码为空",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
			{
				"phone": ""
			}
`,
			wantCode: 200,
			wantBody: web.Result{
				Code: 4,
				Msg:  "手机号有误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			var res web.Result
			err = json.Unmarshal(resp.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			assert.Equal(t, tc.wantBody, res)
			tc.after(t)
		})
	}
}

package web

import (
	"basic-go/webook/internal/service"
	svcmocks "basic-go/webook/internal/service/mocks"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return userSvc
			},
			reqBody: `
	{
		"email": "Yakult@qq.com",
		"password": "TZX5at4nbVHF0",
		"confirmPassword": "TZX5at4nbVHF0"
	}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
	{
		"email": "Yakult@qq.com",
		"password": "TZX5at4nbVHF0"
	`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
	{
		"email": "Yakult",
		"password": "TZX5at4nbVHF0",
		"confirmPassword": "TZX5at4nbVHF0"
	}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式错误",
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
	{
		"email": "Yakult@qq.com",
		"password": "TZX5at4nbVHF0",
		"confirmPassword": "123TZX5at4nbVHF0"
	}`,
			wantCode: http.StatusOK,
			wantBody: "两次密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody: `
	{
		"email": "Yakult@qq.com",
		"password": "123",
		"confirmPassword": "123"
	}`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，且包含数字和字母",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate)
				return userSvc
			},
			reqBody: `
	{
		"email": "Yakult@qq.com",
		"password": "TZX5at4nbVHF0",
		"confirmPassword": "TZX5at4nbVHF0"
	}`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("系统异常"))
				return userSvc
			},
			reqBody: `
	{
		"email": "Yakult@qq.com",
		"password": "TZX5at4nbVHF0",
		"confirmPassword": "TZX5at4nbVHF0"
	}`,
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.mock(control), nil)
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/users/signup",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	repomocks "basic-go/webook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestUserServiceStruct_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		input    domain.User
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "Yakult@qq.com").
					Return(domain.User{
						Email:      "Yakult@qq.com",
						Password:   "$2a$10$eGik8oiZS7oFvsBzMaXh8O5/NxYTW2xIuaeTvcYehVZ7U7NI.VjoK",
						CreateTime: now,
					}, nil)
				return repo
			},
			input: domain.User{
				Email:    "Yakult@qq.com",
				Password: "TZX5at4nbVHF0",
			},
			wantUser: domain.User{
				Email:      "Yakult@qq.com",
				Password:   "$2a$10$eGik8oiZS7oFvsBzMaXh8O5/NxYTW2xIuaeTvcYehVZ7U7NI.VjoK",
				CreateTime: now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "Yakult@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			input: domain.User{
				Email:    "Yakult@qq.com",
				Password: "TZX5at4nbVHF0",
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "Yakult@qq.com").
					Return(domain.User{}, errors.New("mock db错误"))
				return repo
			},
			input: domain.User{
				Email:    "Yakult@qq.com",
				Password: "TZX5at4nbVHF0",
			},
			wantUser: domain.User{},
			wantErr:  errors.New("mock db错误"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "Yakult@qq.com").
					Return(domain.User{
						Email:      "Yakult@qq.com",
						Password:   "$2a$10$eGik8oiZS7oFvsBzMaXh8O5/NxYTW2xIuaeTvcYehVZ7U7NI.VjoK",
						CreateTime: now,
					}, ErrInvalidUserOrPassword)
				return repo
			},
			input: domain.User{
				Email:    "Yakult@qq.com",
				Password: "TZX5at4nbVHF0",
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.input)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

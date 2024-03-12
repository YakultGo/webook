package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	cachemocks "basic-go/webook/internal/repository/cache/mocks"
	"basic-go/webook/internal/repository/dao"
	daomocks "basic-go/webook/internal/repository/dao/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		ctx  context.Context
		id   int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，但是查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDAO(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id:    123,
						Email: "123@qq.com",
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:    123,
					Email: "123@qq.com",
				}).Return(nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:    123,
				Email: "123@qq.com",
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDAO(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:    123,
						Email: "123@qq.com",
					}, nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:    123,
				Email: "123@qq.com",
			},
			wantErr: nil,
		},
		{
			name: "数据库查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDAO(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("数据库查询失败"))
				return d, c
			},
			ctx:      context.Background(),
			id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("数据库查询失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

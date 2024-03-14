//go:build wireinject

package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	myJwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.NewRateLimiter,
		dao.NewUserDAO,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		web.NewOAuth2WechatHandler,
		myJwt.NewRedisJWTHandler,
		ioc.InitSMSService,
		web.NewUserHandler,
		ioc.InitGin,
		ioc.InitMiddlewares,
		ioc.InitOAuth2WechatService,
	)
	return new(gin.Engine)
}

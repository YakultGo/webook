package ioc

import (
	"basic-go/webook/internal/web"
	myJwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/middlewares/ratelimit"
	pkgratelimit "basic-go/webook/pkg/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitGin(handlerFunc []gin.HandlerFunc, handler *web.UserHandler,
	oauth2WeChatHandler *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(handlerFunc...)
	handler.RegisterRoutes(server)
	oauth2WeChatHandler.RegisterRoutes(server)
	return server
}

func InitMiddlewares(limiter pkgratelimit.Limiter, hdl myJwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ratelimit.NewBuilder(limiter).Build(),
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuilder(hdl).
			IgnorePath("/users/signup").
			IgnorePath("/users/login").
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			IgnorePath("/oauth2/wechat/author").
			IgnorePath("/oauth2/wechat/callback").
			IgnorePath("/users/refresh_token").
			Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	})

}

package ioc

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitGin(handlerFunc []gin.HandlerFunc, handler *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(handlerFunc...)
	handler.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
		corsHdl(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePath("/users/signup").
			IgnorePath("/users/login").
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	})

}

package main

import (
	"basic-go/webook/config"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/sms/memory"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	u := initUser(db, rdb)
	server := initWebServer()
	u.RegisterRoutes(server)

	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
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
	}))

	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("vbCzjQ3aud0CFVjDFM91dHARoYbhHo5j"),
	//	[]byte("EUCzqATwmrF00y08rXwoQ4nTBPh4Xnxc"))
	//store, err := redis.NewStore(16, "tcp", "localhost:16379", "",
	//	[]byte("vbCzjQ3aud0CFVjDFM91dHARoYbhHo5j"), []byte("EUCzqATwmrF00y08rXwoQ4nTBPh4Xnxc"))
	//if err != nil {
	//	panic(err)
	//}

	//server.Use(sessions.Sessions("mysession", store))

	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePath("/users/signup").
	//	IgnorePath("/users/login").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePath("/users/signup").
		IgnorePath("/users/login").
		IgnorePath("/users/login_sms/code/send").Build())
	return server
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	c := cache.NewUserCache(rdb, time.Minute*10)
	repo := repository.NewUserRepository(ud, c)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 初始化过程中出现错误，直接退出
		panic(err)
	}
	// 建表
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

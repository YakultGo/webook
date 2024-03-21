package middleware

import (
	myJwt "basic-go/webook/internal/web/jwt"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths map[string]bool
	myJwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHandler myJwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		paths:   make(map[string]bool),
		Handler: jwtHandler,
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	l.paths[path] = true
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 注册time.Time类型，否则session无法存储
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 不需要登录校验
		if ok, val := l.paths[ctx.Request.URL.Path]; ok && val == true {
			return
		}
		// 使用JWT校验
		tokenStr := l.Handler.ExtractToken(ctx)
		claims := &myJwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return myJwt.AccessTokenKey, nil
		})
		if err != nil {
			ctx.String(http.StatusUnauthorized, "未登录")
			fmt.Println("解析token失败", err)
			ctx.Abort()
			return
		}
		if !token.Valid {
			ctx.String(http.StatusUnauthorized, "未登录")
			fmt.Println("token无效")
			ctx.Abort()
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.String(http.StatusUnauthorized, "恶意入侵")
			ctx.Abort()
			return
		}
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			fmt.Println("session校验失败", err)
			ctx.String(http.StatusUnauthorized, "未登录")
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
	}
}

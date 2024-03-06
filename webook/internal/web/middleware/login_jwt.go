package middleware

import (
	"basic-go/webook/internal/web"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths map[string]bool
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{paths: make(map[string]bool)}
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
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.String(http.StatusUnauthorized, "未登录")
			ctx.Abort()
			return
		}
		// Bearer token
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.String(http.StatusUnauthorized, "未登录")
			ctx.Abort()
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("4a1LwMzFjaCW4HrJETQsR8ybdYq82WMV"), nil
		})
		if err != nil {
			ctx.String(http.StatusUnauthorized, "未登录")
			ctx.Abort()
			return
		}
		if !token.Valid {
			ctx.String(http.StatusUnauthorized, "未登录")
			ctx.Abort()
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.String(http.StatusUnauthorized, "恶意入侵")
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
	}
}

package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths map[string]bool
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{paths: make(map[string]bool)}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths[path] = true
	return l
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 注册time.Time类型，否则session无法存储
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 不需要登录校验
		if ok, val := l.paths[ctx.Request.URL.Path]; ok && val == true {
			return
		}
		//if ctx.Request.URL.Path == "/users/signup" || ctx.Request.URL.Path == "/users/login" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.String(http.StatusUnauthorized, "未登录")
			ctx.Abort()
			return
		}
		updateTime := sess.Get("updateTime")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60 * 60,
		})
		now := time.Now()
		// 第一次登录
		if updateTime == nil {
			sess.Set("updateTime", now)
			sess.Save()
			return
		}
		updateTimeVal := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Minute {
			sess.Set("updateTime", now)
			sess.Save()
			return
		}
	}
}

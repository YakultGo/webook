package web

import (
	"basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/author", h.AuthURL)
	g.GET("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "获取授权地址失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	ctx.String(http.StatusOK, "wechat callback")
}

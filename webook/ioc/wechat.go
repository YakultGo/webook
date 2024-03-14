package ioc

import "basic-go/webook/internal/service/oauth2/wechat"

func InitOAuth2WechatService() wechat.Service {
	return wechat.NewOAuthService("appid")
}

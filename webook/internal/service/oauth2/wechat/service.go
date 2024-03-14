package wechat

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
}

type OAuthService struct {
	appid string
}

func NewOAuthService(appid string) Service {
	return &OAuthService{
		appid: appid,
	}
}

func (o *OAuthService) AuthURL(ctx context.Context) (string, error) {
	const urlPattern = `appid=%s&redirect=%s&state=%s`
	const redirectURI = `https://xxxxx/oauth2/wechat/callback`
	state := uuid.New()
	return fmt.Sprintf(urlPattern, o.appid, redirectURI, state), nil
}

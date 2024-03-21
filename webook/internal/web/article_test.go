package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	svcmocks "basic-go/webook/internal/service/mocks"
	myJwt "basic-go/webook/internal/web/jwt"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantBody Result
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 1,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
			{
		"title": "我的标题",
		"content": "我的内容"
			}			
`,
			wantCode: 200,
			wantBody: Result{
				Data: float64(1),
				Msg:  "发表成功",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			control := gomock.NewController(t)
			defer control.Finish()
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("claims", myJwt.UserClaims{
					UserId: 1,
				})
				ctx.Next()
			})
			h := NewArticleHandler(tc.mock(control))
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/articles/publish",
				bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if tc.wantCode == http.StatusOK {
				return
			}
			var res Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}

package integration

import (
	"basic-go/webook/internal/repository/dao"
	myJwt "basic-go/webook/internal/web/jwt"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ArticleTestSuite 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

// SetupTest 测试用例执行前的初始化
func (a *ArticleTestSuite) SetupTest() {
	a.server = InitWebServer()
	a.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", myJwt.UserClaims{
			UserId: 123,
		})
		ctx.Next()
	})
	a.db = initTestGORM()
	a.server.Run(":8080")
}

// TearDownTest 测试用例执行后的资源释放
func (a *ArticleTestSuite) TearDownTest() {
	//a.db.Exec("TRUNCATE TABLE articles")
}
func (a *ArticleTestSuite) TestEdit() {
	t := a.T()
	now := time.Now().UnixMilli()
	testCases := []struct {
		name     string
		art      Article
		before   func(t *testing.T)
		after    func(t *testing.T)
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建帖子-保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				err := a.db.Where("id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					Id:         1,
					Title:      "标题",
					Content:    "内容",
					AuthorId:   1,
					CreateTime: now,
					UpdateTime: now,
				}, art)
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/articles/edit",
				bytes.NewReader(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTEyNjI1ODksIlVzZXJJZCI6MSwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4zNy4wIiwiU3NpZCI6IjViNDdiYzlmLWYwOTItNDYwMS1iNzIwLWJhNmU5YjJmNDdjYiJ9.BW7NZWy3u0VnXB23xSapWB0rzA0OVnsoNJdBp-x_gvu6N5G_WXVSSV0ULWc-YCE7uqQ8qRHRtt8ZhyP_-bDCrA")
			resp := httptest.NewRecorder()
			a.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			tc.after(t)
		})
	}
}
func (a *ArticleTestSuite) TestABC() {
	a.T().Log("test abc")
}
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data T      `json:"data"`
}

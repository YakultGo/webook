package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	myJwt "basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc      service.ArticleService
	interSvc service.InteractiveService
}

func NewArticleHandler(svc service.ArticleService, interSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:      svc,
		interSvc: interSvc,
	}
}
func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
	g.POST("/list", a.List)
	g.GET("/detail/:id", a.Detail)
	g.POST("/like", a.Like)
	g.POST("/fav", a.Favorite)
}
func (a *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	articleId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		zap.S().Error("前端输入的ID不对", zap.Error(err))
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		zap.S().Error("未发现用户的jwt信息")
		return
	}
	article, err := a.svc.Detail(ctx, articleId, claims.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.S().Error("查询文章详情失败", zap.Error(err))
		return
	}
	// 增加阅读计数
	go func() {
		err2 := a.interSvc.IncrReadCount(ctx, articleId)
		if err2 != nil {
			zap.S().Error("增加阅读计数失败", zap.Error(err))
		}
	}()
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: article,
	})
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		zap.S().Error("未发现用户的jwt信息")
		return
	}
	err := a.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.UserId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.S().Error("撤回失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}
func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		zap.S().Error("未发现用户的jwt信息")
		return
	}
	// 检测输入、先跳过
	id, err := a.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.UserId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.S().Error("保存失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		zap.S().Error("未发现用户的jwt信息")
		return
	}
	id, err := a.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.UserId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.S().Error("发表帖子失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (a *ArticleHandler) List(ctx *gin.Context) {
	type ListReq struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}
	var req ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		zap.S().Error("未发现用户的jwt信息")
		return
	}
	articles, err := a.svc.List(ctx, claims.UserId, req.Offset, req.Limit)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		zap.S().Error("查询文章列表失败", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: a.toVo(articles),
	})
}
func (a *ArticleHandler) toVo(articles []domain.Article) []ArticleVO {
	var vos []ArticleVO
	for _, article := range articles {
		vos = append(vos, ArticleVO{
			Id:         article.Id,
			Title:      article.Title,
			Abstract:   article.Abstract(),
			AuthorId:   article.Author.Id,
			AuthorName: article.Author.Name,
			Status:     article.Status.ToUint8(),
			CreateTime: article.CreateTime,
		})
	}
	return vos
}

func (a *ArticleHandler) Like(ctx *gin.Context) {
	type Req struct {
		ArticleId int64 `json:"id"`
		Like      bool  `json:"like"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		zap.S().Error("未发现用户的jwt信息")
		return
	}
	article := domain.Article{
		Id: req.ArticleId,
		Author: domain.Author{
			Id: claims.UserId,
		},
	}
	if req.Like {
		err := a.interSvc.Like(ctx, article)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			zap.S().Error("点赞失败", zap.Error(err))
			return
		}
	} else {
		err := a.interSvc.UnLike(ctx, article)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			zap.S().Error("取消点赞失败", zap.Error(err))
			return
		}

	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (a *ArticleHandler) Favorite(ctx *gin.Context) {
	type Req struct {
		ArticleId int64 `json:"id"`
		Fav       bool  `json:"fav"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "参数错误")
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		zap.S().Error("未发现用户的jwt信息")
		return
	}
	article := domain.Article{
		Id: req.ArticleId,
		Author: domain.Author{
			Id: claims.UserId,
		},
	}
	if req.Fav {
		err := a.interSvc.Fav(ctx, article)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			zap.S().Error("收藏失败", zap.Error(err))
			return
		}
	} else {
		err := a.interSvc.UnFav(ctx, article)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			zap.S().Error("取消收藏失败", zap.Error(err))
			return
		}

	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleVO struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
	AuthorId   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	Status     uint8  `json:"status"`
	CreateTime int64  `json:"create_time"`
}

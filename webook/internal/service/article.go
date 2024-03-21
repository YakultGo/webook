package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
	"github.com/gin-gonic/gin"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	List(ctx context.Context, userId int64, offset int, limit int) ([]domain.Article, error)
	Detail(ctx *gin.Context, articleId int64, userId int64) (domain.Article, error)
}

type ArticleServiceStruct struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &ArticleServiceStruct{
		repo: repo,
	}
}
func (a *ArticleServiceStruct) Detail(ctx *gin.Context, articleId int64, userId int64) (domain.Article, error) {
	return a.repo.GetArticleById(ctx, articleId, userId)
}

func (a *ArticleServiceStruct) Withdraw(ctx context.Context, art domain.Article) error {
	art.Status = domain.ArticleStatusPrivate
	return a.repo.SyncStatus(ctx, art)
}

func (a *ArticleServiceStruct) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnPublished
	if art.Id > 0 {
		return art.Id, a.repo.Update(ctx, art)
	}
	return a.repo.Create(ctx, art)
}

func (a *ArticleServiceStruct) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	if art.Id > 0 {
		return art.Id, a.repo.Update(ctx, art)
	}
	return a.repo.Create(ctx, art)
}

func (a *ArticleServiceStruct) List(ctx context.Context, userId int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, userId, offset, limit)
}

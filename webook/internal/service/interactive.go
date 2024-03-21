package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
	"github.com/gin-gonic/gin"
)

type InteractiveService interface {
	IncrReadCount(ctx context.Context, articleId int64) error
	Like(ctx *gin.Context, article domain.Article) error
	UnLike(ctx *gin.Context, article domain.Article) error
	Fav(ctx *gin.Context, article domain.Article) error
	UnFav(ctx *gin.Context, article domain.Article) error
}

type CachedInteractiveService struct {
	interRepo repository.InteractiveRepository
}

func NewInteractiveService(interRepo repository.InteractiveRepository) InteractiveService {
	return &CachedInteractiveService{
		interRepo: interRepo,
	}
}
func (c *CachedInteractiveService) Fav(ctx *gin.Context, article domain.Article) error {
	return c.interRepo.IncreaseFavCount(ctx, article.Author.Id, article.Id)
}

func (c *CachedInteractiveService) UnFav(ctx *gin.Context, article domain.Article) error {
	return c.interRepo.DecreaseFavCount(ctx, article.Author.Id, article.Id)
}

func (c *CachedInteractiveService) Like(ctx *gin.Context, article domain.Article) error {
	return c.interRepo.IncreaseLikeCount(ctx, article.Author.Id, article.Id)
}

func (c *CachedInteractiveService) UnLike(ctx *gin.Context, article domain.Article) error {
	return c.interRepo.DecreaseLikeCount(ctx, article.Author.Id, article.Id)
}

func (c *CachedInteractiveService) IncrReadCount(ctx context.Context, articleId int64) error {
	return c.interRepo.IncreaseReadCount(ctx, articleId)
}

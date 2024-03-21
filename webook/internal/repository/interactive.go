package repository

import (
	"basic-go/webook/internal/repository/dao"
	"context"
)

type InteractiveRepository interface {
	IncreaseReadCount(ctx context.Context, articleId int64) error
	IncreaseLikeCount(ctx context.Context, userId, articleId int64) error
	DecreaseLikeCount(ctx context.Context, userId, articleId int64) error
	IncreaseFavCount(ctx context.Context, userId, articleId int64) error
	DecreaseFavCount(ctx context.Context, userId, articleId int64) error
}

type CachedInteractiveRepository struct {
	interactiveDao dao.InteractiveDao
}

func NewCachedInteractiveRepository(interactiveDao dao.InteractiveDao) InteractiveRepository {
	return &CachedInteractiveRepository{
		interactiveDao: interactiveDao,
	}
}
func (c *CachedInteractiveRepository) IncreaseFavCount(ctx context.Context, userId, articleId int64) error {
	return c.interactiveDao.InsertFavInfo(ctx, userId, articleId)
}

func (c *CachedInteractiveRepository) DecreaseFavCount(ctx context.Context, userId, articleId int64) error {
	return c.interactiveDao.DeleteFavInfo(ctx, userId, articleId)
}

func (c *CachedInteractiveRepository) DecreaseLikeCount(ctx context.Context, userId, articleId int64) error {
	return c.interactiveDao.DeleteLikeInfo(ctx, userId, articleId)
}

func (c *CachedInteractiveRepository) IncreaseLikeCount(ctx context.Context, userId int64, articleId int64) error {
	return c.interactiveDao.InsertLikeInfo(ctx, userId, articleId)
}

func (c *CachedInteractiveRepository) IncreaseReadCount(ctx context.Context, articleId int64) error {
	return c.interactiveDao.IncrReadCnt(ctx, articleId)
}

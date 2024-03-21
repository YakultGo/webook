package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	SyncStatus(ctx context.Context, art domain.Article) error
	List(ctx context.Context, userId int64, offset int, limit int) ([]domain.Article, error)
	GetArticleById(ctx *gin.Context, articleId int64, userId int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
}

func NewArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache) ArticleRepository {
	return &CachedArticleRepository{
		dao:   dao,
		cache: cache,
	}
}

func (a *CachedArticleRepository) GetArticleById(ctx *gin.Context, articleId int64, userId int64) (domain.Article, error) {
	art, err := a.dao.GetArticleById(ctx, articleId, userId)
	if err != nil {
		return domain.Article{}, err
	}
	return domain.Article{
		Id:         art.Id,
		Title:      art.Title,
		Content:    art.Content,
		CreateTime: art.CreateTime,
		Status:     domain.ArticleStatus(art.Status),
	}, nil
}

func (a *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Insert(ctx, a.toEntity(art))
}

func (a *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	defer func() {
		a.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return a.dao.UpdateById(ctx, a.toEntity(art))
}

func (a *CachedArticleRepository) SyncStatus(ctx context.Context, art domain.Article) error {
	return a.dao.SyncStatus(ctx, a.toEntity(art))
}

func (a *CachedArticleRepository) List(ctx context.Context, userId int64, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		data, err := a.cache.GetFirstPage(ctx, userId)
		if err == nil {
			go func() {
				a.preCache(ctx, data)
			}()
			return data, nil
		}
	}
	arts, err := a.dao.GetByAuthor(ctx, userId, offset, limit)
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	for _, v := range arts {
		res = append(res, domain.Article{
			Id:         v.Id,
			Title:      v.Title,
			Content:    v.Content,
			CreateTime: v.CreateTime,
			Status:     domain.ArticleStatus(v.Status),
			Author: domain.Author{
				Id: v.AuthorId,
			},
		})
	}
	go func() {
		// 回写缓存
		err = a.cache.SetFirstPage(ctx, userId, res)
		if err != nil {
			zap.S().Errorf("回写缓存失败 %v", err)
		}
		a.preCache(ctx, res)
	}()
	return res, nil
}

func (a *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}

func (a *CachedArticleRepository) preCache(ctx context.Context, res []domain.Article) {
	if len(res) > 0 {
		// 缓存第一页
		err := a.cache.Set(ctx, res[0])
		if err != nil {
			zap.S().Errorf("预热缓存失败 %v", err)
		}
	}
}

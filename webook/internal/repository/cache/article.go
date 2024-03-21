package cache

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, author int64, articles []domain.Article) error
	DelFirstPage(ctx context.Context, author int64) error
	Set(ctx context.Context, article domain.Article) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: cmd,
	}
}
func (r *RedisArticleCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error) {
	// 获取缓存
	data, err := r.client.Get(ctx, r.firstPageKey(author)).Bytes()
	if err != nil {
		return nil, err

	}
	var articles []domain.Article
	err = json.Unmarshal(data, &articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, author int64, articles []domain.Article) error {
	for i := 0; i < len(articles); i++ {
		articles[i].Content = articles[i].Abstract()
	}
	data, err := json.Marshal(articles)
	if err != nil {
		return err

	}
	return r.client.Set(ctx, r.firstPageKey(author), data, time.Minute*10).Err()
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, author int64) error {
	// 删除缓存
	return r.client.Del(ctx, r.firstPageKey(author)).Err()
}

func (r *RedisArticleCache) Set(ctx context.Context, art domain.Article) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(art.Id), data, time.Second*30).Err()
}
func (r *RedisArticleCache) firstPageKey(id int64) string {
	return fmt.Sprintf("article:first_page:%d", id)
}
func (r *RedisArticleCache) key(userId int64) string {
	return fmt.Sprintf("article:%d", userId)
}

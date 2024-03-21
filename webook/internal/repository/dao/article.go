package dao

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	Id         int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title      string `gorm:"type:varchar(1024)" bson:"title,omitempty"`
	Content    string `gorm:"type:blob" bson:"content,omitempty"`
	AuthorId   int64  `gorm:"index=aid_ctime" bson:"author_id, omitempty"`
	CreateTime int64  `gorm:"index=aid_ctime;default:NULL" bson:"create_time,omitempty"`
	UpdateTime int64  `gorm:"default:NULL" bson:"update_time,omitempty"`
	Status     uint8  `bson:"status,omitempty"`
}
type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	SyncStatus(ctx context.Context, art Article) error
	GetByAuthor(ctx context.Context, id int64, offset int, limit int) ([]Article, error)
	GetArticleById(ctx *gin.Context, articleId int64, userId int64) (Article, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (g *GORMArticleDAO) GetArticleById(ctx *gin.Context, articleId int64, userId int64) (Article, error) {
	var art Article
	err := g.db.WithContext(ctx).
		Where("id = ? and author_id = ?", articleId, userId).
		First(&art).Error
	return art, err
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}
func (g *GORMArticleDAO) GetByAuthor(ctx context.Context, id int64, offset int, limit int) ([]Article, error) {
	var arts []Article
	err := g.db.WithContext(ctx).
		Model(&Article{}).
		Where("author_id = ?", id).
		Offset(offset).
		Limit(limit).
		Order("update_time desc").
		Find(&arts).Error
	return arts, err
}
func (g *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.CreateTime = now
	art.UpdateTime = now
	result := g.db.WithContext(ctx).Create(&art)
	return art.Id, result.Error
}

func (g *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.UpdateTime = now
	result := g.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? and author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":       art.Title,
			"content":     art.Content,
			"update_time": art.UpdateTime,
			"status":      art.Status,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("更新失败，创作者非法创建%d, author_id %d", art.Id, art.AuthorId)
	}
	return nil
}

func (g *GORMArticleDAO) SyncStatus(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	result := g.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? and author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]interface{}{
			"status":      art.Status,
			"update_time": now,
		})
	return result.Error
}

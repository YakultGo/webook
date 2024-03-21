package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Interactive struct {
	Id         int64 `gorm:"primaryKey,autoIncrement""`
	ArticleId  int64 `gorm:"not null; unique"`
	ReadCnt    int64 `gorm:"default:0"`
	LikeCnt    int64 `gorm:"default:0"`
	FavCnt     int64 `gorm:"default:0"`
	CreateTime int64 `gorm:"default:NULL"`
	UpdateTime int64 `gorm:"default:NULL"`
}

type UserLikeArticle struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 建立联合唯一索引
	UserID    int64 `gorm:"uniqueIndex:like_user_article"`
	ArticleID int64 `gorm:"uniqueIndex:like_user_article"`
	// 0:未点赞 1:已点赞
	Status     uint8 `gorm:"default:0"`
	CreateTime int64 `gorm:"default:NULL"`
	UpdateTime int64 `gorm:"default:NULL"`
}

type UserFavArticle struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 建立联合唯一索引
	UserID    int64 `gorm:"uniqueIndex:fav_user_article"`
	ArticleID int64 `gorm:"uniqueIndex:fav_user_article"`
	// 0:未点赞 1:已点赞
	Status     uint8 `gorm:"default:0"`
	CreateTime int64 `gorm:"default:NULL"`
	UpdateTime int64 `gorm:"default:NULL"`
}
type InteractiveDao interface {
	IncrReadCnt(ctx context.Context, articleId int64) error
	InsertLikeInfo(ctx context.Context, userId, articleId int64) error
	DeleteLikeInfo(ctx context.Context, userId, articleId int64) error
	InsertFavInfo(ctx context.Context, userId, articleId int64) error
	DeleteFavInfo(ctx context.Context, userId, articleId int64) error
}

type GORMInteractiveDao struct {
	db *gorm.DB
}

func NewInteractiveDao(db *gorm.DB) InteractiveDao {
	return &GORMInteractiveDao{
		db: db,
	}
}

func (g *GORMInteractiveDao) InsertFavInfo(ctx context.Context, userId, articleId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"status":      1,
				"update_time": now,
			}),
		}).Create(&UserFavArticle{
			UserID:     userId,
			ArticleID:  articleId,
			Status:     1,
			CreateTime: now,
			UpdateTime: now,
		})
		if result.Error != nil {
			return result.Error
		}
		// 更新文章点赞数
		result = g.db.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"fav_cnt":     gorm.Expr(" `fav_cnt` + ?", 1),
				"update_time": now,
			}),
		}).Create(&Interactive{
			FavCnt:     1,
			ArticleId:  articleId,
			CreateTime: now,
			UpdateTime: now,
		})
		return result.Error
	})
}

func (g *GORMInteractiveDao) DeleteFavInfo(ctx context.Context, userId, articleId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&UserFavArticle{}).
			Where("user_id = ? and article_id = ?", userId, articleId).
			Updates(map[string]any{
				"status":      0,
				"update_time": now,
			})
		if result.Error != nil {
			return result.Error
		}
		// 更新文章点赞数
		result = tx.Model(&Interactive{}).
			Where("article_id = ?", articleId).
			Updates(map[string]any{
				"fav_cnt":     gorm.Expr(" `fav_cnt` - ?", 1),
				"update_time": now,
			})
		return result.Error
	})
}

func (g *GORMInteractiveDao) DeleteLikeInfo(ctx context.Context, userId, articleId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&UserLikeArticle{}).
			Where("user_id = ? and article_id = ?", userId, articleId).
			Updates(map[string]any{
				"status":      0,
				"update_time": now,
			})
		if result.Error != nil {
			return result.Error
		}
		// 更新文章点赞数
		result = tx.Model(&Interactive{}).
			Where("article_id = ?", articleId).
			Updates(map[string]any{
				"like_cnt":    gorm.Expr(" `like_cnt` - ?", 1),
				"update_time": now,
			})
		return result.Error
	})
}

func (g *GORMInteractiveDao) InsertLikeInfo(ctx context.Context, userId, articleId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"status":      1,
				"update_time": now,
			}),
		}).Create(&UserLikeArticle{
			UserID:     userId,
			ArticleID:  articleId,
			Status:     1,
			CreateTime: now,
			UpdateTime: now,
		})
		if result.Error != nil {
			return result.Error
		}
		// 更新文章点赞数
		result = g.db.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt":    gorm.Expr(" `like_cnt` + ?", 1),
				"update_time": now,
			}),
		}).Create(&Interactive{
			LikeCnt:    1,
			ArticleId:  articleId,
			CreateTime: now,
			UpdateTime: now,
		})
		return result.Error
	})
}

func (g *GORMInteractiveDao) IncrReadCnt(ctx context.Context, articleId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt":    gorm.Expr(" `  ` + ?", 1),
			"update_time": now,
		}),
	}).Create(&Interactive{
		LikeCnt:    1,
		ArticleId:  articleId,
		CreateTime: now,
		UpdateTime: now,
	}).Error
}

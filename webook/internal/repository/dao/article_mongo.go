package dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoDBDAO struct {
	//client   *mongo.Client
	//database *mongo.Database
	col  *mongo.Collection
	node *snowflake.Node
}

func (m *MongoDBDAO) GetByAuthor(ctx context.Context, id int64, offset int, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) GetArticleById(ctx *gin.Context, articleId int64, userId int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDAO{
		col:  db.Collection("articles"),
		node: node,
	}
}
func (m *MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	id := m.node.Generate().Int64()
	art.Id = id
	art.CreateTime = now
	art.UpdateTime = now
	_, err := m.col.InsertOne(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *MongoDBDAO) UpdateById(ctx context.Context, art Article) error {
	filter := bson.M{"author_id": art.AuthorId, "id": art.Id}
	update := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title":       art.Title,
		"content":     art.Content,
		"update_time": time.Now().UnixMilli(),
		"status":      art.Status,
	}}}
	//update := bson.D{bson.E{Key: "$set", Value: art}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func (m *MongoDBDAO) SyncStatus(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

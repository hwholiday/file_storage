package fileInfo

import (
	"context"
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"strings"
)

var (
	s *fileInfo
)

type fileInfo struct {
	mgo *mongo.Database
}

type Service interface {
}

func NewFileInfo(mClient *mongo.Client, c *conf.Config) *fileInfo {
	s = &fileInfo{
		mgo: mClient.Database(c.Mongo.DataBase),
	}
	s.CreateIndex()
	return s
}

func GetServer() Service {
	return s
}

func (f *fileInfo) CreateIndex() {
	var (
		//对 fileinfo 库创建索引
		fileInfo = &entity.FileInfo{}
	)
	indexView := s.mgo.Collection(fileInfo.TableName()).Indexes()
	cursor, err := indexView.List(context.Background())
	if err != nil {
		log.GetLogger().Panic("[NewFileInfo] List", zap.Any("name", fileInfo.TableName()), zap.Error(err))
		return
	}
	var hasFidIndex bool
	for cursor.Next(context.Background()) {
		if strings.Contains(cursor.Current.String(), "fid_1") {
			hasFidIndex = true
			break
		}
	}
	if !hasFidIndex {
		_, err := indexView.CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.D{{"fid", 1}},
		})
		if err != nil {
			log.GetLogger().Panic("[NewFileInfo] CreateIndex", zap.Any("name", fileInfo.TableName()), zap.Error(err))
			return
		}
		log.GetLogger().Info("[NewFileInfo] CreateIndex fid success", zap.Any("name", fileInfo.TableName()))

	}
}

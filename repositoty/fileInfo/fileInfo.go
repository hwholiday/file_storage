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

type fileInfo struct {
	mgo *mongo.Database
}

type Service interface {
	InsertFileInfo(f *entity.FileInfo) (err error)
	DelFileInfoByFid(fid int64) (err error)
	GetFileInfoByFid(fid int64) (fileInfo *entity.FileInfo, err error)
	GetFileInfoByMd5(md5 string) (fileInfo *entity.FileInfo, err error)
	UpdateFileInfoStatusByFid(fid int64, status int) (err error)
	UpdateFileInfoByFid(fid int64, change interface{}) (err error)
}

func NewFileInfo(mClient *mongo.Client, c *conf.Config) *fileInfo {
	s := &fileInfo{
		mgo: mClient.Database(c.Mongo.DataBase),
	}
	s.createIndex()
	return s
}

func (f *fileInfo) createIndex() {
	var (
		//对 fileinfo 库创建索引
		fileInfo = &entity.FileInfo{}
	)
	indexView := f.mgo.Collection(fileInfo.TableName()).Indexes()
	cursor, err := indexView.List(context.Background())
	if err != nil {
		log.GetLogger().Panic("[NewFileInfo] List", zap.Any("name", fileInfo.TableName()), zap.Error(err))
		return
	}
	var hasFidIndex bool
	var hasMd5Index bool
	for cursor.Next(context.Background()) {
		if strings.Contains(cursor.Current.String(), "fid_1") {
			hasFidIndex = true
		}
		if strings.Contains(cursor.Current.String(), "md5_1") {
			hasMd5Index = true
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
	if !hasMd5Index {
		_, err := indexView.CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.D{{"md5", 1}},
		})
		if err != nil {
			log.GetLogger().Panic("[NewFileInfo] CreateIndex", zap.Any("name", fileInfo.TableName()), zap.Error(err))
			return
		}
		log.GetLogger().Info("[NewFileInfo] CreateIndex md5 success", zap.Any("name", fileInfo.TableName()))

	}

}

package fileInfo

import (
	"context"
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (s *fileInfo) InsertFileInfo(f *entity.FileInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), conf.MgoContextTimeOut)
	defer cancel()
	if _, err := s.mgo.Collection(f.TableName()).InsertOne(ctx, f); err != nil {
		log.GetLogger().Error("[InsertFileInfo] InsertOne", zap.Any("name", f.TableName()), zap.Any("fileInfo", f), zap.Error(err))
		return err
	}
	return nil
}

func (s *fileInfo) DelFileInfoByFid(fid int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), conf.MgoContextTimeOut)
	defer cancel()
	var tableName = entity.FileInfo{}.TableName()
	if _, err := s.mgo.Collection(tableName).DeleteOne(ctx, bson.M{"fid": fid}); err != nil {
		log.GetLogger().Error("[DelFileInfoByFid] DelFileInfo", zap.Any("name", tableName), zap.Any("fid", fid), zap.Error(err))
		return err
	}
	return nil
}

func (s *fileInfo) GetFileInfoByFid(fid int64) (fileInfo *entity.FileInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), conf.MgoContextTimeOut)
	defer cancel()
	fileInfo = new(entity.FileInfo)
	tableName := fileInfo.TableName()
	if err = s.mgo.Collection(tableName).FindOne(ctx, bson.M{"fid": fid}).Decode(fileInfo); err != nil {
		log.GetLogger().Error("[GetFileInfoByFid] FindOne", zap.Any("name", tableName), zap.Any("fid", fid), zap.Error(err))
		return
	}
	return
}

package fileInfo

import (
	"context"
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (s *fileInfo) InsertFileInfo(f *entity.FileInfo) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), conf.MgoContextTimeOut)
	defer cancel()
	if _, err = s.mgo.Collection(f.TableName()).InsertOne(ctx, f); err != nil {
		log.GetLogger().Error("[InsertFileInfo] InsertOne", zap.Any("name", f.TableName()), zap.Any("fileInfo", f), zap.Error(err))
		return
	}
	return
}

func (s *fileInfo) DelFileInfoByFid(fid int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), conf.MgoContextTimeOut)
	defer cancel()
	var tableName = entity.FileInfo{}.TableName()
	if _, err = s.mgo.Collection(tableName).DeleteOne(ctx, bson.M{"fid": fid}); err != nil {
		log.GetLogger().Error("[DelFileInfoByFid] DelFileInfo", zap.Any("name", tableName), zap.Any("fid", fid), zap.Error(err))
		return
	}
	return
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

func (s *fileInfo) UpdateFileInfoStatusByFid(fid int64, status int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), conf.MgoContextTimeOut)
	defer cancel()
	tableName := entity.FileInfo{}.TableName()
	if _, err = s.mgo.Collection(tableName).UpdateOne(ctx, bson.M{"fid": fid}, bson.D{{"$set", bson.D{{"status", status}}}}); err != nil {
		log.GetLogger().Error("[UpdateFileInfoStatusByFid] UpdateOne", zap.Any("name", tableName), zap.Any("fid", fid), zap.Any("status", status), zap.Error(err))
		return
	}
	return
}

func (s *fileInfo) UpdateFileInfoByFid(fid int64, change interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), conf.MgoContextTimeOut)
	defer cancel()
	tableName := entity.FileInfo{}.TableName()
	if _, err = s.mgo.Collection(tableName).UpdateOne(ctx, bson.M{"fid": fid}, change); err != nil {
		log.GetLogger().Error("[UpdateFileInfoByFid] UpdateOne", zap.Any("name", tableName), zap.Any("fid", fid), zap.Any("status", change), zap.Error(err))
		return
	}
	return
}

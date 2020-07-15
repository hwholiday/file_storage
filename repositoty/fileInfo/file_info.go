package fileInfo

import (
	"context"
	"filesrv/entity"
	"filesrv/library/log"
	"go.uber.org/zap"
	"time"
)

func (s *fileInfo) InsertFileInfo(f *entity.FileInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := s.mgo.Collection(f.TableName()).InsertOne(ctx, f); err != nil {
		log.GetLogger().Error("[InsertFileInfo] InsertOne", zap.Any("name", f.TableName()), zap.Any("fileInfo", f), zap.Error(err))
	}
	return nil
}

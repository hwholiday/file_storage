package storage

import (
	"bytes"
	"filesrv/conf"
	"filesrv/library/log"
	"fmt"
	"github.com/minio/minio-go"
	"go.uber.org/zap"
)

func (s *storage) UpFile(fid int64, bucketName string, data []byte) error {
	size, err := s.mClient.PutObject(bucketName, fmt.Sprint(fid), bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		log.GetLogger().Error("[UpFileNotSlice] PutObject", zap.Any(bucketName, fid), zap.Error(err))
		return err
	}
	if size != int64(len(data)) {
		log.GetLogger().Error("[UpFileNotSlice]   data inconsistently", zap.Any(bucketName, fid))
		return conf.ErrFileSizeInvalid
	}
	log.GetLogger().Debug("[UpFileNotSlice]", zap.Any("Fid", fid), zap.Any("bucketName", bucketName), zap.Any("data", len(data)))
	return nil
}

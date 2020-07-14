package repositoty

import (
	"bytes"
	"filesrv/conf"
	"filesrv/library/log"
	"fmt"
	"github.com/minio/minio-go"
	"go.uber.org/zap"
	"io/ioutil"
)

func (s *Repository) StorageUpFileNotSlice(fid string, bucketName string, data []byte) error {
	size, err := s.mClient.PutObject(bucketName, fid, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		log.GetLogger().Error("[StorageUpFileNotSlice] PutObject", zap.Any(bucketName, fid), zap.Error(err))
		return err
	}
	if size != int64(len(data)) {
		log.GetLogger().Error("[StorageUpFileNotSlice]   data inconsistently", zap.Any(bucketName, fid))
		return conf.ErrFileSizeInvalid
	}
	log.GetLogger().Debug("[StorageUpFileNotSlice]", zap.Any("Fid", fid), zap.Any("bucketName", bucketName), zap.Any("data", len(data)))
	return nil
}

func (s *Repository) StorageGetFileNotSlice(fid int64, bucketName string) (error, []byte) {
	var (
		object *minio.Object
		err    error
		data   []byte
	)
	object, err = s.mClient.GetObject(bucketName, fmt.Sprint(fid), minio.GetObjectOptions{})
	if err != nil {
		log.GetLogger().Error("[StorageGetFileNotSlice] GetObject", zap.Any(bucketName, fid), zap.Error(err))
		return err, nil
	}
	defer object.Close()
	data, err = ioutil.ReadAll(object)
	if err != nil {
		log.GetLogger().Error("[StorageGetFileNotSlice] ioutil.ReadAll", zap.Any(bucketName, fid), zap.Error(err))
		return err, nil
	}
	log.GetLogger().Debug("[StorageGetFileNotSlice]", zap.Any("Fid", fid), zap.Any("bucketName", bucketName))
	return nil, data
}

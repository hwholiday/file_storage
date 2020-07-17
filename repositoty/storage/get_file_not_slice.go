package storage

import (
	"filesrv/library/log"
	"fmt"
	"github.com/minio/minio-go"
	"go.uber.org/zap"
	"io/ioutil"
)

func (s *storage) GetFileNotSlice(fid int64, bucketName string) ([]byte, error) {
	var (
		object *minio.Object
		err    error
		data   []byte
	)
	object, err = s.mClient.GetObject(bucketName, fmt.Sprint(fid), minio.GetObjectOptions{})
	if err != nil {
		log.GetLogger().Error("[GetFileNotSlice] GetObject", zap.Any(bucketName, fid), zap.Error(err))
		return nil, err
	}
	defer object.Close()
	data, err = ioutil.ReadAll(object)
	if err != nil {
		log.GetLogger().Error("[GetFileNotSlice] ioutil.ReadAll", zap.Any(bucketName, fid), zap.Error(err))
		return nil, err
	}
	log.GetLogger().Debug("[GetFileNotSlice]", zap.Any("Fid", fid), zap.Any("bucketName", bucketName))
	return data, nil
}

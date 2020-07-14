package repositoty

import (
	"filesrv/library/log"
	"github.com/minio/minio-go"
	"go.uber.org/zap"
	"io/ioutil"
)

func (s *Repository) StorageGetSliceFile(start, end int64, fileName, bucketName string) (error, []byte) {
	var opt minio.GetObjectOptions
	err := opt.SetRange(start, end)
	if err != nil {
		log.GetLogger().Error("[StorageGetSliceFile] SetRange", zap.Any(bucketName, fileName), zap.Error(err))
		return err, nil
	}
	object, err := s.mClient.GetObject(bucketName, fileName, opt)
	if err != nil {
		log.GetLogger().Error("[StorageGetSliceFile] GetObject", zap.Any(bucketName, fileName), zap.Error(err))
		return err, nil
	}
	defer object.Close()
	data, err := ioutil.ReadAll(object)
	if err != nil {
		log.GetLogger().Error("[StorageGetSliceFile] ReadAll", zap.Any(bucketName, fileName), zap.Error(err))
		return err, nil
	}
	log.GetLogger().Debug("[StorageGetSliceFile]", zap.Any(bucketName, fileName))
	return nil, data
}

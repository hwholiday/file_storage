package storage

import (
	"filesrv/library/log"
	"fmt"
	"github.com/minio/minio-go"
	"go.uber.org/zap"
	"io/ioutil"
)

func (s *storage) GetSliceFile(fid int64, bucketName string, start, end int64) ([]byte, error) {
	var opt minio.GetObjectOptions
	err := opt.SetRange(start, end)
	if err != nil {
		log.GetLogger().Error("[GetSliceFile] SetRange", zap.Any(bucketName, fid), zap.Error(err))
		return nil, err
	}
	object, err := s.mClient.GetObject(bucketName, fmt.Sprint(fid), opt)
	if err != nil {
		log.GetLogger().Error("[GetSliceFile] GetObject", zap.Any(bucketName, fid), zap.Error(err))
		return nil, err
	}
	defer object.Close()
	data, err := ioutil.ReadAll(object)
	if err != nil {
		log.GetLogger().Error("[GetSliceFile] ReadAll", zap.Any(bucketName, fid), zap.Error(err))
		return nil, err
	}
	log.GetLogger().Debug("[GetSliceFile]", zap.Any(bucketName, fid))
	return data, nil
}

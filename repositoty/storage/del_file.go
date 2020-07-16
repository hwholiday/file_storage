package storage

import (
	"filesrv/library/log"
	"go.uber.org/zap"
	"strconv"
)

func (s *storage) DelFile(fid int64, bucketName string) error {
	if err := s.mClient.RemoveObject(bucketName, strconv.Itoa(int(fid))); err != nil {
		log.GetLogger().Error("[DelFile]", zap.Any(bucketName, fid), zap.Error(err))
		return err
	}
	log.GetLogger().Debug("[DelFile]", zap.Any(bucketName, fid))
	return nil
}

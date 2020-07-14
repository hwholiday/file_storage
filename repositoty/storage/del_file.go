package storage

import (
	"filesrv/library/log"
	"go.uber.org/zap"
	"strconv"
)

func (s *storage) DelFile(fid int64, bucketName string) error {
	log.GetLogger().Debug("[DelFile]", zap.Any(bucketName, fid))
	return s.mClient.RemoveObject(bucketName, strconv.Itoa(int(fid)))
}

package storage

import (
	"filesrv/library/log"
	"go.uber.org/zap"
)

func (s *StorageServer) DelFile(bucketName, fileName string) error {
	log.GetLogger().Debug("[DelFile]", zap.Any(bucketName, fileName))
	return s.mClient.RemoveObject(bucketName, fileName)
}

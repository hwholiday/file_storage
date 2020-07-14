package repositoty

import (
	"filesrv/library/log"
	"go.uber.org/zap"
)

func (s *Repository) DbDelFile(fid int64) {
}

func (s *Repository) StorageDelFile(bucketName, fileName string) error {
	log.GetLogger().Debug("[StorageDelFile]", zap.Any(bucketName, fileName))
	return s.mClient.RemoveObject(bucketName, fileName)
}

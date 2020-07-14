package storage

import (
	"filesrv/common/storage/bucket"
	"filesrv/library/log"
	"filesrv/library/utils"
	"go.uber.org/zap"
)

func (s *storage) GetFidAndBucketName() (fid int64, name string) {
	fid = utils.GetSnowFlake().GetId()
	name = bucket.GetStorageBucket().GetRandBucketName()
	log.GetLogger().Debug("[GetFidAndBucketName]", zap.Int64("FID", fid), zap.Any("bucket", name))
	return
}

package service

import (
	"filesrv/common/storage/bucket"
	"filesrv/entity"
	"filesrv/library/utils"
)

func (s *service) ApplyFid(info *entity.FileInfo) (err error) {
	info.Fid = utils.GetSnowFlake().GetId()
	info.BucketName = bucket.GetStorageBucket().GetRandBucketName()
	info.IsImage = utils.IsImage(info.ExName)
	info.ContentType = utils.GetContentType(info.ExName)
	err = s.r.FileInfoServer.InsertFileInfo(info)
	return
}

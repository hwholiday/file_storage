package service

import (
	storage "filesrv/api/pb"
	"filesrv/common/storage/bucket"
	"filesrv/common/storage/manager"
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/log"
	"filesrv/library/utils"
	"go.uber.org/zap"
)

func (s *service) ApplyFid(info *storage.InApplyFid) (out *storage.OutApplyFid, err error) {

	var fInfo = s.convertDataToFileInfo(info)
	out = new(storage.OutApplyFid)
	out.Fid = fInfo.Fid
	status := s.addApplyFidIntoManager(fInfo)
	if status == conf.FileUploading {
		out.Status = status
		log.GetLogger().Info("[ApplyFid] addApplyFidIntoManager", zap.Any("find fid by manager", fInfo.Fid))
		return
	}
	if err = s.r.FileInfoServer.InsertFileInfo(fInfo); err != nil {
		s.f.DelItem(fInfo.Fid) //插入数据库失败，删除文件管理类
		log.GetLogger().Info("[ApplyFid] InsertFileInfo", zap.Any("fid", fInfo.Fid))
		return
	}
	return
}

func (s *service) convertDataToFileInfo(info *storage.InApplyFid) (fInfo *entity.FileInfo) {
	fInfo = &entity.FileInfo{}
	fInfo.Fid = utils.GetSnowFlake().GetId()
	fInfo.BucketName = bucket.GetStorageBucket().GetRandBucketName()
	fInfo.IsImage = utils.IsImage(info.ExName)
	fInfo.ContentType = utils.GetContentType(info.ExName)
	fInfo.Status = conf.FileWaitingForUpload
	fInfo.Name = info.Name
	fInfo.Size = info.Size
	fInfo.ExName = info.ExName
	fInfo.Md5 = info.Md5
	fInfo.SliceTotal = info.SliceTotal
	fInfo.ExpiredTime = info.ExpiredTime
	fInfo.CreateTime = utils.GetTimeUnix()
	fInfo.UpdateTime = utils.GetTimeUnix()
	return
}

func (s *service) addApplyFidIntoManager(info *entity.FileInfo) int32 {
	return s.f.NewItem(&manager.FileItem{
		Fid:        info.Fid,
		BucketName: info.BucketName,
		Size:       info.Size,
		Md5:        info.Md5,
		IsImage:    info.IsImage,
		SliceTotal: info.SliceTotal,
	})
}

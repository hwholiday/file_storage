package service

import (
	storage "filesrv/api/pb"
	"filesrv/common/storage/bucket"
	"filesrv/common/storage/manager"
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/log"
	"filesrv/library/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (s *service) ApplyFid(info *storage.InApplyFid) (out *storage.OutApplyFid, err error) {
	var needNew bool
	out, needNew, err = s.applyFidCheckExistByMd5(info)
	if err != nil {
		return
	}
	if !needNew { //代表存在数据,并且上传已经完成
		log.GetLogger().Debug("[ApplyFid] fid find", zap.Any("out", out))
		return
	}
	//新建文件信息
	var fInfo = s.convertDataToFileInfo(info)
	out.Fid = fInfo.Fid
	s.addApplyFidIntoManager(fInfo)
	if err = s.r.FileInfoServer.InsertFileInfo(fInfo); err != nil {
		s.f.DelItem(fInfo.Fid) //插入数据库失败，删除文件管理类
		log.GetLogger().Info("[ApplyFid] InsertFileInfo", zap.Any("fid", fInfo.Fid))
		return
	}
	log.GetLogger().Info("[ApplyFid] InsertFileInfo", zap.Any("fid", fInfo.Fid))
	out.Status = fInfo.Status
	return
}

func (s *service) applyFidCheckExistByMd5(info *storage.InApplyFid) (out *storage.OutApplyFid, needNew bool, err error) {
	var fileInfo *entity.FileInfo
	out = &storage.OutApplyFid{}
	fileInfo, err = s.GetFileInfoByMd5NotAutoClear(info.Md5)
	if err != nil {
		return
	}
	if fileInfo == nil {
		needNew = true
		return
	}
	out.Fid = fileInfo.Fid
	out.Status = fileInfo.Status
	if fileInfo.ExpiredTime > 0 && utils.GetTimeUnix() > fileInfo.ExpiredTime && fileInfo.Status == conf.FileExists {
		//文件已经过期,但是有人想上传该文件,对该文件续期
		var expiredTime int64
		if info.ExpiredTime == 0 {
			expiredTime = 0
		} else {
			expiredTime = utils.GetTimeUnix() + info.ExpiredTime
		}
		change := bson.D{{"expired_time", expiredTime}, {"update_time", utils.GetTimeUnix()}}
		if err = s.r.FileInfoServer.UpdateFileInfoByFid(fileInfo.Fid, change); err != nil {
			log.GetLogger().Error("[ApplyFid] applyFidCheckExistByMd5", zap.Any("fid", fileInfo.Fid), zap.Error(err))
			needNew = true
			return
		}
	}
	//查询到该文件为等待上传状态,处理这个状态时，申请md5相同的文件直接分配ID让用户上传
	if fileInfo.Status == conf.FileWaitingForUpload {
		needNew = true
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
	fInfo.ExImage = entity.ImageEx{
		Height: info.Height,
		Width:  info.Width,
	}
	fInfo.Md5 = info.Md5
	fInfo.SliceTotal = info.SliceTotal
	fInfo.ExpiredTime = info.ExpiredTime
	fInfo.CreateTime = utils.GetTimeUnix()
	fInfo.UpdateTime = utils.GetTimeUnix()
	return
}

func (s *service) addApplyFidIntoManager(info *entity.FileInfo) {
	s.f.NewItem(&manager.FileItem{
		Fid:        info.Fid,
		BucketName: info.BucketName,
		Size:       info.Size,
		Md5:        info.Md5,
		IsImage:    info.IsImage,
		SliceTotal: info.SliceTotal,
	})
}

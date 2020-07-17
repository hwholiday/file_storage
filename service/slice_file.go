package service

import (
	"crypto/md5"
	storage "filesrv/api/pb"
	"filesrv/common/storage/manager"
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/log"
	"fmt"
	"go.uber.org/zap"
)

func (s *service) UpSliceFile(in *storage.InUpSliceFileItem) (err error) {
	if err = s.f.AddItem(&manager.FileUploadItem{
		Fid:  in.Fid,
		Part: in.Part,
		Data: in.Data,
		Md5:  in.Md5,
	}); err != nil {
		log.GetLogger().Error("[UpSliceFile] AddItem", zap.Any("fid", in.Fid), zap.Any("part", in.Part), zap.Error(err))
		return
	}
	return
}

func (s *service) DownSliceFile(in *storage.InDownSliceFileItem) (out *storage.OutDownSliceFileItem, err error) {
	if in.Limit > 1048576 {
		err = conf.ErrLimitInvalid
		return
	}
	if in.Offset%1024 != 0 {
		err = conf.ErrOffsetInvalid
		return
	}
	var (
		fileInfo *entity.FileInfo
		data     []byte
	)
	out = &storage.OutDownSliceFileItem{}
	fileInfo, err = s.GetFileInfoByFid(in.Fid)
	if err != nil {
		return
	}
	if fileInfo == nil {
		err = conf.ErrFileIdInvalid
		return
	}
	if fileInfo.Status == conf.FileExpired {
		err = conf.ErrFileIdInvalid
		return
	}
	if in.Offset+in.Limit != fileInfo.Size {
		if in.Limit%1024 != 0 {
			err = conf.ErrLimitInvalid
			return
		}
	}
	var (
		start int64
		end   int64
	)
	start = in.Offset
	end = in.Limit + in.Offset
	if in.Offset != 0 {
		start++
	}
	log.GetLogger().Debug("range", zap.Int64("size", fileInfo.Size), zap.Int64("start", start), zap.Int64("end", end))
	data, err = s.r.StorageServer.GetSliceFile(fileInfo.Fid, fileInfo.BucketName, start, end)
	if err != nil {
		return
	}
	out.Fid = in.Fid
	out.Data = data
	out.Md5 = fmt.Sprintf("%x", md5.Sum(data))
	return
}

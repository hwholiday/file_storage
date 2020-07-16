package service

import (
	storage "filesrv/api/pb"
	"filesrv/common/storage/manager"
	"filesrv/library/log"
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

func (s *service) DownSliceFile(in *storage.InUpSliceFileItem) (err error) {
	return
}

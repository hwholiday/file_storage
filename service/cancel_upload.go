package service

import (
	storage "filesrv/api/pb"
	"filesrv/library/log"
	"go.uber.org/zap"
)

func (s *service) CancelByFid(info *storage.InCancelUpload) (err error) {
	if err = s.r.FileInfoServer.DelFileInfoByFid(info.Fid); err != nil {
		log.GetLogger().Error("[CancelByFid] DelFileInfoByFid", zap.Error(err))
		return
	}
	s.f.DelItem(info.Fid)
	log.GetLogger().Info("[CancelByFid]", zap.Any("fid", info.Fid))
	return
}

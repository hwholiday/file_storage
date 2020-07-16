package service

import (
	storage "filesrv/api/pb"
	"filesrv/common/storage/manager"
	"filesrv/conf"
	"filesrv/library/log"
	"filesrv/library/utils"
	"filesrv/repositoty"
	"go.uber.org/zap"
)

var s *service

type service struct {
	c *conf.Config
	r *repositoty.Repository
	f *manager.FileManager
}

type Service interface {
	ApplyFid(info *storage.InApplyFid) (out *storage.OutApplyFid, err error)
}

func NewService(c *conf.Config) {
	s = &service{
		c: c,
		r: repositoty.NewRepository(c),
	}
	err := utils.NewWorker(c.SnowFlakeId)
	if err != nil {
		log.GetLogger().Panic("[NewService] NewWorker", zap.Error(err))
		return
	}
	manager.NewFileManager(s.r)
	s.f = manager.GetFileManager()
}

func GetService() Service {
	return s
}

func (s *service) Close() {
	s.r.Close()
}

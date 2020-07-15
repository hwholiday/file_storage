package service

import (
	"filesrv/common/storage/manager"
	"filesrv/conf"
	"filesrv/library/log"
	"filesrv/library/utils"
	"filesrv/repositoty"
	"go.uber.org/zap"
)

var s *Service

type Service struct {
	c *conf.Config
	r *repositoty.Repository
	f *manager.FileManager
}

func NewService(c *conf.Config) {
	s = &Service{
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

func GetService() *Service {
	return s
}

func (s *Service) Close() {
	s.r.Close()
}

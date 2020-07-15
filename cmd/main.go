package main

import (
	"filesrv/conf"
	"filesrv/library/log"
	"filesrv/library/utils"
	"filesrv/service"
	"flag"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.NewLogger(conf.Conf.Log)
	service.NewService(conf.Conf)
	utils.QuitSignal(func() {
		log.GetLogger().Info("filesrv exit success")
	})
}

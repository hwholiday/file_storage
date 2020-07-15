package main

import (
	"filesrv/conf"
	"filesrv/library/log"
	"filesrv/library/utils"
	"filesrv/repositoty"
	"flag"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.NewLogger(conf.Conf.Log)
	repositoty.NewRepository(conf.Conf)
	utils.QuitSignal(func() {
		log.GetLogger().Info("filesrv exit success")
	})
}

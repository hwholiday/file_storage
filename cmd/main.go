package main

import (
	"filesrv/conf"
	"filesrv/library/log"
	"filesrv/library/utils"
	"flag"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.NewLogger(conf.Conf.Log)
	utils.QuitSignal(func() {
		log.GetLogger().Info("filesrv exit success")
	})
}

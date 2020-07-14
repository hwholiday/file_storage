package conf

import (
	"filesrv/library/database/minio"
	"filesrv/library/log"
	"flag"
	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = new(Config)
)

type Config struct {
	Development bool
	SnowFlakeId int64
	Log         *log.Options
	Minio       *minio.Config
	Http        *httpConf
	Grpc        *grpcConf
}

type httpConf struct {
	Port int
}
type grpcConf struct {
	Port int
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

func Init() error {
	return local()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

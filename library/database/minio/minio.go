package minio

import (
	"filesrv/library/log"
	"github.com/minio/minio-go"
	"go.uber.org/zap"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Location        string
	MaxBucket       int //设置MaxBucket个桶进行负载均衡 每次取一个桶
}

func NewMinio(c *Config) (mClient *minio.Client) {
	var err error
	if mClient, err = minio.New(c.Endpoint, c.AccessKeyID, c.SecretAccessKey, c.UseSSL); err != nil {
		log.GetLogger().Panic("[NewMinio] New", zap.Error(err))
		return
	}
	return
}

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
}

func NewMinio(c *Config) (mClient *minio.Client) {
	var err error
	if mClient, err = minio.New(c.Endpoint, c.AccessKeyID, c.SecretAccessKey, c.UseSSL); err != nil {
		log.GetLogger().Error("[NewMinio] New", zap.Error(err))
	}
	return
}

package repositoty

import (
	"context"
	"filesrv/common/storage/bucket"
	"filesrv/conf"
	m "filesrv/library/database/minio"
	mgo "filesrv/library/database/mongo"
	"filesrv/library/log"
	"filesrv/library/utils"
	"filesrv/repositoty/storage"
	"github.com/minio/minio-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository struct {
	c             *conf.Config
	mClient       *minio.Client
	mDb           *mongo.Client
	storageServer storage.Service
}

func NewRepository(c *conf.Config) (r *Repository) {
	r = &Repository{
		c:       c,
		mClient: m.NewMinio(c.Minio),
		mDb:     mgo.NewMongo(c.Mongo),
	}
	err := utils.NewWorker(c.SnowFlakeId)
	if err != nil {
		log.GetLogger().Panic("[NewRepository] NewWorker", zap.Error(err))
		return
	}
	bucket.NewBucket(r.mClient, c)
	storage.NewStorage(r.mClient)
	r.storageServer = storage.GetServer()
	return r
}

func (r *Repository) Close() {
	_ = r.mDb.Disconnect(context.Background())
}

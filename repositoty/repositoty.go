package repositoty

import (
	"context"
	"filesrv/conf"
	m "filesrv/library/database/minio"
	mgo "filesrv/library/database/mongo"
	"filesrv/repositoty/storage"
	"github.com/minio/minio-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	c             *conf.Config
	mClient       *minio.Client
	mDb           *mongo.Client
	storageServer *storage.StorageServer
}

func NewRepository(c *conf.Config) (r *Repository) {
	r = &Repository{
		c:       c,
		mClient: m.NewMinio(c.Minio),
		mDb:     mgo.NewMongo(c.Mongo),
	}
	r.storageServer = storage.NewStorage(r.mClient)
	return r
}

func (r *Repository) Close() {
	_ = r.mDb.Disconnect(context.Background())
}

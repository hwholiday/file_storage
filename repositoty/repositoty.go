package repositoty

import (
	"context"
	"filesrv/common/storage/bucket"
	"filesrv/conf"
	m "filesrv/library/database/minio"
	mgo "filesrv/library/database/mongo"
	"filesrv/repositoty/fileInfo"
	"filesrv/repositoty/storage"
	"github.com/minio/minio-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	c              *conf.Config
	mClient        *minio.Client
	mgo            *mongo.Client
	StorageServer  storage.Service
	FileInfoServer fileInfo.Service
}

func NewRepository(c *conf.Config) (r *Repository) {
	r = &Repository{
		c:       c,
		mClient: m.NewMinio(c.Minio),
		mgo:     mgo.NewMongo(c.Mongo),
	}
	bucket.NewBucket(r.mClient, c)
	r.StorageServer = storage.NewStorage(r.mClient)
	r.FileInfoServer = fileInfo.NewFileInfo(r.mgo, c)
	return r
}

func (r *Repository) Close() {
	_ = r.mgo.Disconnect(context.Background())
}

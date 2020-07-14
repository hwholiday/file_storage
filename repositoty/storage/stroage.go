package storage

import (
	"github.com/minio/minio-go"
)

type StorageServer struct {
	mClient *minio.Client
}

func NewStorage(mClient *minio.Client) (r *StorageServer) {
	r = &StorageServer{
		mClient: mClient,
	}
	return r
}

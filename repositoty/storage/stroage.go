package storage

import (
	"github.com/minio/minio-go"
)

type storage struct {
	mClient *minio.Client
}

type Service interface {
	GetFidAndBucketName() (int64, string)
	UpFileNotSlice(fid int64, bucketName string, data []byte) error
	GetFileNotSlice(fid int64, bucketName string) ([]byte, error)
	GetSliceFile(fid int64, bucketName string, start, end int64) ([]byte, error)
	DelFile(fid int64, bucketName string) error
}

func NewStorage(mClient *minio.Client) *storage {
	s := &storage{
		mClient: mClient,
	}
	return s
}

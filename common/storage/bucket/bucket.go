package bucket

import (
	"bytes"
	"filesrv/conf"
	"filesrv/library/log"
	"github.com/minio/minio-go"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var storageBucket *StorageBucket

type StorageBucket struct {
	mu        sync.Mutex
	prefix    string
	maxBucket int
	nextID    int
}

func NewBucket(client *minio.Client, c *conf.Config) {
	for i := 1; i <= c.Minio.MaxBucket; i++ {
		var bucketName bytes.Buffer
		bucketName.WriteString(c.AppName)
		bucketName.WriteString("-")
		bucketName.WriteString(strconv.Itoa(int(i)))
		ok, err := client.BucketExists(bucketName.String())
		if err != nil {
			log.GetLogger().Panic("[InitBucket] BucketExists", zap.Any("bucketName", bucketName.String()), zap.Error(err))
			return
		}
		if !ok {
			if err = client.MakeBucket(bucketName.String(), c.Minio.Location); err != nil {
				log.GetLogger().Panic("[InitBucket] MakeBucket", zap.Any("bucketName", bucketName.String()), zap.Any("location", c.Minio.Location), zap.Error(err))
				return
			}
		}
	}
	storageBucket = &StorageBucket{
		prefix:    c.AppName,
		maxBucket: c.Minio.MaxBucket,
	}
}

func GetStorageBucket() *StorageBucket {
	return storageBucket
}

func (s *StorageBucket) GetRandBucketName() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	var bucketName bytes.Buffer
	bucketName.WriteString(s.prefix)
	bucketName.WriteString("-")
	bucketName.WriteString(strconv.Itoa(int(s.nextID)))
	if s.nextID >= s.maxBucket {
		s.nextID = 0
	}
	return bucketName.String()
}

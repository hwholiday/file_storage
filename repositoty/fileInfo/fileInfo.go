package fileInfo

import (
	"filesrv/conf"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	s *storage
)

type storage struct {
	mgo *mongo.Database
}

type Service interface {
}

func NewStorage(mClient *mongo.Client, c *conf.Config) *storage {
	s = &storage{
		mgo: mClient.Database(c.Mongo.DataBase),
	}
	return s
}

func GetServer() Service {
	return s
}

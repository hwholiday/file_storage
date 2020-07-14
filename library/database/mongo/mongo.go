package mongo

import (
	"context"
	"filesrv/library/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Url string
}

func NewMongo(c *Config) (client *mongo.Client) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(c.Url)); err != nil {
		log.GetLogger().Panic("[NewMongo] Connect", zap.Error(err))
		return
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.GetLogger().Panic("[NewMongo] Ping", zap.Error(err))
		return
	}
	return
}

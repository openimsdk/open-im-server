package checks

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/log"
)

type MongoCheck struct {
	Mongo *config.Mongo
}

func CheckMongo(ctx context.Context, config *MongoCheck) error {
	mongoConfig := &mongoutil.Config{
		Uri:         config.Mongo.Uri,
		Address:     config.Mongo.Address,
		Database:    config.Mongo.Database,
		Username:    config.Mongo.Username,
		Password:    config.Mongo.Password,
		MaxPoolSize: config.Mongo.MaxPoolSize,
		MaxRetry:    0,
	}

	log.CInfo(ctx, "Checking MongoDB connection", "URI", mongoConfig.Uri, "Database", mongoConfig.Database)

	err := mongoutil.CheckMongo(ctx, mongoConfig)
	if err != nil {
		log.CInfo(ctx, "MongoDB connection failed", "error", err)
		return err
	}

	log.CInfo(ctx, "MongoDB connection established successfully")
	return nil
}

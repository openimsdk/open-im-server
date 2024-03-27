package checks

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/discovery/zookeeper"
	"github.com/openimsdk/tools/log"
)

type ZookeeperCheck struct {
	Zookeeper *config.Zookeeper
}

func checkZookeeper(ctx context.Context, config *ZookeeperCheck) error {
	zkServers := config.Zookeeper.ZkAddr
	schema := config.Zookeeper.Schema

	authOption := zookeeper.WithUserNameAndPassword(config.Zookeeper.Username, config.Zookeeper.Password)

	log.CInfo(ctx, "Checking Zookeeper connection", "Schema", schema, "ZkServers", zkServers)

	err := zookeeper.CheckZookeeper(ctx, zkServers, config.Zookeeper.Schema, authOption)
	if err != nil {
		log.CInfo(ctx, "Zookeeper connection failed", "error", err)
		return err
	}

	log.CInfo(ctx, "Zookeeper connection established successfully")
	return nil
}

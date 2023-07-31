package key

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/key"
	"google.golang.org/grpc"
)

type keyServer struct {
	authDatabase   controller.KeyDatabase
	keyRpcClient   key.KeyClient
	RegisterCenter discoveryregistry.SvcDiscoveryRegistry
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	keyDB := relation.NewKeyDB(db)
	database := controller.NewKeyDatabase(keyDB)
	key.RegisterKeyServer(server, &keyServer{
		authDatabase:   database,
		keyRpcClient:   nil,
		RegisterCenter: client,
	})
	return nil
}

func (k keyServer) GetKey(ctx context.Context, req *key.GetKeyReq) (*key.KeyResp, error) {
	return &key.KeyResp{}, nil
}
func (k keyServer) GetAllKey(ctx context.Context, req *key.GetAllKeyReq) (*key.GetAllKeyResp, error) {
	m := new(map[string]string)
	return &key.GetAllKeyResp{Keys: *m}, nil
}

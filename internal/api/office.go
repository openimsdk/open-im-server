package api

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/office"
	"google.golang.org/grpc"
)

func NewOffice(c discoveryregistry.SvcDiscoveryRegistry) *Office {
	conn, err := c.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImOfficeName)
	if err != nil {
		panic(err)
	}
	return &Office{conn: conn}
}

type Office struct {
	conn *grpc.ClientConn
}

func (o *Office) client(ctx context.Context) (office.OfficeClient, error) {
	return office.NewOfficeClient(o.conn), nil
}

package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"google.golang.org/grpc"
)

type MetaClient struct {
	// contains filtered or unexported fields
	client          discoveryregistry.SvcDiscoveryRegistry
	rpcRegisterName string
}

func NewMetaClient(client discoveryregistry.SvcDiscoveryRegistry, rpcRegisterName string, opts ...MetaClientOptions) *MetaClient {
	c := &MetaClient{
		client:          client,
		rpcRegisterName: rpcRegisterName,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type MetaClientOptions func(*MetaClient)

func (m *MetaClient) getConn(ctx context.Context) (*grpc.ClientConn, error) {
	return m.client.GetConn(ctx, m.rpcRegisterName)
}

type CommonUser interface {
	GetNickname() string
	GetFaceURL() string
	GetUserID() string
	GetEx() string
}

type CommonGroup interface {
	GetNickname() string
	GetFaceURL() string
	GetGroupID() string
	GetEx() string
}

package rpcclient

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"google.golang.org/grpc"
)

type MetaClient struct {
	// contains filtered or unexported fields
	client          discoveryregistry.SvcDiscoveryRegistry
	rpcRegisterName string
}

func NewMetaClient(client discoveryregistry.SvcDiscoveryRegistry, rpcRegisterName string) *MetaClient {
	return &MetaClient{
		client:          client,
		rpcRegisterName: rpcRegisterName,
	}
}

func (m *MetaClient) getConn() (*grpc.ClientConn, error) {
	return m.client.GetConn(m.rpcRegisterName)
}

type NotificationMsg struct {
	SendID         string
	RecvID         string
	Content        []byte
	MsgFrom        int32
	ContentType    int32
	SessionType    int32
	SenderNickname string
	SenderFaceURL  string
}

type CommonUser interface {
	GetNickname() string
	GetFaceURL() string
	GetUserID() string
	GetEx() string
}

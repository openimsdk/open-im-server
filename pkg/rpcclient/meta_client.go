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

type NotificationMsg struct {
	SendID         string
	RecvID         string
	Content        []byte //  sdkws.TipsComm
	MsgFrom        int32
	ContentType    int32
	SessionType    int32
	SenderNickname string
	SenderFaceURL  string
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

func (m *MetaClient) getRpcRegisterName() string {
	return m.rpcRegisterName
}

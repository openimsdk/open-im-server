package check

import (
	"Open_IM/pkg/common/config"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	"Open_IM/pkg/proto/msg"
	"context"
	"google.golang.org/grpc"
)

type MsgCheck struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func NewMsgCheck(zk discoveryRegistry.SvcDiscoveryRegistry) *MsgCheck {
	return &MsgCheck{zk: zk}
}

func (m *MsgCheck) getConn() (*grpc.ClientConn, error) {
	return m.zk.GetConn(config.Config.RpcRegisterName.OpenImMsgName)
}

func (m *MsgCheck) SendMsg(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).SendMsg(ctx, req)
	return resp, err
}

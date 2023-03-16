package check

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
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

func (m *MsgCheck) GetMaxAndMinSeq(ctx context.Context, req *sdkws.GetMaxAndMinSeqReq) (*sdkws.GetMaxAndMinSeqResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).GetMaxAndMinSeq(ctx, req)
	return resp, err
}

func (m *MsgCheck) PullMessageBySeqList(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).PullMessageBySeqs(ctx, req)
	return resp, err
}

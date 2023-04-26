package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type MsgClient struct {
	*MetaClient
}

func NewMsgClient(zk discoveryregistry.SvcDiscoveryRegistry) *MsgClient {
	return &MsgClient{NewMetaClient(zk, config.Config.RpcRegisterName.OpenImMsgName)}
}

func (m *MsgClient) SendMsg(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).SendMsg(ctx, req)
	return resp, err
}

func (m *MsgClient) GetMaxAndMinSeq(ctx context.Context, req *sdkws.GetMaxAndMinSeqReq) (*sdkws.GetMaxAndMinSeqResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).GetMaxAndMinSeq(ctx, req)
	return resp, err
}

func (m *MsgClient) PullMessageBySeqList(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).PullMessageBySeqs(ctx, req)
	return resp, err
}

func (c *MsgClient) Notification(ctx context.Context, notificationMsg *NotificationMsg) error {
	var err error
	var req msg.SendMsgReq
	var msg sdkws.MsgData
	var offlineInfo sdkws.OfflinePushInfo
	var title, desc, ex string
	var pushEnable, unReadCount bool
	msg.SendID = notificationMsg.SendID
	msg.RecvID = notificationMsg.RecvID
	msg.Content = notificationMsg.Content
	msg.MsgFrom = notificationMsg.MsgFrom
	msg.ContentType = notificationMsg.ContentType
	msg.SessionType = notificationMsg.SessionType
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(notificationMsg.SendID)
	msg.Options = make(map[string]bool, 7)
	msg.SenderNickname = notificationMsg.SenderNickname
	msg.SenderFaceURL = notificationMsg.SenderFaceURL
	utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, unReadCount)
	utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, pushEnable)
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ex = ex
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	_, err = c.SendMsg(ctx, &req)
	return err
}

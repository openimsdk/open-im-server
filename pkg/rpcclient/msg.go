package rpcclient

import (
	"context"
	"encoding/json"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"

	"github.com/golang/protobuf/proto"
	// "google.golang.org/protobuf/proto"
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

func (c *MsgClient) Notification(ctx context.Context, sendID, recvID string, contentType, sessionType int32, m proto.Message, cfg config.NotificationConf, opts ...utils.OptionsOpt) error {
	content, err := json.Marshal(m)
	if err != nil {
		return err
	}
	var req msg.SendMsgReq
	var msg sdkws.MsgData
	var offlineInfo sdkws.OfflinePushInfo
	var title, desc, ex string
	msg.SendID = sendID
	msg.RecvID = recvID
	if sessionType == constant.SuperGroupChatType {
		msg.GroupID = recvID
	}
	msg.Content = content
	msg.MsgFrom = constant.SysMsgType
	msg.ContentType = contentType
	msg.SessionType = sessionType
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(sendID)
	// msg.Options = make(map[string]bool, 7)
	// todo notification get sender name and face url
	// msg.SenderNickname, msg.SenderFaceURL, err = c.getFaceURLAndName(sendID)
	options := config.GetOptionsByNotification(cfg)
	options = utils.WithOptions(options, opts...)
	msg.Options = options
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ex = ex
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	_, err = c.SendMsg(ctx, &req)
	return err
}

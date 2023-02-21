package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	ExcludeContentType = []int{constant.HasReadReceipt, constant.GroupHasReadReceipt}
)

type Validator interface {
	validate(pb *msg.SendMsgReq) (bool, int32, string)
}

type MessageRevoked struct {
	RevokerID                   string `json:"revokerID"`
	RevokerRole                 int32  `json:"revokerRole"`
	ClientMsgID                 string `json:"clientMsgID"`
	RevokerNickname             string `json:"revokerNickname"`
	RevokeTime                  int64  `json:"revokeTime"`
	SourceMessageSendTime       int64  `json:"sourceMessageSendTime"`
	SourceMessageSendID         string `json:"sourceMessageSendID"`
	SourceMessageSenderNickname string `json:"sourceMessageSenderNickname"`
	SessionType                 int32  `json:"sessionType"`
	Seq                         uint32 `json:"seq"`
}
type MsgCallBackReq struct {
	SendID       string `json:"sendID"`
	RecvID       string `json:"recvID"`
	Content      string `json:"content"`
	SendTime     int64  `json:"sendTime"`
	MsgFrom      int32  `json:"msgFrom"`
	ContentType  int32  `json:"contentType"`
	SessionType  int32  `json:"sessionType"`
	PlatformID   int32  `json:"senderPlatformID"`
	MsgID        string `json:"msgID"`
	IsOnlineOnly bool   `json:"isOnlineOnly"`
}
type MsgCallBackResp struct {
	ErrCode         int32  `json:"errCode"`
	ErrMsg          string `json:"errMsg"`
	ResponseErrCode int32  `json:"responseErrCode"`
	ResponseResult  struct {
		ModifiedMsg string `json:"modifiedMsg"`
		Ext         string `json:"ext"`
	}
}

func (m *msgServer) userIsMuteAndIsAdminInGroup(ctx context.Context, groupID, userID string) (isMute bool, err error) {
	groupMemberInfo, err := m.Group.GetGroupMemberInfo(ctx, groupID, userID)
	if err != nil {
		return false, err
	}
	if groupMemberInfo.MuteEndTime >= time.Now().Unix() {
		return true, nil
	}
	return false, nil
}

// 如果禁言了，再看下是否群管理员
func (m *msgServer) groupIsMuted(ctx context.Context, groupID string, userID string) (bool, bool, error) {
	groupInfo, err := m.Group.GetGroupInfo(ctx, groupID)
	if err != nil {
		return false, false, err
	}

	if groupInfo.Status == constant.GroupStatusMuted {
		groupMemberInfo, err := m.Group.GetGroupMemberInfo(ctx, groupID, userID)
		if err != nil {
			return false, false, err
		}
		return true, groupMemberInfo.RoleLevel > constant.GroupOrdinaryUsers, nil
	}
	return false, false, nil
}

func (m *msgServer) GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error) {
	return m.GroupLocalCache.GetGroupMemberIDs(ctx, groupID)
}

func (m *msgServer) messageVerification(ctx context.Context, data *msg.SendMsgReq) ([]string, error) {
	switch data.MsgData.SessionType {
	case constant.SingleChatType:
		if utils.IsContain(data.MsgData.SendID, config.Config.Manager.AppManagerUid) {
			return nil, nil
		}
		if data.MsgData.ContentType <= constant.NotificationEnd && data.MsgData.ContentType >= constant.NotificationBegin {
			return nil, nil
		}
		black, err := m.black.IsBlocked(ctx, data.MsgData.SendID, data.MsgData.RecvID)
		if err != nil {
			return nil, err
		}
		if black {
			return nil, constant.ErrBlockedByPeer.Wrap()
		}
		if *config.Config.MessageVerify.FriendVerify {
			friend, err := m.friend.IsFriend(ctx, data.MsgData.SendID, data.MsgData.RecvID)
			if err != nil {
				return nil, err
			}
			if !friend {
				return nil, constant.ErrNotPeersFriend.Wrap()
			}
			return nil, nil
		}
		return nil, nil
	case constant.GroupChatType:
		if utils.IsContain(data.MsgData.SendID, config.Config.Manager.AppManagerUid) {
			return nil, nil
		}
		userIDList, err := m.GetGroupMemberIDs(ctx, data.MsgData.GroupID)
		if err != nil {
			return nil, err
		}
		if utils.IsContain(data.MsgData.SendID, config.Config.Manager.AppManagerUid) {
			return userIDList, nil
		}
		if data.MsgData.ContentType <= constant.NotificationEnd && data.MsgData.ContentType >= constant.NotificationBegin {
			return userIDList, nil
		}
		if !utils.IsContain(data.MsgData.SendID, userIDList) {
			return nil, constant.ErrNotInGroupYet.Wrap()
		}
		isMute, err := m.userIsMuteAndIsAdminInGroup(ctx, data.MsgData.GroupID, data.MsgData.SendID)
		if err != nil {
			return nil, err
		}
		if isMute {
			return nil, constant.ErrMutedInGroup.Wrap()
		}

		isMute, isAdmin, err := m.groupIsMuted(ctx, data.MsgData.GroupID, data.MsgData.SendID)
		if err != nil {
			return nil, err
		}
		if isAdmin {
			return userIDList, nil
		}

		if isMute {
			return nil, constant.ErrMutedGroup.Wrap()
		}
		return userIDList, nil
	case constant.SuperGroupChatType:
		groupInfo, err := m.Group.GetGroupInfo(ctx, data.MsgData.GroupID)
		if err != nil {
			return nil, err
		}
		if data.MsgData.ContentType == constant.AdvancedRevoke {
			revokeMessage := new(MessageRevoked)
			err := utils.JsonStringToStruct(string(data.MsgData.Content), revokeMessage)
			if err != nil {
				return nil, constant.ErrArgs.Wrap()
			}

			if revokeMessage.RevokerID != revokeMessage.SourceMessageSendID {
				resp, err := m.MsgInterface.GetSuperGroupMsgBySeqs(ctx, data.MsgData.GroupID, []int64{int64(revokeMessage.Seq)})
				if err != nil {
					return nil, err
				}
				if resp[0].ClientMsgID == revokeMessage.ClientMsgID && resp[0].Seq == int64(revokeMessage.Seq) {
					revokeMessage.SourceMessageSendTime = resp[0].SendTime
					revokeMessage.SourceMessageSenderNickname = resp[0].SenderNickname
					revokeMessage.SourceMessageSendID = resp[0].SendID
					data.MsgData.Content = []byte(utils.StructToJsonString(revokeMessage))
				} else {
					return nil, constant.ErrData.Wrap("MsgData")
				}
			}
		}
		if groupInfo.GroupType == constant.SuperGroup {
			return nil, nil
		}

		userIDList, err := m.GetGroupMemberIDs(ctx, data.MsgData.GroupID)
		if err != nil {
			return nil, err
		}
		if utils.IsContain(data.MsgData.SendID, config.Config.Manager.AppManagerUid) {
			return nil, nil
		}
		if data.MsgData.ContentType <= constant.NotificationEnd && data.MsgData.ContentType >= constant.NotificationBegin {
			return userIDList, nil
		} else {
			if !utils.IsContain(data.MsgData.SendID, userIDList) {
				return nil, constant.ErrNotInGroupYet.Wrap()
			}
		}
		isMute, err := m.userIsMuteAndIsAdminInGroup(ctx, data.MsgData.GroupID, data.MsgData.SendID)
		if err != nil {
			return nil, err
		}
		if isMute {
			return nil, constant.ErrMutedInGroup.Wrap()
		}

		isMute, isAdmin, err := m.groupIsMuted(ctx, data.MsgData.GroupID, data.MsgData.SendID)
		if err != nil {
			return nil, err
		}
		if isAdmin {
			return userIDList, nil
		}
		if isMute {
			return nil, constant.ErrMutedGroup.Wrap()
		}
		return userIDList, nil

	default:
		return nil, nil
	}
}
func (m *msgServer) encapsulateMsgData(msg *sdkws.MsgData) {
	msg.ServerMsgID = GetMsgID(msg.SendID)
	msg.SendTime = utils.GetCurrentTimestampByMill()
	switch msg.ContentType {
	case constant.Text:
		fallthrough
	case constant.Picture:
		fallthrough
	case constant.Voice:
		fallthrough
	case constant.Video:
		fallthrough
	case constant.File:
		fallthrough
	case constant.AtText:
		fallthrough
	case constant.Merger:
		fallthrough
	case constant.Card:
		fallthrough
	case constant.Location:
		fallthrough
	case constant.Custom:
		fallthrough
	case constant.Quote:
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, true)
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, true)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderSync, true)
	case constant.Revoke:
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
	case constant.HasReadReceipt:
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
	case constant.Typing:
		utils.SetSwitchFromOptions(msg.Options, constant.IsHistory, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsPersistent, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderSync, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
	}
}

func GetMsgID(sendID string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return utils.Md5(t + "-" + sendID + "-" + strconv.Itoa(rand.Int()))
}

func (m *msgServer) modifyMessageByUserMessageReceiveOpt(ctx context.Context, userID, sourceID string, sessionType int, pb *msg.SendMsgReq) (bool, error) {
	opt, err := m.User.GetUserGlobalMsgRecvOpt(ctx, userID)
	if err != nil {
		return false, err
	}
	switch opt {
	case constant.ReceiveMessage:
	case constant.NotReceiveMessage:
		return false, nil
	case constant.ReceiveNotNotifyMessage:
		if pb.MsgData.Options == nil {
			pb.MsgData.Options = make(map[string]bool, 10)
		}
		utils.SetSwitchFromOptions(pb.MsgData.Options, constant.IsOfflinePush, false)
		return true, nil
	}
	conversationID := utils.GetConversationIDBySessionType(sourceID, sessionType)
	singleOpt, err := m.Conversation.GetSingleConversationRecvMsgOpt(ctx, userID, conversationID)
	if err != nil {
		return false, err
	}
	switch singleOpt {
	case constant.ReceiveMessage:
		return true, nil
	case constant.NotReceiveMessage:
		if utils.IsContainInt(int(pb.MsgData.ContentType), ExcludeContentType) {
			return true, nil
		}
		return false, nil
	case constant.ReceiveNotNotifyMessage:
		if pb.MsgData.Options == nil {
			pb.MsgData.Options = make(map[string]bool, 10)
		}
		utils.SetSwitchFromOptions(pb.MsgData.Options, constant.IsOfflinePush, false)
		return true, nil
	}

	return true, nil
}

func valueCopy(pb *msg.SendMsgReq) *msg.SendMsgReq {
	offlinePushInfo := sdkws.OfflinePushInfo{}
	if pb.MsgData.OfflinePushInfo != nil {
		offlinePushInfo = *pb.MsgData.OfflinePushInfo
	}
	msgData := sdkws.MsgData{}
	msgData = *pb.MsgData
	msgData.OfflinePushInfo = &offlinePushInfo

	options := make(map[string]bool, 10)
	for key, value := range pb.MsgData.Options {
		options[key] = value
	}
	msgData.Options = options
	return &msg.SendMsgReq{MsgData: &msgData}
}

func (m *msgServer) sendMsgToGroupOptimization(ctx context.Context, list []string, groupPB *msg.SendMsgReq, wg *sync.WaitGroup) error {
	msgToMQGroup := msg.MsgDataToMQ{OperationID: tracelog.GetOperationID(ctx), MsgData: groupPB.MsgData}
	tempOptions := make(map[string]bool, 1)
	for k, v := range groupPB.MsgData.Options {
		tempOptions[k] = v
	}
	for _, v := range list {
		groupPB.MsgData.RecvID = v
		options := make(map[string]bool, 1)
		for k, v := range tempOptions {
			options[k] = v
		}
		groupPB.MsgData.Options = options
		isSend, err := m.modifyMessageByUserMessageReceiveOpt(ctx, v, groupPB.MsgData.GroupID, constant.GroupChatType, groupPB)
		if err != nil {
			wg.Done()
			return err
		}
		if isSend {
			if v == "" || groupPB.MsgData.SendID == "" {
				return constant.ErrArgs.Wrap("userID or groupPB.MsgData.SendID is empty")
			}
			err := m.MsgInterface.MsgToMQ(ctx, v, &msgToMQGroup)
			if err != nil {
				wg.Done()
				return err
			}
		}
	}
	wg.Done()
	return nil
}

/*
** description("").
** copyright('Open_IM,www.Open_IM.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/21 15:29).
 */
package gate

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbRtc "Open_IM/pkg/proto/rtc"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"github.com/golang/protobuf/proto"
)

type Req struct {
	ReqIdentifier int32  `json:"reqIdentifier" validate:"required"`
	Token         string `json:"token" `
	SendID        string `json:"sendID" validate:"required"`
	OperationID   string `json:"operationID" validate:"required"`
	MsgIncr       string `json:"msgIncr" validate:"required"`
	Data          []byte `json:"data"`
}
type Resp struct {
	ReqIdentifier int32  `json:"reqIdentifier"`
	MsgIncr       string `json:"msgIncr"`
	OperationID   string `json:"operationID"`
	ErrCode       int32  `json:"errCode"`
	ErrMsg        string `json:"errMsg"`
	Data          []byte `json:"data"`
}

type SeqData struct {
	SeqBegin int64 `mapstructure:"seqBegin" validate:"required"`
	SeqEnd   int64 `mapstructure:"seqEnd" validate:"required"`
}
type MsgData struct {
	PlatformID  int32                  `mapstructure:"platformID" validate:"required"`
	SessionType int32                  `mapstructure:"sessionType" validate:"required"`
	MsgFrom     int32                  `mapstructure:"msgFrom" validate:"required"`
	ContentType int32                  `mapstructure:"contentType" validate:"required"`
	RecvID      string                 `mapstructure:"recvID" validate:"required"`
	ForceList   []string               `mapstructure:"forceList"`
	Content     string                 `mapstructure:"content" validate:"required"`
	Options     map[string]interface{} `mapstructure:"options" validate:"required"`
	ClientMsgID string                 `mapstructure:"clientMsgID" validate:"required"`
	OfflineInfo map[string]interface{} `mapstructure:"offlineInfo" validate:"required"`
	Ext         map[string]interface{} `mapstructure:"ext"`
}
type MaxSeqResp struct {
	MaxSeq int64 `json:"maxSeq"`
}
type PullMessageResp struct {
}
type SeqListData struct {
	SeqList []int64 `mapstructure:"seqList" validate:"required"`
}

func (ws *WServer) argsValidate(m *Req, r int32) (isPass bool, errCode int32, errMsg string, returnData interface{}) {
	switch r {
	case constant.WSSendMsg:
		data := open_im_sdk.MsgData{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			log.ErrorByKv("Decode Data struct  err", "", "err", err.Error(), "reqIdentifier", r)
			return false, 203, err.Error(), nil
		}
		if err := validate.Struct(data); err != nil {
			log.ErrorByKv("data args validate  err", "", "err", err.Error(), "reqIdentifier", r)
			return false, 204, err.Error(), nil

		}
		return true, 0, "", data
	case constant.WSSendSignalMsg:
		data := pbRtc.SignalReq{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			log.ErrorByKv("Decode Data struct  err", "", "err", err.Error(), "reqIdentifier", r)
			return false, 203, err.Error(), nil
		}
		if err := validate.Struct(data); err != nil {
			log.ErrorByKv("data args validate  err", "", "err", err.Error(), "reqIdentifier", r)
			return false, 204, err.Error(), nil

		}
		return true, 0, "", &data
	case constant.WSPullMsgBySeqList:
		data := open_im_sdk.PullMessageBySeqListReq{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			log.ErrorByKv("Decode Data struct  err", "", "err", err.Error(), "reqIdentifier", r)
			return false, 203, err.Error(), nil
		}
		if err := validate.Struct(data); err != nil {
			log.ErrorByKv("data args validate  err", "", "err", err.Error(), "reqIdentifier", r)
			return false, 204, err.Error(), nil

		}
		return true, 0, "", data

	default:
	}

	return false, 204, "args err", nil

	//b := bytes.NewBuffer(m.Data)
	//dec := gob.NewDecoder(b)
	//err := dec.Decode(&data)
	//if err != nil {
	//	log.ErrorByKv("Decode Data struct  err", "", "err", err.Error(), "reqIdentifier", r)
	//	return false, 203, err.Error(), nil
	//}
	//if err := mapstructure.WeakDecode(m.Data, &data); err != nil {
	//	log.ErrorByKv("map to Data struct  err", "", "err", err.Error(), "reqIdentifier", r)
	//	return false, 203, err.Error(), nil
	//} else

}

//func (ws *WServer) signalMessageAssemble(s *open_im_sdk.SignalReq, operationID string) (isPass bool, errCode int32, errMsg string, r *open_im_sdk.SignalResp, msgData *open_im_sdk.MsgData) {
//	var msg open_im_sdk.MsgData
//	var resp open_im_sdk.SignalResp
//	media := open_im_media.NewMedia()
//	msg.MsgFrom = constant.UserMsgType
//	msg.ContentType = constant.SignalingNotification
//	reqData, e := proto.Marshal(s)
//	if e != nil {
//		return false, 201, e.Error(), nil, nil
//	}
//	msg.Content = reqData
//	msg.CreateTime = utils.GetCurrentTimestampByMill()
//	options := make(map[string]bool, 6)
//	utils.SetSwitchFromOptions(options, constant.IsHistory, false)
//	utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
//	utils.SetSwitchFromOptions(options, constant.IsSenderSync, true)
//	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
//	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
//	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
//	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, true)
//	msg.Options = options
//	switch payload := s.Payload.(type) {
//	case *open_im_sdk.SignalReq_Invite:
//		token, liveURL, err2 := media.GetJoinToken(payload.Invite.Invitation.RoomID, payload.Invite.Invitation.InviterUserID, operationID, payload.Invite.Participant)
//		if err2 != nil {
//			return false, 202, err2.Error(), nil, nil
//		}
//		invite := open_im_sdk.SignalResp_Invite{&open_im_sdk.SignalInviteReply{
//			Token:   token,
//			RoomID:  payload.Invite.Invitation.RoomID,
//			LiveURL: liveURL,
//		}}
//		resp.Payload = &invite
//		msg.SenderPlatformID = payload.Invite.Invitation.PlatformID
//		msg.SessionType = payload.Invite.Invitation.SessionType
//		msg.OfflinePushInfo = payload.Invite.OfflinePushInfo
//		msg.SendID = payload.Invite.Invitation.InviterUserID
//		if len(payload.Invite.Invitation.InviteeUserIDList) > 0 {
//			msg.RecvID = payload.Invite.Invitation.InviteeUserIDList[0]
//		} else {
//			return false, 203, errors.New("InviteeUserIDList is null").Error(), nil, nil
//		}
//		msg.ClientMsgID = utils.GetMsgID(payload.Invite.Invitation.InviterUserID)
//		return true, 0, "", &resp, &msg
//	case *open_im_sdk.SignalReq_InviteInGroup:
//		token, liveURL, err2 := media.GetJoinToken(payload.InviteInGroup.Invitation.RoomID, payload.InviteInGroup.Invitation.InviterUserID, operationID, payload.InviteInGroup.Participant)
//		if err2 != nil {
//			return false, 204, err2.Error(), nil, nil
//		}
//		inviteGroup := open_im_sdk.SignalResp_InviteInGroup{&open_im_sdk.SignalInviteInGroupReply{
//			RoomID:  payload.InviteInGroup.Invitation.RoomID,
//			Token:   token,
//			LiveURL: liveURL,
//		}}
//		resp.Payload = &inviteGroup
//		msg.SenderPlatformID = payload.InviteInGroup.Invitation.PlatformID
//		msg.SessionType = payload.InviteInGroup.Invitation.SessionType
//		msg.OfflinePushInfo = payload.InviteInGroup.OfflinePushInfo
//		msg.SendID = payload.InviteInGroup.Invitation.InviterUserID
//		if len(payload.InviteInGroup.Invitation.InviteeUserIDList) > 0 {
//			msg.GroupID = payload.InviteInGroup.Invitation.GroupID
//		} else {
//			return false, 205, errors.New("InviteeUserIDList is null").Error(), nil, nil
//		}
//		msg.ClientMsgID = utils.GetMsgID(payload.InviteInGroup.Invitation.InviterUserID)
//
//		return true, 0, "", &resp, &msg
//	case *open_im_sdk.SignalReq_Cancel:
//		cancel := open_im_sdk.SignalResp_Cancel{&open_im_sdk.SignalCancelReply{}}
//		resp.Payload = &cancel
//		msg.OfflinePushInfo = payload.Cancel.OfflinePushInfo
//		msg.SendID = payload.Cancel.Invitation.InviterUserID
//		msg.SenderPlatformID = payload.Cancel.Invitation.PlatformID
//		msg.SessionType = payload.Cancel.Invitation.SessionType
//		if len(payload.Cancel.Invitation.InviteeUserIDList) > 0 {
//			switch payload.Cancel.Invitation.SessionType {
//			case constant.SingleChatType:
//				msg.RecvID = payload.Cancel.Invitation.InviteeUserIDList[0]
//			case constant.GroupChatType:
//				msg.GroupID = payload.Cancel.Invitation.GroupID
//			}
//		} else {
//			return false, 206, errors.New("InviteeUserIDList is null").Error(), nil, nil
//		}
//		msg.ClientMsgID = utils.GetMsgID(payload.Cancel.OpUserID)
//		return true, 0, "", &resp, &msg
//	case *open_im_sdk.SignalReq_Accept:
//		token, liveURL, err2 := media.GetJoinToken(payload.Accept.Invitation.RoomID, payload.Accept.OpUserID, operationID, payload.Accept.Participant)
//		if err2 != nil {
//			return false, 207, err2.Error(), nil, nil
//		}
//		accept := open_im_sdk.SignalResp_Accept{&open_im_sdk.SignalAcceptReply{
//			Token:   token,
//			LiveURL: liveURL,
//			RoomID:  payload.Accept.Invitation.RoomID,
//		}}
//		resp.Payload = &accept
//		msg.OfflinePushInfo = payload.Accept.OfflinePushInfo
//		msg.SendID = payload.Accept.OpUserID
//		msg.SenderPlatformID = payload.Accept.Invitation.PlatformID
//		msg.SessionType = payload.Accept.Invitation.SessionType
//		if len(payload.Accept.Invitation.InviteeUserIDList) > 0 {
//			switch payload.Accept.Invitation.SessionType {
//			case constant.SingleChatType:
//				msg.RecvID = payload.Accept.Invitation.InviterUserID
//			case constant.GroupChatType:
//				msg.GroupID = payload.Accept.Invitation.GroupID
//			}
//		} else {
//			return false, 208, errors.New("InviteeUserIDList is null").Error(), nil, nil
//		}
//		msg.ClientMsgID = utils.GetMsgID(payload.Accept.OpUserID)
//		return true, 0, "", &resp, &msg
//	case *open_im_sdk.SignalReq_HungUp:
//	case *open_im_sdk.SignalReq_Reject:
//		reject := open_im_sdk.SignalResp_Reject{&open_im_sdk.SignalRejectReply{}}
//		resp.Payload = &reject
//		msg.OfflinePushInfo = payload.Reject.OfflinePushInfo
//		msg.SendID = payload.Reject.OpUserID
//		msg.SenderPlatformID = payload.Reject.Invitation.PlatformID
//		msg.SessionType = payload.Reject.Invitation.SessionType
//		if len(payload.Reject.Invitation.InviteeUserIDList) > 0 {
//			switch payload.Reject.Invitation.SessionType {
//			case constant.SingleChatType:
//				msg.RecvID = payload.Reject.Invitation.InviterUserID
//			case constant.GroupChatType:
//				msg.GroupID = payload.Reject.Invitation.GroupID
//			}
//		} else {
//			return false, 209, errors.New("InviteeUserIDList is null").Error(), nil, nil
//		}
//		msg.ClientMsgID = utils.GetMsgID(payload.Reject.OpUserID)
//		return true, 0, "", &resp, &msg
//	}
//	return false, 210, errors.New("InviteeUserIDList is null").Error(), nil, nil
//}

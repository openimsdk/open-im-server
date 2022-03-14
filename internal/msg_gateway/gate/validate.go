/*
** description("").
** copyright('Open_IM,www.Open_IM.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/21 15:29).
 */
package gate

import (
	"Open_IM/internal/msg_gateway/gate/open_im_media"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"errors"
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
		data := open_im_sdk.SignalReq{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			log.ErrorByKv("Decode Data struct  err", "", "err", err.Error(), "reqIdentifier", r)
			return false, 203, err.Error(), nil
		}
		if err := validate.Struct(data); err != nil {
			log.ErrorByKv("data args validate  err", "", "err", err.Error(), "reqIdentifier", r)
			return false, 204, err.Error(), nil

		}
		return true, 0, "", data
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

func (ws *WServer) signalMessageAssemble(s *open_im_sdk.SignalReq) (isPass bool, errCode int32, errMsg string, r *open_im_sdk.SignalResp, msgData *open_im_sdk.MsgData) {
	var msg open_im_sdk.MsgData
	var resp open_im_sdk.SignalResp
	media := open_im_media.NewMedia()
	msg.MsgFrom = constant.UserMsgType
	msg.ContentType = constant.SignalingNotification
	reqData, e := proto.Marshal(s)
	if e != nil {
		return false, 201, e.Error(), nil, nil
	}
	msg.Content = reqData
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	options := make(map[string]bool, 6)
	utils.SetSwitchFromOptions(options, constant.IsHistory, false)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, true)
	msg.Options = options
	switch payload := s.Payload.(type) {
	case *open_im_sdk.SignalReq_Invite:
		_, err := media.CreateRoom(payload.Invite.Invitation.RoomID)
		if err != nil {
			return false, 201, err.Error(), nil, nil

		}
		token, err2 := media.GetJoinToken(payload.Invite.Invitation.RoomID, payload.Invite.Invitation.InviterUserID)
		if err2 != nil {
			return false, 201, err2.Error(), nil, nil
		}
		invite := open_im_sdk.SignalResp_Invite{&open_im_sdk.SignalInviteReply{
			Token:   token,
			LiveURL: media.GetUrl(),
		}}
		resp.Payload = &invite
		msg.SenderPlatformID = payload.Invite.Invitation.PlatformID
		msg.SessionType = payload.Invite.Invitation.SessionType
		msg.OfflinePushInfo = payload.Invite.OfflinePushInfo
		msg.SendID = payload.Invite.Invitation.InviterUserID
		if len(payload.Invite.Invitation.InviteeUserIDList) > 0 {
			msg.RecvID = payload.Invite.Invitation.InviteeUserIDList[0]
		} else {
			return false, 201, errors.New("InviteeUserIDList is null").Error(), nil, nil
		}
		msg.ClientMsgID = utils.GetMsgID(payload.Invite.Invitation.InviterUserID)
		return true, 0, "", &resp, &msg
	case *open_im_sdk.SignalReq_InviteInGroup:
		_, err := media.CreateRoom(payload.InviteInGroup.Invitation.RoomID)
		if err != nil {
			return false, 201, err.Error(), nil, nil

		}
		token, err2 := media.GetJoinToken(payload.InviteInGroup.Invitation.RoomID, payload.InviteInGroup.Invitation.InviterUserID)
		if err2 != nil {
			return false, 201, err2.Error(), nil, nil
		}
		inviteGroup := open_im_sdk.SignalResp_InviteInGroup{&open_im_sdk.SignalInviteInGroupReply{
			RoomID:  payload.InviteInGroup.Invitation.RoomID,
			Token:   token,
			LiveURL: media.GetUrl(),
		}}
		resp.Payload = &inviteGroup
		msg.SenderPlatformID = payload.InviteInGroup.Invitation.PlatformID
		msg.SessionType = payload.InviteInGroup.Invitation.SessionType
		msg.OfflinePushInfo = payload.InviteInGroup.OfflinePushInfo
		msg.SendID = payload.InviteInGroup.Invitation.InviterUserID
		if len(payload.InviteInGroup.Invitation.InviteeUserIDList) > 0 {
			msg.GroupID = payload.InviteInGroup.Invitation.GroupID
		} else {
			return false, 201, errors.New("InviteeUserIDList is null").Error(), nil, nil
		}
		msg.ClientMsgID = utils.GetMsgID(payload.InviteInGroup.Invitation.InviterUserID)

		return true, 0, "", &resp, &msg
	case *open_im_sdk.SignalReq_Cancel:
		cancel := open_im_sdk.SignalResp_Cancel{&open_im_sdk.SignalCancelReply{}}
		resp.Payload = &cancel
		msg.OfflinePushInfo = payload.Cancel.Invitation.OfflinePushInfo
		msg.SendID = payload.Cancel.Invitation.Invitation.InviterUserID
		msg.SenderPlatformID = payload.Cancel.Invitation.Invitation.PlatformID
		msg.SessionType = payload.Cancel.Invitation.Invitation.SessionType
		if len(payload.Cancel.Invitation.Invitation.InviteeUserIDList) > 0 {
			switch payload.Cancel.Invitation.Invitation.SessionType {
			case constant.SingleChatType:
				msg.RecvID = payload.Cancel.Invitation.Invitation.InviteeUserIDList[0]
			case constant.GroupChatType:
				msg.GroupID = payload.Cancel.Invitation.Invitation.GroupID
			}
		} else {
			return false, 201, errors.New("InviteeUserIDList is null").Error(), nil, nil
		}
		msg.ClientMsgID = utils.GetMsgID(payload.Cancel.InviterUserID)
		return true, 0, "", &resp, &msg
	case *open_im_sdk.SignalReq_Accept:
		token, err2 := media.GetJoinToken(payload.Accept.Invitation.Invitation.RoomID, payload.Accept.InviteeUserID)
		if err2 != nil {
			return false, 201, err2.Error(), nil, nil
		}
		cancel := open_im_sdk.SignalResp_Accept{&open_im_sdk.SignalAcceptReply{
			Token:   token,
			LiveURL: media.GetUrl(),
			RoomID:  payload.Accept.Invitation.Invitation.RoomID,
		}}
		resp.Payload = &cancel
		msg.OfflinePushInfo = payload.Accept.Invitation.OfflinePushInfo
		msg.SendID = payload.Accept.InviteeUserID
		msg.SenderPlatformID = payload.Accept.Invitation.Invitation.PlatformID
		msg.SessionType = payload.Accept.Invitation.Invitation.SessionType
		if len(payload.Accept.Invitation.Invitation.InviteeUserIDList) > 0 {
			switch payload.Accept.Invitation.Invitation.SessionType {
			case constant.SingleChatType:
				msg.RecvID = payload.Accept.Invitation.Invitation.InviterUserID
			case constant.GroupChatType:
				msg.GroupID = payload.Accept.Invitation.Invitation.GroupID
			}
		} else {
			return false, 201, errors.New("InviteeUserIDList is null").Error(), nil, nil
		}
		msg.ClientMsgID = utils.GetMsgID(payload.Accept.InviteeUserID)
		return true, 0, "", &resp, &msg
	case *open_im_sdk.SignalReq_HungUp:
	case *open_im_sdk.SignalReq_Reject:
		cancel := open_im_sdk.SignalResp_Reject{&open_im_sdk.SignalRejectReply{}}
		resp.Payload = &cancel
		msg.OfflinePushInfo = payload.Reject.Invitation.OfflinePushInfo
		msg.SendID = payload.Reject.InviteeUserID
		msg.SenderPlatformID = payload.Reject.Invitation.Invitation.PlatformID
		msg.SessionType = payload.Reject.Invitation.Invitation.SessionType
		if len(payload.Reject.Invitation.Invitation.InviteeUserIDList) > 0 {
			switch payload.Reject.Invitation.Invitation.SessionType {
			case constant.SingleChatType:
				msg.RecvID = payload.Reject.Invitation.Invitation.InviterUserID
			case constant.GroupChatType:
				msg.GroupID = payload.Reject.Invitation.Invitation.GroupID
			}
		} else {
			return false, 201, errors.New("InviteeUserIDList is null").Error(), nil, nil
		}
		msg.ClientMsgID = utils.GetMsgID(payload.Reject.InviteeUserID)
		return true, 0, "", &resp, &msg
	}
	return false, 201, errors.New("InviteeUserIDList is null").Error(), nil, nil
}

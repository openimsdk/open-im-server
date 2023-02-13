package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tokenverify"
	cacheRpc "Open_IM/pkg/proto/cache"
	"Open_IM/pkg/proto/msg"
	pbPush "Open_IM/pkg/proto/push"
	pbRelay "Open_IM/pkg/proto/relay"
	sdkws "Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	go_redis "github.com/go-redis/redis/v8"
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

func isMessageHasReadEnabled(pb *msg.SendMsgReq) (bool, int32, string) {
	switch pb.MsgData.ContentType {
	case constant.HasReadReceipt:
		if config.Config.SingleMessageHasReadReceiptEnable {
			return true, 0, ""
		} else {
			return false, constant.ErrMessageHasReadDisable.ErrCode, constant.ErrMessageHasReadDisable.ErrMsg
		}
	case constant.GroupHasReadReceipt:
		if config.Config.GroupMessageHasReadReceiptEnable {
			return true, 0, ""
		} else {
			return false, constant.ErrMessageHasReadDisable.ErrCode, constant.ErrMessageHasReadDisable.ErrMsg
		}
	}
	return true, 0, ""
}

func userIsMuteAndIsAdminInGroup(ctx context.Context, groupID, userID string) (isMute bool, isAdmin bool, err error) {
	groupMemberInfo, err := rocksCache.GetGroupMemberInfoFromCache(ctx, groupID, userID)
	if err != nil {
		return false, false, utils.Wrap(err, "")
	}

	if groupMemberInfo.MuteEndTime.Unix() >= time.Now().Unix() {
		return true, groupMemberInfo.RoleLevel > constant.GroupOrdinaryUsers, nil
	}
	return false, groupMemberInfo.RoleLevel > constant.GroupOrdinaryUsers, nil
}

func groupIsMuted(ctx context.Context, groupID string) (bool, error) {
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
	if err != nil {
		return false, utils.Wrap(err, "GetGroupInfoFromCache failed")
	}
	if groupInfo.Status == constant.GroupStatusMuted {
		return true, nil
	}
	return false, nil
}

func (rpc *msgServer) messageVerification(ctx context.Context, data *pbChat.SendMsgReq) (bool, int32, string, []string) {
	switch data.MsgData.SessionType {
	case constant.SingleChatType:
		if utils.IsContain(data.MsgData.SendID, config.Config.Manager.AppManagerUid) {
			return true, 0, "", nil
		}
		if data.MsgData.ContentType <= constant.NotificationEnd && data.MsgData.ContentType >= constant.NotificationBegin {
			return true, 0, "", nil
		}
		log.NewDebug(data.OperationID, *config.Config.MessageVerify.FriendVerify)
		reqGetBlackIDListFromCache := &cacheRpc.GetBlackIDListFromCacheReq{UserID: data.MsgData.RecvID, OperationID: data.OperationID}
		etcdConn, err := rpc.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImCacheName)
		if err != nil {
			errMsg := data.OperationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(data.OperationID, errMsg)
			return true, 0, "", nil
		}

		cacheClient := cacheRpc.NewCacheClient(etcdConn)
		cacheResp, err := cacheClient.GetBlackIDListFromCache(context.Background(), reqGetBlackIDListFromCache)
		if err != nil {
			log.NewError(data.OperationID, "GetBlackIDListFromCache rpc call failed ", err.Error())
		} else {
			if cacheResp.CommonResp.ErrCode != 0 {
				log.NewError(data.OperationID, "GetBlackIDListFromCache rpc logic call failed ", cacheResp.String())
			} else {
				if utils.IsContain(data.MsgData.SendID, cacheResp.UserIDList) {
					return false, 600, "in black list", nil
				}
			}
		}
		log.NewDebug(data.OperationID, *config.Config.MessageVerify.FriendVerify)
		if *config.Config.MessageVerify.FriendVerify {
			reqGetFriendIDListFromCache := &cacheRpc.GetFriendIDListFromCacheReq{UserID: data.MsgData.RecvID, OperationID: data.OperationID}
			etcdConn, err := rpc.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImCacheName)
			if err != nil {
				errMsg := data.OperationID + "getcdv3.GetDefaultConn == nil"
				log.NewError(data.OperationID, errMsg)
				return true, 0, "", nil
			}
			cacheClient := cacheRpc.NewCacheClient(etcdConn)
			cacheResp, err := cacheClient.GetFriendIDListFromCache(context.Background(), reqGetFriendIDListFromCache)
			if err != nil {
				log.NewError(data.OperationID, "GetFriendIDListFromCache rpc call failed ", err.Error())
			} else {
				if cacheResp.CommonResp.ErrCode != 0 {
					log.NewError(data.OperationID, "GetFriendIDListFromCache rpc logic call failed ", cacheResp.String())
				} else {
					if !utils.IsContain(data.MsgData.SendID, cacheResp.UserIDList) {
						return false, 601, "not friend", nil
					}
				}
			}
			return true, 0, "", nil
		} else {
			return true, 0, "", nil
		}
	case constant.GroupChatType:
		userIDList, err := utils.GetGroupMemberUserIDList(ctx, data.MsgData.GroupID, data.OperationID)
		if err != nil {
			errMsg := data.OperationID + err.Error()
			log.NewError(data.OperationID, errMsg)
			return false, 201, errMsg, nil
		}
		if tokenverify.IsManagerUserID(data.MsgData.SendID) {
			return true, 0, "", userIDList
		}
		if data.MsgData.ContentType <= constant.NotificationEnd && data.MsgData.ContentType >= constant.NotificationBegin {
			return true, 0, "", userIDList
		} else {
			if !utils.IsContain(data.MsgData.SendID, userIDList) {
				//return returnMsg(&replay, pb, 202, "you are not in group", "", 0)
				return false, 202, "you are not in group", nil
			}
		}
		isMute, isAdmin, err := userIsMuteAndIsAdminInGroup(ctx, data.MsgData.GroupID, data.MsgData.SendID)
		if err != nil {
			errMsg := data.OperationID + err.Error()
			return false, 223, errMsg, nil
		}
		if isMute {
			return false, 224, "you are muted", nil
		}
		if isAdmin {
			return true, 0, "", userIDList
		}
		isMute, err = groupIsMuted(ctx, data.MsgData.GroupID)
		if err != nil {
			errMsg := data.OperationID + err.Error()
			return false, 223, errMsg, nil
		}
		if isMute {
			return false, 225, "group id muted", nil
		}
		return true, 0, "", userIDList
	case constant.SuperGroupChatType:
		groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, data.MsgData.GroupID)
		if err != nil {
			return false, 201, err.Error(), nil
		}

		if data.MsgData.ContentType == constant.AdvancedRevoke {
			revokeMessage := new(MessageRevoked)
			err := utils.JsonStringToStruct(string(data.MsgData.Content), revokeMessage)
			if err != nil {
				log.Error(data.OperationID, "json unmarshal err:", err.Error())
				return false, 201, err.Error(), nil
			}
			log.Debug(data.OperationID, "revoke message is", *revokeMessage)
			if revokeMessage.RevokerID != revokeMessage.SourceMessageSendID {
				req := pbChat.GetSuperGroupMsgReq{OperationID: data.OperationID, Seq: revokeMessage.Seq, GroupID: data.MsgData.GroupID}
				resp, err := rpc.GetSuperGroupMsg(context.Background(), &req)
				if err != nil {
					log.Error(data.OperationID, "GetSuperGroupMsgReq err:", err.Error())
				} else if resp.ErrCode != 0 {
					log.Error(data.OperationID, "GetSuperGroupMsgReq err:", resp.ErrCode, resp.ErrMsg)
				} else {
					if resp.MsgData != nil && resp.MsgData.ClientMsgID == revokeMessage.ClientMsgID && resp.MsgData.Seq == revokeMessage.Seq {
						revokeMessage.SourceMessageSendTime = resp.MsgData.SendTime
						revokeMessage.SourceMessageSenderNickname = resp.MsgData.SenderNickname
						revokeMessage.SourceMessageSendID = resp.MsgData.SendID
						log.Debug(data.OperationID, "new revoke message is ", revokeMessage)
						data.MsgData.Content = []byte(utils.StructToJsonString(revokeMessage))
					} else {
						return false, 201, errors.New("msg err").Error(), nil
					}
				}
			}
		}
		if groupInfo.GroupType == constant.SuperGroup {
			return true, 0, "", nil
		} else {
			userIDList, err := utils.GetGroupMemberUserIDList(ctx, data.MsgData.GroupID, data.OperationID)
			if err != nil {
				errMsg := data.OperationID + err.Error()
				log.NewError(data.OperationID, errMsg)
				return false, 201, errMsg, nil
			}
			if tokenverify.IsManagerUserID(data.MsgData.SendID) {
				return true, 0, "", userIDList
			}
			if data.MsgData.ContentType <= constant.NotificationEnd && data.MsgData.ContentType >= constant.NotificationBegin {
				return true, 0, "", userIDList
			} else {
				if !utils.IsContain(data.MsgData.SendID, userIDList) {
					//return returnMsg(&replay, pb, 202, "you are not in group", "", 0)
					return false, 202, "you are not in group", nil
				}
			}
			isMute, isAdmin, err := userIsMuteAndIsAdminInGroup(ctx, data.MsgData.GroupID, data.MsgData.SendID)
			if err != nil {
				errMsg := data.OperationID + err.Error()
				return false, 223, errMsg, nil
			}
			if isMute {
				return false, 224, "you are muted", nil
			}
			if isAdmin {
				return true, 0, "", userIDList
			}
			isMute, err = groupIsMuted(ctx, data.MsgData.GroupID)
			if err != nil {
				errMsg := data.OperationID + err.Error()
				return false, 223, errMsg, nil
			}
			if isMute {
				return false, 225, "group id muted", nil
			}
			return true, 0, "", userIDList
		}
	default:
		return true, 0, "", nil
	}

}
func (rpc *msgServer) encapsulateMsgData(msg *sdkws.MsgData) {
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
		log.Info("", "this is a test start", msg, msg.Options)
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
		log.Info("", "this is a test end", msg, msg.Options)
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

func (rpc *msgServer) sendMsgToWriter(ctx context.Context, m *pbChat.MsgDataToMQ, key string, status string) error {
	switch status {
	case constant.OnlineStatus:
		if m.MsgData.ContentType == constant.SignalingNotification {
			rpcPushMsg := pbPush.PushMsgReq{OperationID: m.OperationID, MsgData: m.MsgData, PushToUserID: key}
			grpcConn, err := rpc.GetConn(ctx, config.Config.RpcRegisterName.OpenImPushName)
			if err != nil {
				return err
			}
			msgClient := pbPush.NewPushMsgServiceClient(grpcConn)
			_, err = msgClient.PushMsg(context.Background(), &rpcPushMsg)
			if err != nil {
				log.Error(rpcPushMsg.OperationID, "rpc send failed", rpcPushMsg.OperationID, "push data", rpcPushMsg.String(), "err", err.Error())
				return err
			} else {
				return nil
			}
		}
		pid, offset, err := rpc.messageWriter.SendMessage(m, key, m.OperationID)
		if err != nil {
			log.Error(m.OperationID, "kafka send failed", "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "key", key, status)
		} else {
			//	log.NewWarn(m.OperationID, "sendMsgToWriter   client msgID ", m.MsgData.ClientMsgID)
		}
		return err
	case constant.OfflineStatus:
		pid, offset, err := rpc.messageWriter.SendMessage(m, key, m.OperationID)
		if err != nil {
			log.Error(m.OperationID, "kafka send failed", "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "key", key, status)
		}
		return err
	}
	return errors.New("status error")
}
func GetMsgID(sendID string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return utils.Md5(t + "-" + sendID + "-" + strconv.Itoa(rand.Int()))
}

func returnMsg(replay *pbChat.SendMsgResp, pb *pbChat.SendMsgReq, errCode int32, errMsg, serverMsgID string, sendTime int64) (*pbChat.SendMsgResp, error) {
	replay.ErrCode = errCode
	replay.ErrMsg = errMsg
	replay.ServerMsgID = serverMsgID
	replay.ClientMsgID = pb.MsgData.ClientMsgID
	replay.SendTime = sendTime
	return replay, nil
}

func modifyMessageByUserMessageReceiveOpt(userID, sourceID string, sessionType int, pb *pbChat.SendMsgReq) bool {
	opt, err := db.DB.GetUserGlobalMsgRecvOpt(userID)
	if err != nil {
		log.NewError(pb.OperationID, "GetUserGlobalMsgRecvOpt from redis err", userID, pb.String(), err.Error())

	}
	switch opt {
	case constant.ReceiveMessage:
	case constant.NotReceiveMessage:
		return false
	case constant.ReceiveNotNotifyMessage:
		if pb.MsgData.Options == nil {
			pb.MsgData.Options = make(map[string]bool, 10)
		}
		utils.SetSwitchFromOptions(pb.MsgData.Options, constant.IsOfflinePush, false)
		return true
	}
	conversationID := utils.GetConversationIDBySessionType(sourceID, sessionType)
	singleOpt, sErr := db.DB.GetSingleConversationRecvMsgOpt(userID, conversationID)
	if sErr != nil && sErr != go_redis.Nil {
		log.NewError(pb.OperationID, "GetSingleConversationMsgOpt from redis err", conversationID, pb.String(), sErr.Error())
		return true
	}
	switch singleOpt {
	case constant.ReceiveMessage:
		return true
	case constant.NotReceiveMessage:
		if utils.IsContainInt(int(pb.MsgData.ContentType), ExcludeContentType) {
			return true
		}
		return false
	case constant.ReceiveNotNotifyMessage:
		if pb.MsgData.Options == nil {
			pb.MsgData.Options = make(map[string]bool, 10)
		}
		utils.SetSwitchFromOptions(pb.MsgData.Options, constant.IsOfflinePush, false)
		return true
	}

	return true
}

func modifyMessageByUserMessageReceiveOptoptimization(userID, sourceID string, sessionType int, operationID string, options *map[string]bool) bool {
	conversationID := utils.GetConversationIDBySessionType(sourceID, sessionType)
	opt, err := db.DB.GetSingleConversationRecvMsgOpt(userID, conversationID)
	if err != nil && err != go_redis.Nil {
		log.NewError(operationID, "GetSingleConversationMsgOpt from redis err", userID, conversationID, err.Error())
		return true
	}

	switch opt {
	case constant.ReceiveMessage:
		return true
	case constant.NotReceiveMessage:
		return false
	case constant.ReceiveNotNotifyMessage:
		if *options == nil {
			*options = make(map[string]bool, 10)
		}
		utils.SetSwitchFromOptions(*options, constant.IsOfflinePush, false)
		return true
	}
	return true
}

func getOnlineAndOfflineUserIDList(memberList []string, m map[string][]string, operationID string) {
	var onllUserIDList, offlUserIDList []string
	var wsResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	req := &pbRelay.GetUsersOnlineStatusReq{}
	req.UserIDList = memberList
	req.OperationID = operationID
	req.OpUserID = config.Config.Manager.AppManagerUid[0]
	flag := false
	grpcCons := rpc.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), operationID)
	for _, v := range grpcCons {
		client := pbRelay.NewRelayClient(v)
		reply, err := client.GetUsersOnlineStatus(context.Background(), req)
		if err != nil {
			log.NewError(operationID, "GetUsersOnlineStatus rpc  err", req.String(), err.Error())
			continue
		} else {
			if reply.ErrCode == 0 {
				wsResult = append(wsResult, reply.SuccessResult...)
			}
		}
	}
	log.NewInfo(operationID, "call GetUsersOnlineStatus rpc server is success", wsResult)
	//Online data merge of each node
	for _, v1 := range memberList {
		flag = false

		for _, v2 := range wsResult {
			if v2.UserID == v1 {
				flag = true
				onllUserIDList = append(onllUserIDList, v1)
			}

		}
		if !flag {
			offlUserIDList = append(offlUserIDList, v1)
		}
	}
	m[constant.OnlineStatus] = onllUserIDList
	m[constant.OfflineStatus] = offlUserIDList
}

func valueCopy(pb *pbChat.SendMsgReq) *pbChat.SendMsgReq {
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
	return &pbChat.SendMsgReq{Token: pb.Token, OperationID: pb.OperationID, MsgData: &msgData}
}

func (rpc *msgServer) sendMsgToGroup(ctx context.Context, list []string, pb pbChat.SendMsgReq, status string, sendTag *bool, wg *sync.WaitGroup) {
	//	log.Debug(pb.OperationID, "split userID ", list)
	offlinePushInfo := sdkws.OfflinePushInfo{}
	if pb.MsgData.OfflinePushInfo != nil {
		offlinePushInfo = *pb.MsgData.OfflinePushInfo
	}
	msgData := sdkws.MsgData{}
	msgData = *pb.MsgData
	msgData.OfflinePushInfo = &offlinePushInfo

	groupPB := pbChat.SendMsgReq{Token: pb.Token, OperationID: pb.OperationID, MsgData: &msgData}
	msgToMQGroup := pbChat.MsgDataToMQ{Token: pb.Token, OperationID: pb.OperationID, MsgData: &msgData}
	for _, v := range list {
		options := make(map[string]bool, 10)
		for key, value := range pb.MsgData.Options {
			options[key] = value
		}
		groupPB.MsgData.RecvID = v
		groupPB.MsgData.Options = options
		isSend := modifyMessageByUserMessageReceiveOpt(v, msgData.GroupID, constant.GroupChatType, &groupPB)
		if isSend {
			msgToMQGroup.MsgData = groupPB.MsgData
			//	log.Debug(groupPB.OperationID, "sendMsgToWriter, ", v, groupID, msgToMQGroup.String())
			err := rpc.sendMsgToWriter(ctx, &msgToMQGroup, v, status)
			if err != nil {
				log.NewError(msgToMQGroup.OperationID, "kafka send msg err:UserId", v, msgToMQGroup.String())
			} else {
				*sendTag = true
			}
		} else {
			log.Debug(groupPB.OperationID, "not sendMsgToWriter, ", v)
		}
	}
	wg.Done()
}

func (rpc *msgServer) sendMsgToGroupOptimization(ctx context.Context, list []string, groupPB *pbChat.SendMsgReq, status string, sendTag *bool, wg *sync.WaitGroup) {
	msgToMQGroup := pbChat.MsgDataToMQ{Token: groupPB.Token, OperationID: groupPB.OperationID, MsgData: groupPB.MsgData}
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
		isSend := modifyMessageByUserMessageReceiveOpt(v, groupPB.MsgData.GroupID, constant.GroupChatType, groupPB)
		if isSend {
			if v == "" || groupPB.MsgData.SendID == "" {
				log.Error(msgToMQGroup.OperationID, "sendMsgToGroupOptimization userID nil ", msgToMQGroup.String())
				continue
			}
			err := rpc.sendMsgToWriter(ctx, &msgToMQGroup, v, status)
			if err != nil {
				log.NewError(msgToMQGroup.OperationID, "kafka send msg err:UserId", v, msgToMQGroup.String())
			} else {
				*sendTag = true
			}
		} else {
			log.Debug(groupPB.OperationID, "not sendMsgToWriter, ", v)
		}
	}
	wg.Done()
}

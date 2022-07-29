package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	cacheRpc "Open_IM/pkg/proto/cache"
	pbCache "Open_IM/pkg/proto/cache"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbChat "Open_IM/pkg/proto/msg"
	pbRelay "Open_IM/pkg/proto/relay"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

//When the number of group members is greater than this value，Online users will be sent first，Guaranteed service availability
const GroupMemberNum = 500

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

func isMessageHasReadEnabled(pb *pbChat.SendMsgReq) (bool, int32, string) {
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

func messageVerification(data *pbChat.SendMsgReq) (bool, int32, string, []string) {
	switch data.MsgData.SessionType {
	case constant.SingleChatType:
		if utils.IsContain(data.MsgData.SendID, config.Config.Manager.AppManagerUid) {
			return true, 0, "", nil
		}
		if data.MsgData.ContentType <= constant.NotificationEnd && data.MsgData.ContentType >= constant.NotificationBegin {
			return true, 0, "", nil
		}
		log.NewDebug(data.OperationID, config.Config.MessageVerify.FriendVerify)
		reqGetBlackIDListFromCache := &cacheRpc.GetBlackIDListFromCacheReq{UserID: data.MsgData.RecvID, OperationID: data.OperationID}
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, data.OperationID)
		if etcdConn == nil {
			errMsg := data.OperationID + "getcdv3.GetConn == nil"
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
		log.NewDebug(data.OperationID, config.Config.MessageVerify.FriendVerify)
		if config.Config.MessageVerify.FriendVerify {
			reqGetFriendIDListFromCache := &cacheRpc.GetFriendIDListFromCacheReq{UserID: data.MsgData.RecvID, OperationID: data.OperationID}
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, data.OperationID)
			if etcdConn == nil {
				errMsg := data.OperationID + "getcdv3.GetConn == nil"
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
		fallthrough
	case constant.SuperGroupChatType:
		getGroupMemberIDListFromCacheReq := &pbCache.GetGroupMemberIDListFromCacheReq{OperationID: data.OperationID, GroupID: data.MsgData.GroupID}
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, data.OperationID)
		if etcdConn == nil {
			errMsg := data.OperationID + "getcdv3.GetConn == nil"
			log.NewError(data.OperationID, errMsg)
			//return returnMsg(&replay, pb, 201, errMsg, "", 0)
			return false, 201, errMsg, nil
		}
		client := pbCache.NewCacheClient(etcdConn)
		cacheResp, err := client.GetGroupMemberIDListFromCache(context.Background(), getGroupMemberIDListFromCacheReq)
		if err != nil {
			log.NewError(data.OperationID, "GetGroupMemberIDListFromCache rpc call failed ", err.Error())
			//return returnMsg(&replay, pb, 201, "GetGroupMemberIDListFromCache failed", "", 0)
			return false, 201, err.Error(), nil
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			log.NewError(data.OperationID, "GetGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
			//return returnMsg(&replay, pb, 201, "GetGroupMemberIDListFromCache logic failed", "", 0)
			return false, cacheResp.CommonResp.ErrCode, cacheResp.CommonResp.ErrMsg, nil
		}
		if !token_verify.IsManagerUserID(data.MsgData.SendID) {
			if !utils.IsContain(data.MsgData.SendID, cacheResp.UserIDList) {
				//return returnMsg(&replay, pb, 202, "you are not in group", "", 0)
				return false, 202, "you are not in group", nil
			}
		}
		return true, 0, "", cacheResp.UserIDList

	default:
		return true, 0, "", nil
	}

}
func (rpc *rpcChat) encapsulateMsgData(msg *sdk_ws.MsgData) {
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
func (rpc *rpcChat) SendMsg(_ context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, error) {
	replay := pbChat.SendMsgResp{}
	newTime := db.GetCurrentTimestampByMill()
	log.Info(pb.OperationID, "rpc sendMsg come here ", pb.String())
	flag, errCode, errMsg := isMessageHasReadEnabled(pb)
	log.Info(pb.OperationID, "isMessageHasReadEnabled ", flag)
	if !flag {
		return returnMsg(&replay, pb, errCode, errMsg, "", 0)
	}
	flag, errCode, errMsg, _ = messageVerification(pb)
	log.Info(pb.OperationID, "userRelationshipVerification ", flag)
	if !flag {
		return returnMsg(&replay, pb, errCode, errMsg, "", 0)
	}
	rpc.encapsulateMsgData(pb.MsgData)
	msgToMQSingle := pbChat.MsgDataToMQ{Token: pb.Token, OperationID: pb.OperationID, MsgData: pb.MsgData}
	// callback
	callbackResp := callbackWordFilter(pb)
	log.Info(pb.OperationID, "callbackWordFilter ", callbackResp)
	if callbackResp.ErrCode != 0 {
		log.Error(pb.OperationID, utils.GetSelfFuncName(), "callbackWordFilter resp: ", callbackResp)
	}
	log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackResp: ", callbackResp)
	if callbackResp.ActionCode != constant.ActionAllow {
		if callbackResp.ErrCode == 0 {
			callbackResp.ErrCode = 201
		}
		log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackWordFilter result", "end rpc and return", pb.MsgData)
		return returnMsg(&replay, pb, int32(callbackResp.ErrCode), callbackResp.ErrMsg, "", 0)
	}
	switch pb.MsgData.SessionType {
	case constant.SingleChatType:
		// callback
		callbackResp := callbackBeforeSendSingleMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg resp: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			if callbackResp.ErrCode == 0 {
				callbackResp.ErrCode = 201
			}
			log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
			return returnMsg(&replay, pb, int32(callbackResp.ErrCode), callbackResp.ErrMsg, "", 0)
		}
		isSend := modifyMessageByUserMessageReceiveOpt(pb.MsgData.RecvID, pb.MsgData.SendID, constant.SingleChatType, pb)
		if isSend {
			msgToMQSingle.MsgData = pb.MsgData
			log.NewInfo(msgToMQSingle.OperationID, msgToMQSingle)
			err1 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.RecvID, constant.OnlineStatus)
			if err1 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err :RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String(), err1.Error())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			err2 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.SendID, constant.OnlineStatus)
			if err2 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err:SendID", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		// callback
		callbackResp = callbackAfterSendSingleMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackAfterSendSingleMsg resp: ", callbackResp)
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)
	case constant.GroupChatType:
		// callback
		callbackResp := callbackBeforeSendGroupMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendGroupMsg resp:", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			if callbackResp.ErrCode == 0 {
				callbackResp.ErrCode = 201
			}
			log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
			return returnMsg(&replay, pb, int32(callbackResp.ErrCode), callbackResp.ErrMsg, "", 0)
		}
		var memberUserIDList []string
		if flag, errCode, errMsg, memberUserIDList = messageVerification(pb); !flag {
			return returnMsg(&replay, pb, errCode, errMsg, "", 0)
		}
		log.Debug(pb.OperationID, "GetGroupAllMember userID list", memberUserIDList, "len: ", len(memberUserIDList))
		var addUidList []string
		switch pb.MsgData.ContentType {
		case constant.MemberKickedNotification:
			var tips sdk_ws.TipsComm
			var memberKickedTips sdk_ws.MemberKickedTips
			err := proto.Unmarshal(pb.MsgData.Content, &tips)
			if err != nil {
				log.Error(pb.OperationID, "Unmarshal err", err.Error())
			}
			err = proto.Unmarshal(tips.Detail, &memberKickedTips)
			if err != nil {
				log.Error(pb.OperationID, "Unmarshal err", err.Error())
			}
			log.Info(pb.OperationID, "data is ", memberKickedTips)
			for _, v := range memberKickedTips.KickedUserList {
				addUidList = append(addUidList, v.UserID)
			}
		case constant.MemberQuitNotification:
			addUidList = append(addUidList, pb.MsgData.SendID)

		default:
		}
		if len(addUidList) > 0 {
			memberUserIDList = append(memberUserIDList, addUidList...)
		}
		m := make(map[string][]string, 2)
		m[constant.OnlineStatus] = memberUserIDList
		log.Debug(pb.OperationID, "send msg cost time1 ", db.GetCurrentTimestampByMill()-newTime, pb.MsgData.ClientMsgID)
		newTime = db.GetCurrentTimestampByMill()

		//split  parallel send
		var wg sync.WaitGroup
		var sendTag bool
		var split = 20
		for k, v := range m {
			remain := len(v) % split
			for i := 0; i < len(v)/split; i++ {
				wg.Add(1)
				tmp := valueCopy(pb)
				//	go rpc.sendMsgToGroup(v[i*split:(i+1)*split], *pb, k, &sendTag, &wg)
				go rpc.sendMsgToGroupOptimization(v[i*split:(i+1)*split], tmp, k, &sendTag, &wg)
			}
			if remain > 0 {
				wg.Add(1)
				tmp := valueCopy(pb)
				//	go rpc.sendMsgToGroup(v[split*(len(v)/split):], *pb, k, &sendTag, &wg)
				go rpc.sendMsgToGroupOptimization(v[split*(len(v)/split):], tmp, k, &sendTag, &wg)
			}
		}
		log.Debug(pb.OperationID, "send msg cost time22 ", db.GetCurrentTimestampByMill()-newTime, pb.MsgData.ClientMsgID, "uidList : ", len(addUidList))
		//wg.Add(1)
		//go rpc.sendMsgToGroup(addUidList, *pb, constant.OnlineStatus, &sendTag, &wg)
		wg.Wait()
		newTime = db.GetCurrentTimestampByMill()
		// callback
		callbackResp = callbackAfterSendGroupMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackAfterSendGroupMsg resp: ", callbackResp)
		}
		if !sendTag {
			log.NewWarn(pb.OperationID, "send tag is ", sendTag)
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		} else {
			if pb.MsgData.ContentType == constant.AtText {
				go func() {
					var conversationReq pbConversation.ModifyConversationFieldReq
					var tag bool
					var atUserID []string
					conversation := pbConversation.Conversation{
						OwnerUserID:      pb.MsgData.SendID,
						ConversationID:   utils.GetConversationIDBySessionType(pb.MsgData.GroupID, constant.GroupChatType),
						ConversationType: constant.GroupChatType,
						GroupID:          pb.MsgData.GroupID,
					}
					conversationReq.Conversation = &conversation
					conversationReq.OperationID = pb.OperationID
					conversationReq.FieldType = constant.FieldGroupAtType
					tagAll := utils.IsContain(constant.AtAllString, pb.MsgData.AtUserIDList)
					if tagAll {
						atUserID = utils.DifferenceString([]string{constant.AtAllString}, pb.MsgData.AtUserIDList)
						if len(atUserID) == 0 { //just @everyone
							conversationReq.UserIDList = memberUserIDList
							conversation.GroupAtType = constant.AtAll
						} else { //@Everyone and @other people
							conversationReq.UserIDList = atUserID
							conversation.GroupAtType = constant.AtAllAtMe
							tag = true
						}
					} else {
						conversationReq.UserIDList = pb.MsgData.AtUserIDList
						conversation.GroupAtType = constant.AtMe
					}
					etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImConversationName, pb.OperationID)
					if etcdConn == nil {
						errMsg := pb.OperationID + "getcdv3.GetConn == nil"
						log.NewError(pb.OperationID, errMsg)
						return
					}
					client := pbConversation.NewConversationClient(etcdConn)
					conversationReply, err := client.ModifyConversationField(context.Background(), &conversationReq)
					if err != nil {
						log.NewError(conversationReq.OperationID, "ModifyConversationField rpc failed, ", conversationReq.String(), err.Error())
					} else if conversationReply.CommonResp.ErrCode != 0 {
						log.NewError(conversationReq.OperationID, "ModifyConversationField rpc failed, ", conversationReq.String(), conversationReply.String())
					}
					if tag {
						conversationReq.UserIDList = utils.DifferenceString(atUserID, memberUserIDList)
						conversation.GroupAtType = constant.AtAll
						etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImConversationName, pb.OperationID)
						if etcdConn == nil {
							errMsg := pb.OperationID + "getcdv3.GetConn == nil"
							log.NewError(pb.OperationID, errMsg)
							return
						}
						client := pbConversation.NewConversationClient(etcdConn)
						conversationReply, err := client.ModifyConversationField(context.Background(), &conversationReq)
						if err != nil {
							log.NewError(conversationReq.OperationID, "ModifyConversationField rpc failed, ", conversationReq.String(), err.Error())
						} else if conversationReply.CommonResp.ErrCode != 0 {
							log.NewError(conversationReq.OperationID, "ModifyConversationField rpc failed, ", conversationReq.String(), conversationReply.String())
						}
					}
				}()
			}
			log.Debug(pb.OperationID, "send msg cost time3 ", db.GetCurrentTimestampByMill()-newTime, pb.MsgData.ClientMsgID)
			return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)
		}
	case constant.NotificationChatType:
		msgToMQSingle.MsgData = pb.MsgData
		log.NewInfo(msgToMQSingle.OperationID, msgToMQSingle)
		err1 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.RecvID, constant.OnlineStatus)
		if err1 != nil {
			log.NewError(msgToMQSingle.OperationID, "kafka send msg err:RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		}

		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			err2 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.SendID, constant.OnlineStatus)
			if err2 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err:SendID", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}

		log.Debug(pb.OperationID, "send msg cost time ", db.GetCurrentTimestampByMill()-newTime, pb.MsgData.ClientMsgID)
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)
	case constant.SuperGroupChatType:
		// callback
		callbackResp := callbackBeforeSendGroupMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSuperGroupMsg resp: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			if callbackResp.ErrCode == 0 {
				callbackResp.ErrCode = 201
			}
			log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSuperGroupMsg result", "end rpc and return", callbackResp)
			return returnMsg(&replay, pb, int32(callbackResp.ErrCode), callbackResp.ErrMsg, "", 0)
		}
		if flag, errCode, errMsg, _ = messageVerification(pb); !flag {
			return returnMsg(&replay, pb, errCode, errMsg, "", 0)
		}
		msgToMQSingle.MsgData = pb.MsgData
		log.NewInfo(msgToMQSingle.OperationID, msgToMQSingle)
		err1 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.GroupID, constant.OnlineStatus)
		if err1 != nil {
			log.NewError(msgToMQSingle.OperationID, "kafka send msg err:RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		}
		// callback
		callbackResp = callbackAfterSendSingleMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackAfterSendSuperGroupMsg resp: ", callbackResp)
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)

	default:
		return returnMsg(&replay, pb, 203, "unknown sessionType", "", 0)
	}
}

func (rpc *rpcChat) sendMsgToKafka(m *pbChat.MsgDataToMQ, key string, status string) error {
	switch status {
	case constant.OnlineStatus:
		pid, offset, err := rpc.onlineProducer.SendMessage(m, key, m.OperationID)
		if err != nil {
			log.Error(m.OperationID, "kafka send failed", "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "key", key, status)
		} else {
			//	log.NewWarn(m.OperationID, "sendMsgToKafka   client msgID ", m.MsgData.ClientMsgID)
		}
		return err
	case constant.OfflineStatus:
		pid, offset, err := rpc.onlineProducer.SendMessage(m, key, m.OperationID)
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

type NotificationMsg struct {
	SendID         string
	RecvID         string
	Content        []byte //  open_im_sdk.TipsComm
	MsgFrom        int32
	ContentType    int32
	SessionType    int32
	OperationID    string
	SenderNickname string
	SenderFaceURL  string
}

func Notification(n *NotificationMsg) {
	var req pbChat.SendMsgReq
	var msg sdk_ws.MsgData
	var offlineInfo sdk_ws.OfflinePushInfo
	var title, desc, ex string
	var pushSwitch, unReadCount bool
	var reliabilityLevel int
	req.OperationID = n.OperationID
	msg.SendID = n.SendID
	msg.RecvID = n.RecvID
	msg.Content = n.Content
	msg.MsgFrom = n.MsgFrom
	msg.ContentType = n.ContentType
	msg.SessionType = n.SessionType
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(n.SendID)
	msg.Options = make(map[string]bool, 7)
	msg.SenderNickname = n.SenderNickname
	msg.SenderFaceURL = n.SenderFaceURL
	switch n.SessionType {
	case constant.GroupChatType, constant.SuperGroupChatType:
		msg.RecvID = ""
		msg.GroupID = n.RecvID
	}
	offlineInfo.IOSBadgeCount = config.Config.IOSPush.BadgeCount
	offlineInfo.IOSPushSound = config.Config.IOSPush.PushSound
	switch msg.ContentType {
	case constant.GroupCreatedNotification:
		pushSwitch = config.Config.Notification.GroupCreated.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupCreated.OfflinePush.Title
		desc = config.Config.Notification.GroupCreated.OfflinePush.Desc
		ex = config.Config.Notification.GroupCreated.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupCreated.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupCreated.Conversation.UnreadCount
	case constant.GroupInfoSetNotification:
		pushSwitch = config.Config.Notification.GroupInfoSet.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupInfoSet.OfflinePush.Title
		desc = config.Config.Notification.GroupInfoSet.OfflinePush.Desc
		ex = config.Config.Notification.GroupInfoSet.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupInfoSet.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupInfoSet.Conversation.UnreadCount
	case constant.JoinGroupApplicationNotification:
		pushSwitch = config.Config.Notification.JoinGroupApplication.OfflinePush.PushSwitch
		title = config.Config.Notification.JoinGroupApplication.OfflinePush.Title
		desc = config.Config.Notification.JoinGroupApplication.OfflinePush.Desc
		ex = config.Config.Notification.JoinGroupApplication.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.JoinGroupApplication.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.JoinGroupApplication.Conversation.UnreadCount
	case constant.MemberQuitNotification:
		pushSwitch = config.Config.Notification.MemberQuit.OfflinePush.PushSwitch
		title = config.Config.Notification.MemberQuit.OfflinePush.Title
		desc = config.Config.Notification.MemberQuit.OfflinePush.Desc
		ex = config.Config.Notification.MemberQuit.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.MemberQuit.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.MemberQuit.Conversation.UnreadCount
	case constant.GroupApplicationAcceptedNotification:
		pushSwitch = config.Config.Notification.GroupApplicationAccepted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Title
		desc = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Desc
		ex = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupApplicationAccepted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupApplicationAccepted.Conversation.UnreadCount
	case constant.GroupApplicationRejectedNotification:
		pushSwitch = config.Config.Notification.GroupApplicationRejected.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupApplicationRejected.OfflinePush.Title
		desc = config.Config.Notification.GroupApplicationRejected.OfflinePush.Desc
		ex = config.Config.Notification.GroupApplicationRejected.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupApplicationRejected.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupApplicationRejected.Conversation.UnreadCount
	case constant.GroupOwnerTransferredNotification:
		pushSwitch = config.Config.Notification.GroupOwnerTransferred.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Title
		desc = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Desc
		ex = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupOwnerTransferred.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupOwnerTransferred.Conversation.UnreadCount
	case constant.MemberKickedNotification:
		pushSwitch = config.Config.Notification.MemberKicked.OfflinePush.PushSwitch
		title = config.Config.Notification.MemberKicked.OfflinePush.Title
		desc = config.Config.Notification.MemberKicked.OfflinePush.Desc
		ex = config.Config.Notification.MemberKicked.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.MemberKicked.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.MemberKicked.Conversation.UnreadCount
	case constant.MemberInvitedNotification:
		pushSwitch = config.Config.Notification.MemberInvited.OfflinePush.PushSwitch
		title = config.Config.Notification.MemberInvited.OfflinePush.Title
		desc = config.Config.Notification.MemberInvited.OfflinePush.Desc
		ex = config.Config.Notification.MemberInvited.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.MemberInvited.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.MemberInvited.Conversation.UnreadCount
	case constant.MemberEnterNotification:
		pushSwitch = config.Config.Notification.MemberEnter.OfflinePush.PushSwitch
		title = config.Config.Notification.MemberEnter.OfflinePush.Title
		desc = config.Config.Notification.MemberEnter.OfflinePush.Desc
		ex = config.Config.Notification.MemberEnter.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.MemberEnter.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.MemberEnter.Conversation.UnreadCount
	case constant.UserInfoUpdatedNotification:
		pushSwitch = config.Config.Notification.UserInfoUpdated.OfflinePush.PushSwitch
		title = config.Config.Notification.UserInfoUpdated.OfflinePush.Title
		desc = config.Config.Notification.UserInfoUpdated.OfflinePush.Desc
		ex = config.Config.Notification.UserInfoUpdated.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.UserInfoUpdated.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.UserInfoUpdated.Conversation.UnreadCount
	case constant.FriendApplicationNotification:
		pushSwitch = config.Config.Notification.FriendApplication.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendApplication.OfflinePush.Title
		desc = config.Config.Notification.FriendApplication.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplication.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendApplication.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendApplication.Conversation.UnreadCount
	case constant.FriendApplicationApprovedNotification:
		pushSwitch = config.Config.Notification.FriendApplicationApproved.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendApplicationApproved.OfflinePush.Title
		desc = config.Config.Notification.FriendApplicationApproved.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplicationApproved.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendApplicationApproved.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendApplicationApproved.Conversation.UnreadCount
	case constant.FriendApplicationRejectedNotification:
		pushSwitch = config.Config.Notification.FriendApplicationRejected.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendApplicationRejected.OfflinePush.Title
		desc = config.Config.Notification.FriendApplicationRejected.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplicationRejected.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendApplicationRejected.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendApplicationRejected.Conversation.UnreadCount
	case constant.FriendAddedNotification:
		pushSwitch = config.Config.Notification.FriendAdded.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendAdded.OfflinePush.Title
		desc = config.Config.Notification.FriendAdded.OfflinePush.Desc
		ex = config.Config.Notification.FriendAdded.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendAdded.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendAdded.Conversation.UnreadCount
	case constant.FriendDeletedNotification:
		pushSwitch = config.Config.Notification.FriendDeleted.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendDeleted.OfflinePush.Title
		desc = config.Config.Notification.FriendDeleted.OfflinePush.Desc
		ex = config.Config.Notification.FriendDeleted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendDeleted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendDeleted.Conversation.UnreadCount
	case constant.FriendRemarkSetNotification:
		pushSwitch = config.Config.Notification.FriendRemarkSet.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendRemarkSet.OfflinePush.Title
		desc = config.Config.Notification.FriendRemarkSet.OfflinePush.Desc
		ex = config.Config.Notification.FriendRemarkSet.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendRemarkSet.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendRemarkSet.Conversation.UnreadCount
	case constant.BlackAddedNotification:
		pushSwitch = config.Config.Notification.BlackAdded.OfflinePush.PushSwitch
		title = config.Config.Notification.BlackAdded.OfflinePush.Title
		desc = config.Config.Notification.BlackAdded.OfflinePush.Desc
		ex = config.Config.Notification.BlackAdded.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.BlackAdded.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.BlackAdded.Conversation.UnreadCount
	case constant.BlackDeletedNotification:
		pushSwitch = config.Config.Notification.BlackDeleted.OfflinePush.PushSwitch
		title = config.Config.Notification.BlackDeleted.OfflinePush.Title
		desc = config.Config.Notification.BlackDeleted.OfflinePush.Desc
		ex = config.Config.Notification.BlackDeleted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.BlackDeleted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.BlackDeleted.Conversation.UnreadCount
	case constant.ConversationOptChangeNotification:
		pushSwitch = config.Config.Notification.ConversationOptUpdate.OfflinePush.PushSwitch
		title = config.Config.Notification.ConversationOptUpdate.OfflinePush.Title
		desc = config.Config.Notification.ConversationOptUpdate.OfflinePush.Desc
		ex = config.Config.Notification.ConversationOptUpdate.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.ConversationOptUpdate.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.ConversationOptUpdate.Conversation.UnreadCount

	case constant.GroupDismissedNotification:
		pushSwitch = config.Config.Notification.GroupDismissed.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupDismissed.OfflinePush.Title
		desc = config.Config.Notification.GroupDismissed.OfflinePush.Desc
		ex = config.Config.Notification.GroupDismissed.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupDismissed.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupDismissed.Conversation.UnreadCount

	case constant.GroupMutedNotification:
		pushSwitch = config.Config.Notification.GroupMuted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMuted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupMuted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupMuted.Conversation.UnreadCount

	case constant.GroupCancelMutedNotification:
		pushSwitch = config.Config.Notification.GroupCancelMuted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupCancelMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupCancelMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupCancelMuted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupCancelMuted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupCancelMuted.Conversation.UnreadCount

	case constant.GroupMemberMutedNotification:
		pushSwitch = config.Config.Notification.GroupMemberMuted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupMemberMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberMuted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupMemberMuted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupMemberMuted.Conversation.UnreadCount

	case constant.GroupMemberCancelMutedNotification:
		pushSwitch = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupMemberCancelMuted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupMemberCancelMuted.Conversation.UnreadCount

	case constant.GroupMemberInfoSetNotification:
		pushSwitch = config.Config.Notification.GroupMemberInfoSet.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupMemberInfoSet.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupMemberInfoSet.Conversation.UnreadCount

	case constant.OrganizationChangedNotification:
		pushSwitch = config.Config.Notification.OrganizationChanged.OfflinePush.PushSwitch
		title = config.Config.Notification.OrganizationChanged.OfflinePush.Title
		desc = config.Config.Notification.OrganizationChanged.OfflinePush.Desc
		ex = config.Config.Notification.OrganizationChanged.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.OrganizationChanged.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.OrganizationChanged.Conversation.UnreadCount

	case constant.WorkMomentNotification:
		pushSwitch = config.Config.Notification.WorkMomentsNotification.OfflinePush.PushSwitch
		title = config.Config.Notification.WorkMomentsNotification.OfflinePush.Title
		desc = config.Config.Notification.WorkMomentsNotification.OfflinePush.Desc
		ex = config.Config.Notification.WorkMomentsNotification.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.WorkMomentsNotification.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.WorkMomentsNotification.Conversation.UnreadCount

	case constant.ConversationPrivateChatNotification:
		pushSwitch = config.Config.Notification.ConversationSetPrivate.OfflinePush.PushSwitch
		title = config.Config.Notification.ConversationSetPrivate.OfflinePush.Title
		desc = config.Config.Notification.ConversationSetPrivate.OfflinePush.Desc
		ex = config.Config.Notification.ConversationSetPrivate.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.ConversationSetPrivate.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.ConversationSetPrivate.Conversation.UnreadCount
	case constant.DeleteMessageNotification:
		reliabilityLevel = constant.ReliableNotificationNoMsg
	case constant.SuperGroupUpdateNotification:
		reliabilityLevel = constant.UnreliableNotification
	}
	switch reliabilityLevel {
	case constant.UnreliableNotification:
		utils.SetSwitchFromOptions(msg.Options, constant.IsHistory, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsPersistent, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
	case constant.ReliableNotificationNoMsg:
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
	case constant.ReliableNotificationMsg:

	}
	utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, unReadCount)
	utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, pushSwitch)
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ex = ex
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		return
	}

	client := pbChat.NewMsgClient(etcdConn)
	reply, err := client.SendMsg(context.Background(), &req)
	if err != nil {
		log.NewError(req.OperationID, "SendMsg rpc failed, ", req.String(), err.Error())
	} else if reply.ErrCode != 0 {
		log.NewError(req.OperationID, "SendMsg rpc failed, ", req.String(), reply.ErrCode, reply.ErrMsg)
	}
}
func getOnlineAndOfflineUserIDList(memberList []string, m map[string][]string, operationID string) {
	var onllUserIDList, offlUserIDList []string
	var wsResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	req := &pbRelay.GetUsersOnlineStatusReq{}
	req.UserIDList = memberList
	req.OperationID = operationID
	req.OpUserID = config.Config.Manager.AppManagerUid[0]
	flag := false
	grpcCons := getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImRelayName)
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
	offlinePushInfo := sdk_ws.OfflinePushInfo{}
	if pb.MsgData.OfflinePushInfo != nil {
		offlinePushInfo = *pb.MsgData.OfflinePushInfo
	}
	msgData := sdk_ws.MsgData{}
	msgData = *pb.MsgData
	msgData.OfflinePushInfo = &offlinePushInfo

	options := make(map[string]bool, 10)
	for key, value := range pb.MsgData.Options {
		options[key] = value
	}
	msgData.Options = options
	return &pbChat.SendMsgReq{Token: pb.Token, OperationID: pb.OperationID, MsgData: &msgData}
}

func (rpc *rpcChat) sendMsgToGroup(list []string, pb pbChat.SendMsgReq, status string, sendTag *bool, wg *sync.WaitGroup) {
	//	log.Debug(pb.OperationID, "split userID ", list)
	offlinePushInfo := sdk_ws.OfflinePushInfo{}
	if pb.MsgData.OfflinePushInfo != nil {
		offlinePushInfo = *pb.MsgData.OfflinePushInfo
	}
	msgData := sdk_ws.MsgData{}
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
			//	log.Debug(groupPB.OperationID, "sendMsgToKafka, ", v, groupID, msgToMQGroup.String())
			err := rpc.sendMsgToKafka(&msgToMQGroup, v, status)
			if err != nil {
				log.NewError(msgToMQGroup.OperationID, "kafka send msg err:UserId", v, msgToMQGroup.String())
			} else {
				*sendTag = true
			}
		} else {
			log.Debug(groupPB.OperationID, "not sendMsgToKafka, ", v)
		}
	}
	wg.Done()
}

func (rpc *rpcChat) sendMsgToGroupOptimization(list []string, groupPB *pbChat.SendMsgReq, status string, sendTag *bool, wg *sync.WaitGroup) {
	msgToMQGroup := pbChat.MsgDataToMQ{Token: groupPB.Token, OperationID: groupPB.OperationID, MsgData: groupPB.MsgData}
	for _, v := range list {
		groupPB.MsgData.RecvID = v
		isSend := modifyMessageByUserMessageReceiveOpt(v, groupPB.MsgData.GroupID, constant.GroupChatType, groupPB)
		if isSend {
			if v == "" || groupPB.MsgData.SendID == "" {
				log.Error(msgToMQGroup.OperationID, "sendMsgToGroupOptimization userID nil ", msgToMQGroup.String())
				continue
			}
			err := rpc.sendMsgToKafka(&msgToMQGroup, v, status)
			if err != nil {
				log.NewError(msgToMQGroup.OperationID, "kafka send msg err:UserId", v, msgToMQGroup.String())
			} else {
				*sendTag = true
			}
		} else {
			log.Debug(groupPB.OperationID, "not sendMsgToKafka, ", v)
		}
	}
	wg.Done()
}

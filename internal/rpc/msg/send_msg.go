package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	cacheRpc "Open_IM/pkg/proto/cache"
	pbChat "Open_IM/pkg/proto/chat"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbGroup "Open_IM/pkg/proto/group"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

func userRelationshipVerification(data *pbChat.SendMsgReq) (bool, int32, string) {
	if data.MsgData.SessionType == constant.GroupChatType {
		return true, 0, ""
	}
	log.NewDebug(data.OperationID, config.Config.MessageVerify.FriendVerify)
	reqGetBlackIDListFromCache := &cacheRpc.GetBlackIDListFromCacheReq{UserID: data.MsgData.RecvID, OperationID: data.OperationID}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName)
	cacheClient := cacheRpc.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.GetBlackIDListFromCache(context.Background(), reqGetBlackIDListFromCache)
	if err != nil {
		log.NewError(data.OperationID, "GetBlackIDListFromCache rpc call failed ", err.Error())
	} else {
		if cacheResp.CommonResp.ErrCode != 0 {
			log.NewError(data.OperationID, "GetBlackIDListFromCache rpc logic call failed ", cacheResp.String())
		} else {
			if utils.IsContain(data.MsgData.SendID, cacheResp.UserIDList) {
				return false, 600, "in black list"
			}
		}
	}
	log.NewDebug(data.OperationID, config.Config.MessageVerify.FriendVerify)
	if config.Config.MessageVerify.FriendVerify {
		reqGetFriendIDListFromCache := &cacheRpc.GetFriendIDListFromCacheReq{UserID: data.MsgData.RecvID, OperationID: data.OperationID}
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName)
		cacheClient := cacheRpc.NewCacheClient(etcdConn)
		cacheResp, err := cacheClient.GetFriendIDListFromCache(context.Background(), reqGetFriendIDListFromCache)
		if err != nil {
			log.NewError(data.OperationID, "GetFriendIDListFromCache rpc call failed ", err.Error())
		} else {
			if cacheResp.CommonResp.ErrCode != 0 {
				log.NewError(data.OperationID, "GetFriendIDListFromCache rpc logic call failed ", cacheResp.String())
			} else {
				if !utils.IsContain(data.MsgData.SendID, cacheResp.UserIDList) {
					return false, 601, "not friend"
				}
			}
		}
		return true, 0, ""
	} else {
		return true, 0, ""
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
	log.NewDebug(pb.OperationID, "rpc sendMsg come here", pb.String())
	flag, errCode, errMsg := userRelationshipVerification(pb)
	if !flag {
		return returnMsg(&replay, pb, errCode, errMsg, "", 0)
	}
	//if !utils.VerifyToken(pb.Token, pb.SendID) {
	//	return returnMsg(&replay, pb, http.StatusUnauthorized, "token validate err,not authorized", "", 0)
	rpc.encapsulateMsgData(pb.MsgData)
	log.Info("", "this is a test MsgData ", pb.MsgData)
	msgToMQSingle := pbChat.MsgDataToMQ{Token: pb.Token, OperationID: pb.OperationID, MsgData: pb.MsgData}
	//options := utils.JsonStringToMap(pbData.Options)
	isHistory := utils.GetSwitchFromOptions(pb.MsgData.Options, constant.IsHistory)
	mReq := MsgCallBackReq{
		SendID:      pb.MsgData.SendID,
		RecvID:      pb.MsgData.RecvID,
		Content:     string(pb.MsgData.Content),
		SendTime:    pb.MsgData.SendTime,
		MsgFrom:     pb.MsgData.MsgFrom,
		ContentType: pb.MsgData.ContentType,
		SessionType: pb.MsgData.SessionType,
		PlatformID:  pb.MsgData.SenderPlatformID,
		MsgID:       pb.MsgData.ClientMsgID,
	}
	if !isHistory {
		mReq.IsOnlineOnly = true
	}

	// callback
	canSend, err := callbackWordFilter(pb)
	if err != nil {
		log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackWordFilter failed", err.Error(), pb.MsgData)
	}
	if !canSend {
		log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackWordFilter result", canSend, "end rpc and return", pb.MsgData)
		return returnMsg(&replay, pb, 201, "callbackWordFilter result stop rpc and return", "", 0)
	}
	switch pb.MsgData.SessionType {
	case constant.SingleChatType:
		// callback
		canSend, err := callbackBeforeSendSingleMsg(pb)
		if err != nil {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg failed", err.Error())
		}
		if !canSend {
			log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg result", canSend, "end rpc and return")
			return returnMsg(&replay, pb, 201, "callbackBeforeSendSingleMsg result stop rpc and return", "", 0)
		}
		isSend := modifyMessageByUserMessageReceiveOpt(pb.MsgData.RecvID, pb.MsgData.SendID, constant.SingleChatType, pb)
		if isSend {
			msgToMQSingle.MsgData = pb.MsgData
			log.NewInfo(msgToMQSingle.OperationID, msgToMQSingle)
			err1 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.RecvID)
			if err1 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err:RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			err2 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.SendID)
			if err2 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err:SendID", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		// callback
		if err := callbackAfterSendSingleMsg(pb); err != nil {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackAfterSendSingleMsg failed", err.Error())
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)
	case constant.GroupChatType:
		// callback
		canSend, err := callbackBeforeSendGroupMsg(pb)
		if err != nil {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendGroupMsg failed", err.Error())
		}
		if !canSend {
			log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendGroupMsg result", canSend, "end rpc and return")
			return returnMsg(&replay, pb, 201, "callbackBeforeSendGroupMsg result stop rpc and return", "", 0)
		}
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
		client := pbGroup.NewGroupClient(etcdConn)
		req := &pbGroup.GetGroupAllMemberReq{
			GroupID:     pb.MsgData.GroupID,
			OperationID: pb.OperationID,
		}
		reply, err := client.GetGroupAllMember(context.Background(), req)
		if err != nil {
			log.Error(pb.Token, pb.OperationID, "rpc send_msg getGroupInfo failed, err = %s", err.Error())
			return returnMsg(&replay, pb, 201, err.Error(), "", 0)
		}
		if reply.ErrCode != 0 {
			log.Error(pb.Token, pb.OperationID, "rpc send_msg getGroupInfo failed, err = %s", reply.ErrMsg)
			return returnMsg(&replay, pb, reply.ErrCode, reply.ErrMsg, "", 0)
		}
		memberUserIDList := func(all []*sdk_ws.GroupMemberFullInfo) (result []string) {
			for _, v := range all {
				result = append(result, v.UserID)
			}
			return result
		}(reply.MemberList)
		log.Debug(pb.OperationID, "GetGroupAllMember userID list", memberUserIDList)
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
		groupID := pb.MsgData.GroupID
		//split  parallel send
		var wg sync.WaitGroup
		var sendTag bool
		var split = 10
		remain := len(memberUserIDList) % split
		for i := 0; i < len(memberUserIDList)/split; i++ {
			wg.Add(1)
			go func(list []string) {
				log.Debug(pb.OperationID, "split userID ", list)
				groupPB := pbChat.SendMsgReq{Token: pb.Token, OperationID: pb.OperationID, MsgData: &sdk_ws.MsgData{OfflinePushInfo: &sdk_ws.OfflinePushInfo{}}}
				*groupPB.MsgData = *pb.MsgData
				*groupPB.MsgData.OfflinePushInfo = *pb.MsgData.OfflinePushInfo
				msgToMQGroup := pbChat.MsgDataToMQ{Token: groupPB.Token, OperationID: groupPB.OperationID, MsgData: groupPB.MsgData}
				for _, v := range list {
					groupPB.MsgData.RecvID = v
					isSend := modifyMessageByUserMessageReceiveOpt(v, groupID, constant.GroupChatType, &groupPB)
					if isSend {
						msgToMQGroup.MsgData = groupPB.MsgData
						log.Debug(groupPB.OperationID, "sendMsgToKafka, ", v, groupID, msgToMQGroup.String())
						err := rpc.sendMsgToKafka(&msgToMQGroup, v)
						if err != nil {
							log.NewError(msgToMQGroup.OperationID, "kafka send msg err:UserId", v, msgToMQGroup.String())
						} else {
							sendTag = true
						}
					} else {
						log.Debug(groupPB.OperationID, "not sendMsgToKafka, ", v)
					}
				}
				wg.Done()
			}(memberUserIDList[i*split : (i+1)*split])
		}
		if remain > 0 {
			wg.Add(1)
			go func(list []string) {
				log.Debug(pb.OperationID, "split userID ", list)
				groupPB := pbChat.SendMsgReq{Token: pb.Token, OperationID: pb.OperationID, MsgData: &sdk_ws.MsgData{OfflinePushInfo: &sdk_ws.OfflinePushInfo{}}}
				*groupPB.MsgData = *pb.MsgData
				*groupPB.MsgData.OfflinePushInfo = *pb.MsgData.OfflinePushInfo
				msgToMQGroup := pbChat.MsgDataToMQ{Token: groupPB.Token, OperationID: groupPB.OperationID, MsgData: groupPB.MsgData}
				for _, v := range list {
					groupPB.MsgData.RecvID = v
					isSend := modifyMessageByUserMessageReceiveOpt(v, groupID, constant.GroupChatType, &groupPB)
					if isSend {
						msgToMQGroup.MsgData = groupPB.MsgData
						log.Debug(groupPB.OperationID, "sendMsgToKafka, ", v, groupID, msgToMQGroup.String())
						err := rpc.sendMsgToKafka(&msgToMQGroup, v)
						if err != nil {
							log.NewError(msgToMQGroup.OperationID, "kafka send msg err:UserId", v, msgToMQGroup.String())
						} else {
							sendTag = true
						}
					} else {
						log.Debug(groupPB.OperationID, "not sendMsgToKafka, ", v)
					}
				}
				wg.Done()
			}(memberUserIDList[split*(len(memberUserIDList)/split):])
		}
		wg.Wait()

		log.Info(msgToMQSingle.OperationID, "addUidList", addUidList)
		for _, v := range addUidList {
			pb.MsgData.RecvID = v
			isSend := modifyMessageByUserMessageReceiveOpt(v, groupID, constant.GroupChatType, pb)
			log.Info(msgToMQSingle.OperationID, "isSend", isSend)
			if isSend {
				msgToMQSingle.MsgData = pb.MsgData
				err := rpc.sendMsgToKafka(&msgToMQSingle, v)
				if err != nil {
					log.NewError(msgToMQSingle.OperationID, "kafka send msg err:UserId", v, msgToMQSingle.String())
				} else {
					sendTag = true
				}
			}
		}
		// callback
		if err := callbackAfterSendGroupMsg(pb); err != nil {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackAfterSendGroupMsg failed", err.Error())
		}
		if !sendTag {
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
					etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImConversationName)
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
						etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImConversationName)
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
			return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)

		}
	case constant.NotificationChatType:
		msgToMQSingle.MsgData = pb.MsgData
		log.NewInfo(msgToMQSingle.OperationID, msgToMQSingle)
		err1 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.RecvID)
		if err1 != nil {
			log.NewError(msgToMQSingle.OperationID, "kafka send msg err:RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		}

		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			err2 := rpc.sendMsgToKafka(&msgToMQSingle, msgToMQSingle.MsgData.SendID)
			if err2 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err:SendID", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)
	default:
		return returnMsg(&replay, pb, 203, "unkonwn sessionType", "", 0)
	}
}

func (rpc *rpcChat) sendMsgToKafka(m *pbChat.MsgDataToMQ, key string) error {
	pid, offset, err := rpc.producer.SendMessage(m, key)
	if err != nil {
		log.ErrorByKv("kafka send failed", m.OperationID, "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "key", key)
	}
	return err
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
	conversationID := utils.GetConversationIDBySessionType(sourceID, sessionType)
	opt, err := db.DB.GetSingleConversationRecvMsgOpt(userID, conversationID)
	if err != nil && err != redis.ErrNil {
		log.NewError(pb.OperationID, "GetSingleConversationMsgOpt from redis err", conversationID, pb.String(), err.Error())
		return true
	}
	switch opt {
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

type NotificationMsg struct {
	SendID      string
	RecvID      string
	Content     []byte //  open_im_sdk.TipsComm
	MsgFrom     int32
	ContentType int32
	SessionType int32
	OperationID string
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
	switch n.SessionType {
	case constant.GroupChatType:
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
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	client := pbChat.NewChatClient(etcdConn)
	reply, err := client.SendMsg(context.Background(), &req)
	if err != nil {
		log.NewError(req.OperationID, "SendMsg rpc failed, ", req.String(), err.Error())
	} else if reply.ErrCode != 0 {
		log.NewError(req.OperationID, "SendMsg rpc failed, ", req.String(), reply.ErrCode, reply.ErrMsg)
	}
}

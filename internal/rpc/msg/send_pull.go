package msg

import (
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbChat "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdkws"
	"context"
	go_redis "github.com/go-redis/redis/v8"
	"strings"
	"sync"
	"time"
)

func (m *msgServer) SendMsg(ctx context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, error) {
	replay := pbChat.SendMsgResp{}

	flag, errCode, errMsg := isMessageHasReadEnabled(pb)
	if !flag {
		return returnMsg(&replay, pb, errCode, errMsg, "", 0)
	}
	t1 := time.Now()
	m.encapsulateMsgData(pb.MsgData)
	log.Debug(pb.OperationID, "encapsulateMsgData ", " cost time: ", time.Since(t1))
	msgToMQSingle := pbChat.MsgDataToMQ{Token: pb.Token, OperationID: pb.OperationID, MsgData: pb.MsgData}
	// callback
	t1 = time.Now()
	callbackResp := callbackMsgModify(pb)
	log.Debug(pb.OperationID, "callbackMsgModify ", callbackResp, "cost time: ", time.Since(t1))
	if callbackResp.ErrCode != 0 {
		log.Error(pb.OperationID, utils.GetSelfFuncName(), "callbackMsgModify resp: ", callbackResp)
	}
	log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackResp: ", callbackResp)
	if callbackResp.ActionCode != constant.ActionAllow {
		if callbackResp.ErrCode == 0 {
			callbackResp.ErrCode = 201
		}
		log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackMsgModify result", "end rpc and return", pb.MsgData)
		return returnMsg(&replay, pb, int32(callbackResp.ErrCode), callbackResp.ErrMsg, "", 0)
	}
	switch pb.MsgData.SessionType {
	case constant.SingleChatType:
		promePkg.PromeInc(promePkg.SingleChatMsgRecvSuccessCounter)
		// callback
		t1 = time.Now()
		callbackResp := callbackBeforeSendSingleMsg(pb)
		log.Debug(pb.OperationID, "callbackBeforeSendSingleMsg ", " cost time: ", time.Since(t1))
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg resp: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			if callbackResp.ErrCode == 0 {
				callbackResp.ErrCode = 201
			}
			log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
			promePkg.PromeInc(promePkg.SingleChatMsgProcessFailedCounter)
			return returnMsg(&replay, pb, int32(callbackResp.ErrCode), callbackResp.ErrMsg, "", 0)
		}
		t1 = time.Now()
		flag, errCode, errMsg, _ = rpc.messageVerification(ctx, pb)
		log.Debug(pb.OperationID, "messageVerification ", flag, " cost time: ", time.Since(t1))
		if !flag {
			return returnMsg(&replay, pb, errCode, errMsg, "", 0)
		}
		t1 = time.Now()
		isSend := modifyMessageByUserMessageReceiveOpt(pb.MsgData.RecvID, pb.MsgData.SendID, constant.SingleChatType, pb)
		log.Info(pb.OperationID, "modifyMessageByUserMessageReceiveOpt ", " cost time: ", time.Since(t1))
		if isSend {
			msgToMQSingle.MsgData = pb.MsgData
			log.NewInfo(msgToMQSingle.OperationID, msgToMQSingle)
			t1 = time.Now()
			err1 := rpc.sendMsgToWriter(ctx, &msgToMQSingle, msgToMQSingle.MsgData.RecvID, constant.OnlineStatus)
			log.Info(pb.OperationID, "sendMsgToWriter ", " cost time: ", time.Since(t1))
			if err1 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err :RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String(), err1.Error())
				promePkg.PromeInc(promePkg.SingleChatMsgProcessFailedCounter)
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			t1 = time.Now()
			err2 := rpc.sendMsgToWriter(ctx, &msgToMQSingle, msgToMQSingle.MsgData.SendID, constant.OnlineStatus)
			log.Info(pb.OperationID, "sendMsgToWriter ", " cost time: ", time.Since(t1))
			if err2 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err:SendID", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				promePkg.PromeInc(promePkg.SingleChatMsgProcessFailedCounter)
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		// callback
		t1 = time.Now()
		callbackResp = callbackAfterSendSingleMsg(pb)
		log.Info(pb.OperationID, "callbackAfterSendSingleMsg ", " cost time: ", time.Since(t1))
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackAfterSendSingleMsg resp: ", callbackResp)
		}
		promePkg.PromeInc(promePkg.SingleChatMsgProcessSuccessCounter)
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)
	case constant.GroupChatType:
		// callback
		promePkg.PromeInc(promePkg.GroupChatMsgRecvSuccessCounter)
		callbackResp := callbackBeforeSendGroupMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendGroupMsg resp:", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			if callbackResp.ErrCode == 0 {
				callbackResp.ErrCode = 201
			}
			log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
			promePkg.PromeInc(promePkg.GroupChatMsgProcessFailedCounter)
			return returnMsg(&replay, pb, int32(callbackResp.ErrCode), callbackResp.ErrMsg, "", 0)
		}
		var memberUserIDList []string
		if flag, errCode, errMsg, memberUserIDList = rpc.messageVerification(ctx, pb); !flag {
			promePkg.PromeInc(promePkg.GroupChatMsgProcessFailedCounter)
			return returnMsg(&replay, pb, errCode, errMsg, "", 0)
		}
		log.Debug(pb.OperationID, "GetGroupAllMember userID list", memberUserIDList, "len: ", len(memberUserIDList))
		var addUidList []string
		switch pb.MsgData.ContentType {
		case constant.MemberKickedNotification:
			var tips sdkws.TipsComm
			var memberKickedTips sdkws.MemberKickedTips
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
		t1 = time.Now()

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
				go rpc.sendMsgToGroupOptimization(ctx, v[i*split:(i+1)*split], tmp, k, &sendTag, &wg)
			}
			if remain > 0 {
				wg.Add(1)
				tmp := valueCopy(pb)
				//	go rpc.sendMsgToGroup(v[split*(len(v)/split):], *pb, k, &sendTag, &wg)
				go rpc.sendMsgToGroupOptimization(ctx, v[split*(len(v)/split):], tmp, k, &sendTag, &wg)
			}
		}
		log.Debug(pb.OperationID, "send msg cost time22 ", time.Since(t1), pb.MsgData.ClientMsgID, "uidList : ", len(addUidList))
		//wg.Add(1)
		//go rpc.sendMsgToGroup(addUidList, *pb, constant.OnlineStatus, &sendTag, &wg)
		wg.Wait()
		t1 = time.Now()
		// callback
		callbackResp = callbackAfterSendGroupMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackAfterSendGroupMsg resp: ", callbackResp)
		}
		if !sendTag {
			log.NewWarn(pb.OperationID, "send tag is ", sendTag)
			promePkg.PromeInc(promePkg.GroupChatMsgProcessFailedCounter)
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
					etcdConn, err := rpc.GetConn(ctx, config.Config.RpcRegisterName.OpenImConversationName)
					if err != nil {
						errMsg := pb.OperationID + "getcdv3.GetDefaultConn == nil"
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
						etcdConn := rpc.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImConversationName, pb.OperationID)
						if etcdConn == nil {
							errMsg := pb.OperationID + "getcdv3.GetDefaultConn == nil"
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
			log.Debug(pb.OperationID, "send msg cost time3 ", time.Since(t1), pb.MsgData.ClientMsgID)
			promePkg.PromeInc(promePkg.GroupChatMsgProcessSuccessCounter)
			return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)
		}
	case constant.NotificationChatType:
		t1 = time.Now()
		msgToMQSingle.MsgData = pb.MsgData
		log.NewInfo(msgToMQSingle.OperationID, msgToMQSingle)
		err1 := rpc.sendMsgToWriter(ctx, &msgToMQSingle, msgToMQSingle.MsgData.RecvID, constant.OnlineStatus)
		if err1 != nil {
			log.NewError(msgToMQSingle.OperationID, "kafka send msg err:RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		}

		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			err2 := rpc.sendMsgToWriter(ctx, &msgToMQSingle, msgToMQSingle.MsgData.SendID, constant.OnlineStatus)
			if err2 != nil {
				log.NewError(msgToMQSingle.OperationID, "kafka send msg err:SendID", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}

		log.Debug(pb.OperationID, "send msg cost time ", time.Since(t1), pb.MsgData.ClientMsgID)
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)
	case constant.SuperGroupChatType:
		promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgRecvSuccessCounter)
		// callback
		callbackResp := callbackBeforeSendGroupMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSuperGroupMsg resp: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			if callbackResp.ErrCode == 0 {
				callbackResp.ErrCode = 201
			}
			promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgProcessFailedCounter)
			log.NewDebug(pb.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSuperGroupMsg result", "end rpc and return", callbackResp)
			return returnMsg(&replay, pb, int32(callbackResp.ErrCode), callbackResp.ErrMsg, "", 0)
		}
		if flag, errCode, errMsg, _ = rpc.messageVerification(ctx, pb); !flag {
			promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgProcessFailedCounter)
			return returnMsg(&replay, pb, errCode, errMsg, "", 0)
		}
		msgToMQSingle.MsgData = pb.MsgData
		log.NewInfo(msgToMQSingle.OperationID, msgToMQSingle)
		err1 := rpc.sendMsgToWriter(ctx, &msgToMQSingle, msgToMQSingle.MsgData.GroupID, constant.OnlineStatus)
		if err1 != nil {
			log.NewError(msgToMQSingle.OperationID, "kafka send msg err:RecvID", msgToMQSingle.MsgData.RecvID, msgToMQSingle.String())
			promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgProcessFailedCounter)
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		}
		// callback
		callbackResp = callbackAfterSendGroupMsg(pb)
		if callbackResp.ErrCode != 0 {
			log.NewError(pb.OperationID, utils.GetSelfFuncName(), "callbackAfterSendSuperGroupMsg resp: ", callbackResp)
		}
		promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgProcessSuccessCounter)
		return returnMsg(&replay, pb, 0, "", msgToMQSingle.MsgData.ServerMsgID, msgToMQSingle.MsgData.SendTime)

	default:
		return returnMsg(&replay, pb, 203, "unknown sessionType", "", 0)
	}
}

func (rpc *rpcChat) GetMaxAndMinSeq(_ context.Context, in *sdkws.GetMaxAndMinSeqReq) (*sdkws.GetMaxAndMinSeqResp, error) {
	log.NewInfo(in.OperationID, "rpc getMaxAndMinSeq is arriving", in.String())
	resp := new(sdkws.GetMaxAndMinSeqResp)
	m := make(map[string]*sdkws.MaxAndMinSeq)
	var maxSeq, minSeq uint64
	var err1, err2 error
	maxSeq, err1 = commonDB.DB.GetUserMaxSeq(in.UserID)
	minSeq, err2 = commonDB.DB.GetUserMinSeq(in.UserID)
	if (err1 != nil && err1 != go_redis.Nil) || (err2 != nil && err2 != go_redis.Nil) {
		log.NewError(in.OperationID, "getMaxSeq from redis error", in.String())
		if err1 != nil {
			log.NewError(in.OperationID, utils.GetSelfFuncName(), err1.Error())
		}
		if err2 != nil {
			log.NewError(in.OperationID, utils.GetSelfFuncName(), err2.Error())
		}
		resp.ErrCode = 200
		resp.ErrMsg = "redis get err"
		return resp, nil
	}
	resp.MaxSeq = uint32(maxSeq)
	resp.MinSeq = uint32(minSeq)
	for _, groupID := range in.GroupIDList {
		x := new(sdkws.MaxAndMinSeq)
		maxSeq, _ := commonDB.DB.GetGroupMaxSeq(groupID)
		minSeq, _ := commonDB.DB.GetGroupUserMinSeq(groupID, in.UserID)
		x.MaxSeq = uint32(maxSeq)
		x.MinSeq = uint32(minSeq)
		m[groupID] = x
	}
	resp.GroupMaxAndMinSeq = m
	return resp, nil
}

func (rpc *rpcChat) PullMessageBySeqList(_ context.Context, in *sdkws.PullMessageBySeqListReq) (*sdkws.PullMessageBySeqListResp, error) {
	log.NewInfo(in.OperationID, "rpc PullMessageBySeqList is arriving", in.String())
	resp := new(sdkws.PullMessageBySeqListResp)
	m := make(map[string]*sdkws.MsgDataList)
	redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(in.UserID, in.SeqList, in.OperationID)
	if err != nil {
		if err != go_redis.Nil {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
			log.Error(in.OperationID, "get message from redis exception", err.Error(), failedSeqList)
		} else {
			log.Debug(in.OperationID, "get message from redis is nil", failedSeqList)
		}
		msgList, err1 := commonDB.DB.GetMsgBySeqListMongo2(in.UserID, failedSeqList, in.OperationID)
		if err1 != nil {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
			log.Error(in.OperationID, "PullMessageBySeqList data error", in.String(), err1.Error())
			resp.ErrCode = 201
			resp.ErrMsg = err1.Error()
			return resp, nil
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
			redisMsgList = append(redisMsgList, msgList...)
			resp.List = redisMsgList
		}
	} else {
		promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
		resp.List = redisMsgList
	}

	for k, v := range in.GroupSeqList {
		x := new(sdkws.MsgDataList)
		redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(k, v.SeqList, in.OperationID)
		if err != nil {
			if err != go_redis.Nil {
				promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
				log.Error(in.OperationID, "get message from redis exception", err.Error(), failedSeqList)
			} else {
				log.Debug(in.OperationID, "get message from redis is nil", failedSeqList)
			}
			msgList, err1 := commonDB.DB.GetSuperGroupMsgBySeqListMongo(k, failedSeqList, in.OperationID)
			if err1 != nil {
				promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
				log.Error(in.OperationID, "PullMessageBySeqList data error", in.String(), err1.Error())
				resp.ErrCode = 201
				resp.ErrMsg = err1.Error()
				return resp, nil
			} else {
				promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
				redisMsgList = append(redisMsgList, msgList...)
				x.MsgDataList = redisMsgList
				m[k] = x
			}
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
			x.MsgDataList = redisMsgList
			m[k] = x
		}
	}
	resp.GroupMsgDataList = m
	return resp, nil
}

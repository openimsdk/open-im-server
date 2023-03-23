package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	go_redis "github.com/go-redis/redis/v8"

	"time"
)

func (rpc *rpcChat) SetMessageReactionExtensions(ctx context.Context, req *msg.SetMessageReactionExtensionsReq) (resp *msg.SetMessageReactionExtensionsResp, err error) {
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rpc args is:", req.String())
	var rResp msg.SetMessageReactionExtensionsResp
	rResp.ClientMsgID = req.ClientMsgID
	rResp.MsgFirstModifyTime = req.MsgFirstModifyTime
	callbackResp := callbackSetMessageReactionExtensions(req)
	if callbackResp.ActionCode != constant.ActionAllow || callbackResp.ErrCode != 0 {
		rResp.ErrCode = int32(callbackResp.ErrCode)
		rResp.ErrMsg = callbackResp.ErrMsg
		for _, value := range req.ReactionExtensionList {
			temp := new(msg.KeyValueResp)
			temp.KeyValue = value
			temp.ErrMsg = callbackResp.ErrMsg
			temp.ErrCode = 100
			rResp.Result = append(rResp.Result, temp)
		}
		return &rResp, nil
	}
	//if ExternalExtension
	if req.IsExternalExtensions {
		var isHistory bool
		if req.IsReact {
			isHistory = false
		} else {
			isHistory = true
		}
		rResp.MsgFirstModifyTime = callbackResp.MsgFirstModifyTime
		rResp.Result = callbackResp.ResultReactionExtensionList
		ExtendMessageUpdatedNotification(req.OperationID, req.OpUserID, req.SourceID, req.OpUserIDPlatformID, req.SessionType, req, &rResp, isHistory, false)
		return &rResp, nil
	}
	for _, v := range callbackResp.ResultReactionExtensionList {
		if v.ErrCode == 0 {
			req.ReactionExtensionList[v.KeyValue.TypeKey] = v.KeyValue
		} else {
			delete(req.ReactionExtensionList, v.KeyValue.TypeKey)
			rResp.Result = append(rResp.Result, v)
		}
	}
	isExists, err := db.DB.JudgeMessageReactionEXISTS(req.ClientMsgID, req.SessionType)
	if err != nil {
		rResp.ErrCode = 100
		rResp.ErrMsg = err.Error()
		for _, value := range req.ReactionExtensionList {
			temp := new(msg.KeyValueResp)
			temp.KeyValue = value
			temp.ErrMsg = err.Error()
			temp.ErrCode = 100
			rResp.Result = append(rResp.Result, temp)
		}
		return &rResp, nil
	}

	if !isExists {
		if !req.IsReact {
			log.Debug(req.OperationID, "redis handle firstly", req.String())
			rResp.MsgFirstModifyTime = utils.GetCurrentTimestampByMill()
			for k, v := range req.ReactionExtensionList {
				err := rpc.dMessageLocker.LockMessageTypeKey(req.ClientMsgID, k)
				if err != nil {
					setKeyResultInfo(&rResp, 100, err.Error(), req.ClientMsgID, k, v)
					continue
				}
				v.LatestUpdateTime = utils.GetCurrentTimestampByMill()
				newerr := db.DB.SetMessageTypeKeyValue(req.ClientMsgID, req.SessionType, k, utils.StructToJsonString(v))
				if newerr != nil {
					setKeyResultInfo(&rResp, 201, newerr.Error(), req.ClientMsgID, k, v)
					continue
				}
				setKeyResultInfo(&rResp, 0, "", req.ClientMsgID, k, v)
			}
			rResp.IsReact = true
			_, err := db.DB.SetMessageReactionExpire(req.ClientMsgID, req.SessionType, time.Duration(24*3)*time.Hour)
			if err != nil {
				log.Error(req.OperationID, "SetMessageReactionExpire err:", err.Error(), req.String())
			}
		} else {
			err := rpc.dMessageLocker.LockGlobalMessage(req.ClientMsgID)
			if err != nil {
				rResp.ErrCode = 100
				rResp.ErrMsg = err.Error()
				for _, value := range req.ReactionExtensionList {
					temp := new(msg.KeyValueResp)
					temp.KeyValue = value
					temp.ErrMsg = err.Error()
					temp.ErrCode = 100
					rResp.Result = append(rResp.Result, temp)
				}
				return &rResp, nil
			}
			mongoValue, err := db.DB.GetExtendMsg(req.SourceID, req.SessionType, req.ClientMsgID, req.MsgFirstModifyTime)
			if err != nil {
				rResp.ErrCode = 200
				rResp.ErrMsg = err.Error()
				for _, value := range req.ReactionExtensionList {
					temp := new(msg.KeyValueResp)
					temp.KeyValue = value
					temp.ErrMsg = err.Error()
					temp.ErrCode = 100
					rResp.Result = append(rResp.Result, temp)
				}
				return &rResp, nil
			}
			setValue := make(map[string]*server_api_params.KeyValue)
			for k, v := range req.ReactionExtensionList {

				temp := new(server_api_params.KeyValue)
				if vv, ok := mongoValue.ReactionExtensionList[k]; ok {
					utils.CopyStructFields(temp, &vv)
					if v.LatestUpdateTime != vv.LatestUpdateTime {
						setKeyResultInfo(&rResp, 300, "message have update", req.ClientMsgID, k, temp)
						continue
					}
				}
				temp.TypeKey = k
				temp.Value = v.Value
				temp.LatestUpdateTime = utils.GetCurrentTimestampByMill()
				setValue[k] = temp
			}
			err = db.DB.InsertOrUpdateReactionExtendMsgSet(req.SourceID, req.SessionType, req.ClientMsgID, req.MsgFirstModifyTime, setValue)
			if err != nil {
				for _, value := range setValue {
					temp := new(msg.KeyValueResp)
					temp.KeyValue = value
					temp.ErrMsg = err.Error()
					temp.ErrCode = 100
					rResp.Result = append(rResp.Result, temp)
				}
			} else {
				for _, value := range setValue {
					temp := new(msg.KeyValueResp)
					temp.KeyValue = value
					rResp.Result = append(rResp.Result, temp)
				}
			}
			lockErr := rpc.dMessageLocker.UnLockGlobalMessage(req.ClientMsgID)
			if lockErr != nil {
				log.Error(req.OperationID, "UnLockGlobalMessage err:", lockErr.Error())
			}
		}

	} else {
		log.Debug(req.OperationID, "redis handle secondly", req.String())

		for k, v := range req.ReactionExtensionList {
			err := rpc.dMessageLocker.LockMessageTypeKey(req.ClientMsgID, k)
			if err != nil {
				setKeyResultInfo(&rResp, 100, err.Error(), req.ClientMsgID, k, v)
				continue
			}
			redisValue, err := db.DB.GetMessageTypeKeyValue(req.ClientMsgID, req.SessionType, k)
			if err != nil && err != go_redis.Nil {
				setKeyResultInfo(&rResp, 200, err.Error(), req.ClientMsgID, k, v)
				continue
			}
			temp := new(server_api_params.KeyValue)
			utils.JsonStringToStruct(redisValue, temp)
			if v.LatestUpdateTime != temp.LatestUpdateTime {
				setKeyResultInfo(&rResp, 300, "message have update", req.ClientMsgID, k, temp)
				continue
			} else {
				v.LatestUpdateTime = utils.GetCurrentTimestampByMill()
				newerr := db.DB.SetMessageTypeKeyValue(req.ClientMsgID, req.SessionType, k, utils.StructToJsonString(v))
				if newerr != nil {
					setKeyResultInfo(&rResp, 201, newerr.Error(), req.ClientMsgID, k, temp)
					continue
				}
				setKeyResultInfo(&rResp, 0, "", req.ClientMsgID, k, v)
			}

		}
	}
	if !isExists {
		if !req.IsReact {
			ExtendMessageUpdatedNotification(req.OperationID, req.OpUserID, req.SourceID, req.OpUserIDPlatformID, req.SessionType, req, &rResp, true, true)
		} else {
			ExtendMessageUpdatedNotification(req.OperationID, req.OpUserID, req.SourceID, req.OpUserIDPlatformID, req.SessionType, req, &rResp, false, false)
		}
	} else {
		ExtendMessageUpdatedNotification(req.OperationID, req.OpUserID, req.SourceID, req.OpUserIDPlatformID, req.SessionType, req, &rResp, false, true)
	}
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rpc return is:", rResp.String())
	return &rResp, nil

}
func setKeyResultInfo(r *msg.SetMessageReactionExtensionsResp, errCode int32, errMsg, clientMsgID, typeKey string, keyValue *server_api_params.KeyValue) {
	temp := new(msg.KeyValueResp)
	temp.KeyValue = keyValue
	temp.ErrCode = errCode
	temp.ErrMsg = errMsg
	r.Result = append(r.Result, temp)
	_ = db.DB.UnLockMessageTypeKey(clientMsgID, typeKey)
}
func setDeleteKeyResultInfo(r *msg.DeleteMessageListReactionExtensionsResp, errCode int32, errMsg, clientMsgID, typeKey string, keyValue *server_api_params.KeyValue) {
	temp := new(msg.KeyValueResp)
	temp.KeyValue = keyValue
	temp.ErrCode = errCode
	temp.ErrMsg = errMsg
	r.Result = append(r.Result, temp)
	_ = db.DB.UnLockMessageTypeKey(clientMsgID, typeKey)
}

func (rpc *rpcChat) GetMessageListReactionExtensions(ctx context.Context, req *msg.GetMessageListReactionExtensionsReq) (resp *msg.GetMessageListReactionExtensionsResp, err error) {
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rpc args is:", req.String())
	var rResp msg.GetMessageListReactionExtensionsResp
	if req.IsExternalExtensions {
		callbackResp := callbackGetMessageListReactionExtensions(req)
		if callbackResp.ActionCode != constant.ActionAllow || callbackResp.ErrCode != 0 {
			rResp.ErrCode = int32(callbackResp.ErrCode)
			rResp.ErrMsg = callbackResp.ErrMsg
			return &rResp, nil
		} else {
			rResp.SingleMessageResult = callbackResp.MessageResultList
			return &rResp, nil
		}
	}
	for _, messageValue := range req.MessageReactionKeyList {
		var oneMessage msg.SingleMessageExtensionResult
		oneMessage.ClientMsgID = messageValue.ClientMsgID

		isExists, err := db.DB.JudgeMessageReactionEXISTS(messageValue.ClientMsgID, req.SessionType)
		if err != nil {
			rResp.ErrCode = 100
			rResp.ErrMsg = err.Error()
			return &rResp, nil
		}
		if isExists {
			redisValue, err := db.DB.GetOneMessageAllReactionList(messageValue.ClientMsgID, req.SessionType)
			if err != nil {
				oneMessage.ErrCode = 100
				oneMessage.ErrMsg = err.Error()
				rResp.SingleMessageResult = append(rResp.SingleMessageResult, &oneMessage)
				continue
			}
			keyMap := make(map[string]*server_api_params.KeyValue)

			for k, v := range redisValue {
				temp := new(server_api_params.KeyValue)
				utils.JsonStringToStruct(v, temp)
				keyMap[k] = temp
			}
			oneMessage.ReactionExtensionList = keyMap

		} else {
			mongoValue, err := db.DB.GetExtendMsg(req.SourceID, req.SessionType, messageValue.ClientMsgID, messageValue.MsgFirstModifyTime)
			if err != nil {
				oneMessage.ErrCode = 100
				oneMessage.ErrMsg = err.Error()
				rResp.SingleMessageResult = append(rResp.SingleMessageResult, &oneMessage)
				continue
			}
			keyMap := make(map[string]*server_api_params.KeyValue)

			for k, v := range mongoValue.ReactionExtensionList {
				temp := new(server_api_params.KeyValue)
				temp.TypeKey = v.TypeKey
				temp.Value = v.Value
				temp.LatestUpdateTime = v.LatestUpdateTime
				keyMap[k] = temp
			}
			oneMessage.ReactionExtensionList = keyMap
		}
		rResp.SingleMessageResult = append(rResp.SingleMessageResult, &oneMessage)
	}
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rpc return is:", rResp.String())
	return &rResp, nil

}

func (rpc *rpcChat) AddMessageReactionExtensions(ctx context.Context, req *msg.AddMessageReactionExtensionsReq) (resp *msg.AddMessageReactionExtensionsResp, err error) {
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rpc args is:", req.String())
	var rResp msg.AddMessageReactionExtensionsResp
	rResp.ClientMsgID = req.ClientMsgID
	rResp.MsgFirstModifyTime = req.MsgFirstModifyTime
	callbackResp := callbackAddMessageReactionExtensions(req)
	if callbackResp.ActionCode != constant.ActionAllow || callbackResp.ErrCode != 0 {
		rResp.ErrCode = int32(callbackResp.ErrCode)
		rResp.ErrMsg = callbackResp.ErrMsg
		for _, value := range callbackResp.ResultReactionExtensionList {
			temp := new(msg.KeyValueResp)
			temp.KeyValue = value.KeyValue
			temp.ErrMsg = value.ErrMsg
			temp.ErrCode = value.ErrCode
			rResp.Result = append(rResp.Result, temp)
		}
		return &rResp, nil
	}

	//if !req.IsExternalExtensions {
	//	rResp.ErrCode = 200
	//	rResp.ErrMsg = "only extenalextensions message can be used"
	//	for _, value := range req.ReactionExtensionList {
	//		temp := new(msg.KeyValueResp)
	//		temp.KeyValue = value
	//		temp.ErrMsg = callbackResp.ErrMsg
	//		temp.ErrCode = 100
	//		rResp.Result = append(rResp.Result, temp)
	//	}
	//	return &rResp, nil
	//}
	//if ExternalExtension
	var isHistory bool
	if req.IsReact {
		isHistory = false
	} else {
		isHistory = true
	}
	rResp.MsgFirstModifyTime = callbackResp.MsgFirstModifyTime
	rResp.Result = callbackResp.ResultReactionExtensionList
	rResp.IsReact = callbackResp.IsReact
	ExtendMessageAddedNotification(req.OperationID, req.OpUserID, req.SourceID, req.OpUserIDPlatformID, req.SessionType, req, &rResp, isHistory, false)
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rpc return is:", resp.String())
	return &rResp, nil
}

func (rpc *rpcChat) DeleteMessageReactionExtensions(ctx context.Context, req *msg.DeleteMessageListReactionExtensionsReq) (resp *msg.DeleteMessageListReactionExtensionsResp, err error) {
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rpc args is:", req.String())
	var rResp msg.DeleteMessageListReactionExtensionsResp
	callbackResp := callbackDeleteMessageReactionExtensions(req)
	if callbackResp.ActionCode != constant.ActionAllow || callbackResp.ErrCode != 0 {
		rResp.ErrCode = int32(callbackResp.ErrCode)
		rResp.ErrMsg = callbackResp.ErrMsg
		for _, value := range req.ReactionExtensionList {
			temp := new(msg.KeyValueResp)
			temp.KeyValue = value
			temp.ErrMsg = callbackResp.ErrMsg
			temp.ErrCode = 100
			rResp.Result = append(rResp.Result, temp)
		}
		return &rResp, nil
	}
	//if ExternalExtension
	if req.IsExternalExtensions {
		rResp.Result = callbackResp.ResultReactionExtensionList
		ExtendMessageDeleteNotification(req.OperationID, req.OpUserID, req.SourceID, req.OpUserIDPlatformID, req.SessionType, req, &rResp, false, false)
		return &rResp, nil

	}
	for _, v := range callbackResp.ResultReactionExtensionList {
		if v.ErrCode != 0 {
			func(req *[]*server_api_params.KeyValue, typeKey string) {
				for i := 0; i < len(*req); i++ {
					if (*req)[i].TypeKey == typeKey {
						*req = append((*req)[:i], (*req)[i+1:]...)
					}
				}
			}(&req.ReactionExtensionList, v.KeyValue.TypeKey)
			rResp.Result = append(rResp.Result, v)
		}
	}
	isExists, err := db.DB.JudgeMessageReactionEXISTS(req.ClientMsgID, req.SessionType)
	if err != nil {
		rResp.ErrCode = 100
		rResp.ErrMsg = err.Error()
		for _, value := range req.ReactionExtensionList {
			temp := new(msg.KeyValueResp)
			temp.KeyValue = value
			temp.ErrMsg = err.Error()
			temp.ErrCode = 100
			rResp.Result = append(rResp.Result, temp)
		}
		return &rResp, nil
	}

	if isExists {
		log.Debug(req.OperationID, "redis handle this delete", req.String())
		for _, v := range req.ReactionExtensionList {
			err := rpc.dMessageLocker.LockMessageTypeKey(req.ClientMsgID, v.TypeKey)
			if err != nil {
				setDeleteKeyResultInfo(&rResp, 100, err.Error(), req.ClientMsgID, v.TypeKey, v)
				continue
			}

			redisValue, err := db.DB.GetMessageTypeKeyValue(req.ClientMsgID, req.SessionType, v.TypeKey)
			if err != nil && err != go_redis.Nil {
				setDeleteKeyResultInfo(&rResp, 200, err.Error(), req.ClientMsgID, v.TypeKey, v)
				continue
			}
			temp := new(server_api_params.KeyValue)
			utils.JsonStringToStruct(redisValue, temp)
			if v.LatestUpdateTime != temp.LatestUpdateTime {
				setDeleteKeyResultInfo(&rResp, 300, "message have update", req.ClientMsgID, v.TypeKey, temp)
				continue
			} else {
				newErr := db.DB.DeleteOneMessageKey(req.ClientMsgID, req.SessionType, v.TypeKey)
				if newErr != nil {
					setDeleteKeyResultInfo(&rResp, 201, newErr.Error(), req.ClientMsgID, v.TypeKey, temp)
					continue
				}
				setDeleteKeyResultInfo(&rResp, 0, "", req.ClientMsgID, v.TypeKey, v)
			}
		}
	} else {
		err := rpc.dMessageLocker.LockGlobalMessage(req.ClientMsgID)
		if err != nil {
			rResp.ErrCode = 100
			rResp.ErrMsg = err.Error()
			for _, value := range req.ReactionExtensionList {
				temp := new(msg.KeyValueResp)
				temp.KeyValue = value
				temp.ErrMsg = err.Error()
				temp.ErrCode = 100
				rResp.Result = append(rResp.Result, temp)
			}
			return &rResp, nil
		}
		mongoValue, err := db.DB.GetExtendMsg(req.SourceID, req.SessionType, req.ClientMsgID, req.MsgFirstModifyTime)
		if err != nil {
			rResp.ErrCode = 200
			rResp.ErrMsg = err.Error()
			for _, value := range req.ReactionExtensionList {
				temp := new(msg.KeyValueResp)
				temp.KeyValue = value
				temp.ErrMsg = err.Error()
				temp.ErrCode = 100
				rResp.Result = append(rResp.Result, temp)
			}
			return &rResp, nil
		}
		setValue := make(map[string]*server_api_params.KeyValue)
		for _, v := range req.ReactionExtensionList {

			temp := new(server_api_params.KeyValue)
			if vv, ok := mongoValue.ReactionExtensionList[v.TypeKey]; ok {
				utils.CopyStructFields(temp, &vv)
				if v.LatestUpdateTime != vv.LatestUpdateTime {
					setDeleteKeyResultInfo(&rResp, 300, "message have update", req.ClientMsgID, v.TypeKey, temp)
					continue
				}
			} else {
				setDeleteKeyResultInfo(&rResp, 400, "key not in", req.ClientMsgID, v.TypeKey, v)
				continue
			}
			temp.TypeKey = v.TypeKey
			setValue[v.TypeKey] = temp
		}
		err = db.DB.DeleteReactionExtendMsgSet(req.SourceID, req.SessionType, req.ClientMsgID, req.MsgFirstModifyTime, setValue)
		if err != nil {
			for _, value := range setValue {
				temp := new(msg.KeyValueResp)
				temp.KeyValue = value
				temp.ErrMsg = err.Error()
				temp.ErrCode = 100
				rResp.Result = append(rResp.Result, temp)
			}
		} else {
			for _, value := range setValue {
				temp := new(msg.KeyValueResp)
				temp.KeyValue = value
				rResp.Result = append(rResp.Result, temp)
			}
		}
		lockErr := rpc.dMessageLocker.UnLockGlobalMessage(req.ClientMsgID)
		if lockErr != nil {
			log.Error(req.OperationID, "UnLockGlobalMessage err:", lockErr.Error())
		}

	}
	ExtendMessageDeleteNotification(req.OperationID, req.OpUserID, req.SourceID, req.OpUserIDPlatformID, req.SessionType, req, &rResp, false, isExists)
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "rpc return is:", rResp.String())
	return &rResp, nil
}

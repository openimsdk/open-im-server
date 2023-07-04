// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func (m *msgServer) SetMessageReactionExtensions(
	ctx context.Context,
	req *msg.SetMessageReactionExtensionsReq,
) (resp *msg.SetMessageReactionExtensionsResp, err error) {
	//resp = &msg.SetMessageReactionExtensionsResp{}
	////resp.ClientMsgID = req.ClientMsgID
	////resp.MsgFirstModifyTime = req.MsgFirstModifyTime
	//
	//if err := CallbackSetMessageReactionExtensions(ctx, req); err != nil {
	//	return nil, err
	//}
	////if ExternalExtension
	//if req.IsExternalExtensions {
	//	resp.MsgFirstModifyTime = req.MsgFirstModifyTime
	// 	notification.ExtendMessageUpdatedNotification(req.OperationID, req.OpUserID, req.conversationID,
	// req.SessionType, req, &resp, !req.IsReact, false)
	//	return resp, nil
	//}
	//isExists, err := m.MsgDatabase.JudgeMessageReactionExist(ctx, req.ClientMsgID, req.SessionType)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if !isExists {
	//	if !req.IsReact {
	//		resp.MsgFirstModifyTime = utils.GetCurrentTimestampByMill()
	//		for k, v := range req.ReactionExtensions {
	//			err := m.MessageLocker.LockMessageTypeKey(ctx, req.ClientMsgID, k)
	//			if err != nil {
	//				return nil, err
	//			}
	//			v.LatestUpdateTime = utils.GetCurrentTimestampByMill()
	// 			if err := m.MsgDatabase.SetMessageTypeKeyValue(ctx, req.ClientMsgID, req.SessionType, k,
	// utils.StructToJsonString(v)); err != nil {
	//				return nil, err
	//			}
	//		}
	//		resp.IsReact = true
	// 		_, err := m.MsgDatabase.SetMessageReactionExpire(ctx, req.ClientMsgID, req.SessionType,
	// time.Duration(24*3)*time.Hour)
	//		if err != nil {
	//			return nil, err
	//		}
	//	} else {
	//		err := m.MessageLocker.LockGlobalMessage(ctx, req.ClientMsgID)
	//		if err != nil {
	//			return nil, err
	//		}
	// 		mongoValue, err := m.MsgDatabase.GetExtendMsg(ctx, req.conversationID, req.SessionType, req.ClientMsgID,
	// req.MsgFirstModifyTime)
	//		if err != nil {
	//			return nil, err
	//		}
	//		setValue := make(map[string]*sdkws.KeyValue)
	//		for k, v := range req.ReactionExtensions {
	//
	//			temp := new(sdkws.KeyValue)
	//			if vv, ok := mongoValue.ReactionExtensions[k]; ok {
	//				utils.CopyStructFields(temp, &vv)
	//				if v.LatestUpdateTime != vv.LatestUpdateTime {
	//					setKeyResultInfo(&resp, 300, "message have update", req.ClientMsgID, k, temp)
	//					continue
	//				}
	//			}
	//			temp.TypeKey = k
	//			temp.Value = v.Value
	//			temp.LatestUpdateTime = utils.GetCurrentTimestampByMill()
	//			setValue[k] = temp
	//		}
	// 		err = db.DB.InsertOrUpdateReactionExtendMsgSet(req.conversationID, req.SessionType, req.ClientMsgID,
	// req.MsgFirstModifyTime, setValue)
	//		if err != nil {
	//			for _, value := range setValue {
	//				temp := new(msg.KeyValueResp)
	//				temp.KeyValue = value
	//				temp.ErrMsg = err.Error()
	//				temp.ErrCode = 100
	//				resp.Result = append(resp.Result, temp)
	//			}
	//		} else {
	//			for _, value := range setValue {
	//				temp := new(msg.KeyValueResp)
	//				temp.KeyValue = value
	//				resp.Result = append(resp.Result, temp)
	//			}
	//		}
	//		lockErr := m.dMessageLocker.UnLockGlobalMessage(req.ClientMsgID)
	//		if lockErr != nil {
	//			log.Error(req.OperationID, "UnLockGlobalMessage err:", lockErr.Error())
	//		}
	//	}
	//
	//} else {
	//	log.Debug(req.OperationID, "redis handle secondly", req.String())
	//
	//	for k, v := range req.Pb2Model {
	//		err := m.dMessageLocker.LockMessageTypeKey(req.ClientMsgID, k)
	//		if err != nil {
	//			setKeyResultInfo(&resp, 100, err.Error(), req.ClientMsgID, k, v)
	//			continue
	//		}
	//		redisValue, err := db.DB.GetMessageTypeKeyValue(req.ClientMsgID, req.SessionType, k)
	//		if err != nil && err != go_redis.Nil {
	//			setKeyResultInfo(&resp, 200, err.Error(), req.ClientMsgID, k, v)
	//			continue
	//		}
	//		temp := new(sdkws.KeyValue)
	//		utils.JsonStringToStruct(redisValue, temp)
	//		if v.LatestUpdateTime != temp.LatestUpdateTime {
	//			setKeyResultInfo(&resp, 300, "message have update", req.ClientMsgID, k, temp)
	//			continue
	//		} else {
	//			v.LatestUpdateTime = utils.GetCurrentTimestampByMill()
	// 			newerr := db.DB.SetMessageTypeKeyValue(req.ClientMsgID, req.SessionType, k, utils.StructToJsonString(v))
	//			if newerr != nil {
	//				setKeyResultInfo(&resp, 201, newerr.Error(), req.ClientMsgID, k, temp)
	//				continue
	//			}
	//			setKeyResultInfo(&resp, 0, "", req.ClientMsgID, k, v)
	//		}
	//
	//	}
	//}
	//if !isExists {
	//	if !req.IsReact {
	// 		notification.ExtendMessageUpdatedNotification(req.OperationID, req.OpUserID, req.conversationID,
	// req.SessionType, req, &resp, true, true)
	//	} else {
	// 		notification.ExtendMessageUpdatedNotification(req.OperationID, req.OpUserID, req.conversationID,
	// req.SessionType, req, &resp, false, false)
	//	}
	//} else {
	// 	notification.ExtendMessageUpdatedNotification(req.OperationID, req.OpUserID, req.conversationID,
	// req.SessionType, req, &resp, false, true)
	//}
	//log.Debug(req.OperationID, utils.GetSelfFuncName(), "m return is:", resp.String())
	return resp, nil

}

func (m *msgServer) setKeyResultInfo(
	ctx context.Context,
	r *msg.SetMessageReactionExtensionsResp,
	errCode int32,
	errMsg, clientMsgID, typeKey string,
	keyValue *sdkws.KeyValue,
) {
	temp := new(msg.KeyValueResp)
	temp.KeyValue = keyValue
	temp.ErrCode = errCode
	temp.ErrMsg = errMsg
	r.Result = append(r.Result, temp)
	_ = m.MessageLocker.UnLockMessageTypeKey(ctx, clientMsgID, typeKey)
}

func (m *msgServer) setDeleteKeyResultInfo(
	ctx context.Context,
	r *msg.DeleteMessagesReactionExtensionsResp,
	errCode int32,
	errMsg, clientMsgID, typeKey string,
	keyValue *sdkws.KeyValue,
) {
	temp := new(msg.KeyValueResp)
	temp.KeyValue = keyValue
	temp.ErrCode = errCode
	temp.ErrMsg = errMsg
	r.Result = append(r.Result, temp)
	_ = m.MessageLocker.UnLockMessageTypeKey(ctx, clientMsgID, typeKey)
}

func (m *msgServer) GetMessagesReactionExtensions(
	ctx context.Context,
	req *msg.GetMessagesReactionExtensionsReq,
) (resp *msg.GetMessagesReactionExtensionsResp, err error) {
	//log.Debug(req.OperationID, utils.GetSelfFuncName(), "m args is:", req.String())
	//var rResp msg.GetMessageListReactionExtensionsResp
	//for _, messageValue := range req.MessageReactionKeyList {
	//	var oneMessage msg.SingleMessageExtensionResult
	//	oneMessage.ClientMsgID = messageValue.ClientMsgID
	//
	//	isExists, err := db.DB.JudgeMessageReactionExist(messageValue.ClientMsgID, req.SessionType)
	//	if err != nil {
	//		rResp.ErrCode = 100
	//		rResp.ErrMsg = err.Error()
	//		return &rResp, nil
	//	}
	//	if isExists {
	//		redisValue, err := db.DB.GetOneMessageAllReactionList(messageValue.ClientMsgID, req.SessionType)
	//		if err != nil {
	//			oneMessage.ErrCode = 100
	//			oneMessage.ErrMsg = err.Error()
	//			rResp.SingleMessageResult = append(rResp.SingleMessageResult, &oneMessage)
	//			continue
	//		}
	//		keyMap := make(map[string]*sdkws.KeyValue)
	//
	//		for k, v := range redisValue {
	//			temp := new(sdkws.KeyValue)
	//			utils.JsonStringToStruct(v, temp)
	//			keyMap[k] = temp
	//		}
	//		oneMessage.Pb2Model = keyMap
	//
	//	} else {
	// 		mongoValue, err := db.DB.GetExtendMsg(req.conversationID, req.SessionType, messageValue.ClientMsgID,
	// messageValue.MsgFirstModifyTime)
	//		if err != nil {
	//			oneMessage.ErrCode = 100
	//			oneMessage.ErrMsg = err.Error()
	//			rResp.SingleMessageResult = append(rResp.SingleMessageResult, &oneMessage)
	//			continue
	//		}
	//		keyMap := make(map[string]*sdkws.KeyValue)
	//
	//		for k, v := range mongoValue.Pb2Model {
	//			temp := new(sdkws.KeyValue)
	//			temp.TypeKey = v.TypeKey
	//			temp.Value = v.Value
	//			temp.LatestUpdateTime = v.LatestUpdateTime
	//			keyMap[k] = temp
	//		}
	//		oneMessage.Pb2Model = keyMap
	//	}
	//	rResp.SingleMessageResult = append(rResp.SingleMessageResult, &oneMessage)
	//}
	//log.Debug(req.OperationID, utils.GetSelfFuncName(), "m return is:", rResp.String())
	return resp, nil

}

func (m *msgServer) AddMessageReactionExtensions(
	ctx context.Context,
	req *msg.ModifyMessageReactionExtensionsReq,
) (resp *msg.ModifyMessageReactionExtensionsResp, err error) {
	return
}

func (m *msgServer) DeleteMessageReactionExtensions(
	ctx context.Context,
	req *msg.DeleteMessagesReactionExtensionsReq,
) (resp *msg.DeleteMessagesReactionExtensionsResp, err error) {
	//log.Debug(req.OperationID, utils.GetSelfFuncName(), "m args is:", req.String())
	//var rResp msg.DeleteMessagesReactionExtensionsResp
	//callbackResp := notification.callbackDeleteMessageReactionExtensions(req)
	//if callbackResp.ActionCode != constant.ActionAllow || callbackResp.ErrCode != 0 {
	//	rResp.ErrCode = int32(callbackResp.ErrCode)
	//	rResp.ErrMsg = callbackResp.ErrMsg
	//	for _, value := range req.Pb2Model {
	//		temp := new(msg.KeyValueResp)
	//		temp.KeyValue = value
	//		temp.ErrMsg = callbackResp.ErrMsg
	//		temp.ErrCode = 100
	//		rResp.Result = append(rResp.Result, temp)
	//	}
	//	return &rResp, nil
	//}
	////if ExternalExtension
	//if req.IsExternalExtensions {
	//	rResp.Result = callbackResp.ResultReactionExtensionList
	// 	notification.ExtendMessageDeleteNotification(req.OperationID, req.OpUserID, req.conversationID, req.SessionType,
	// req, &rResp, false, false)
	//	return &rResp, nil
	//
	//}
	//for _, v := range callbackResp.ResultReactionExtensions {
	//	if v.ErrCode != 0 {
	//		func(req *[]*sdkws.KeyValue, typeKey string) {
	//			for i := 0; i < len(*req); i++ {
	//				if (*req)[i].TypeKey == typeKey {
	//					*req = append((*req)[:i], (*req)[i+1:]...)
	//				}
	//			}
	//		}(&req.Pb2Model, v.KeyValue.TypeKey)
	//		rResp.Result = append(rResp.Result, v)
	//	}
	//}
	//isExists, err := db.DB.JudgeMessageReactionExist(req.ClientMsgID, req.SessionType)
	//if err != nil {
	//	rResp.ErrCode = 100
	//	rResp.ErrMsg = err.Error()
	//	for _, value := range req.Pb2Model {
	//		temp := new(msg.KeyValueResp)
	//		temp.KeyValue = value
	//		temp.ErrMsg = err.Error()
	//		temp.ErrCode = 100
	//		rResp.Result = append(rResp.Result, temp)
	//	}
	//	return &rResp, nil
	//}
	//
	//if isExists {
	//	log.Debug(req.OperationID, "redis handle this delete", req.String())
	//	for _, v := range req.Pb2Model {
	//		err := m.dMessageLocker.LockMessageTypeKey(req.ClientMsgID, v.TypeKey)
	//		if err != nil {
	//			setDeleteKeyResultInfo(&rResp, 100, err.Error(), req.ClientMsgID, v.TypeKey, v)
	//			continue
	//		}
	//
	//		redisValue, err := db.DB.GetMessageTypeKeyValue(req.ClientMsgID, req.SessionType, v.TypeKey)
	//		if err != nil && err != go_redis.Nil {
	//			setDeleteKeyResultInfo(&rResp, 200, err.Error(), req.ClientMsgID, v.TypeKey, v)
	//			continue
	//		}
	//		temp := new(sdkws.KeyValue)
	//		utils.JsonStringToStruct(redisValue, temp)
	//		if v.LatestUpdateTime != temp.LatestUpdateTime {
	//			setDeleteKeyResultInfo(&rResp, 300, "message have update", req.ClientMsgID, v.TypeKey, temp)
	//			continue
	//		} else {
	//			newErr := db.DB.DeleteOneMessageKey(req.ClientMsgID, req.SessionType, v.TypeKey)
	//			if newErr != nil {
	//				setDeleteKeyResultInfo(&rResp, 201, newErr.Error(), req.ClientMsgID, v.TypeKey, temp)
	//				continue
	//			}
	//			setDeleteKeyResultInfo(&rResp, 0, "", req.ClientMsgID, v.TypeKey, v)
	//		}
	//	}
	//} else {
	//	err := m.dMessageLocker.LockGlobalMessage(req.ClientMsgID)
	//	if err != nil {
	//		rResp.ErrCode = 100
	//		rResp.ErrMsg = err.Error()
	//		for _, value := range req.Pb2Model {
	//			temp := new(msg.KeyValueResp)
	//			temp.KeyValue = value
	//			temp.ErrMsg = err.Error()
	//			temp.ErrCode = 100
	//			rResp.Result = append(rResp.Result, temp)
	//		}
	//		return &rResp, nil
	//	}
	// 	mongoValue, err := db.DB.GetExtendMsg(req.conversationID, req.SessionType, req.ClientMsgID,
	// req.MsgFirstModifyTime)
	//	if err != nil {
	//		rResp.ErrCode = 200
	//		rResp.ErrMsg = err.Error()
	//		for _, value := range req.Pb2Model {
	//			temp := new(msg.KeyValueResp)
	//			temp.KeyValue = value
	//			temp.ErrMsg = err.Error()
	//			temp.ErrCode = 100
	//			rResp.Result = append(rResp.Result, temp)
	//		}
	//		return &rResp, nil
	//	}
	//	setValue := make(map[string]*sdkws.KeyValue)
	//	for _, v := range req.Pb2Model {
	//
	//		temp := new(sdkws.KeyValue)
	//		if vv, ok := mongoValue.Pb2Model[v.TypeKey]; ok {
	//			utils.CopyStructFields(temp, &vv)
	//			if v.LatestUpdateTime != vv.LatestUpdateTime {
	//				setDeleteKeyResultInfo(&rResp, 300, "message have update", req.ClientMsgID, v.TypeKey, temp)
	//				continue
	//			}
	//		} else {
	//			setDeleteKeyResultInfo(&rResp, 400, "key not in", req.ClientMsgID, v.TypeKey, v)
	//			continue
	//		}
	//		temp.TypeKey = v.TypeKey
	//		setValue[v.TypeKey] = temp
	//	}
	// 	err = db.DB.DeleteReactionExtendMsgSet(req.conversationID, req.SessionType, req.ClientMsgID,
	// req.MsgFirstModifyTime, setValue)
	//	if err != nil {
	//		for _, value := range setValue {
	//			temp := new(msg.KeyValueResp)
	//			temp.KeyValue = value
	//			temp.ErrMsg = err.Error()
	//			temp.ErrCode = 100
	//			rResp.Result = append(rResp.Result, temp)
	//		}
	//	} else {
	//		for _, value := range setValue {
	//			temp := new(msg.KeyValueResp)
	//			temp.KeyValue = value
	//			rResp.Result = append(rResp.Result, temp)
	//		}
	//	}
	//	lockErr := m.dMessageLocker.UnLockGlobalMessage(req.ClientMsgID)
	//	if lockErr != nil {
	//		log.Error(req.OperationID, "UnLockGlobalMessage err:", lockErr.Error())
	//	}
	//
	//}
	// notification.ExtendMessageDeleteNotification(req.OperationID, req.OpUserID, req.conversationID, req.SessionType,
	// req, &rResp, false, isExists)
	//log.Debug(req.OperationID, utils.GetSelfFuncName(), "m return is:", rResp.String())
	return resp, nil
}

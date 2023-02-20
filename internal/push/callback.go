package push

import (
	"Open_IM/pkg/callbackstruct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/tracelog"
	common "Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
)

func url() string {
	return config.Config.Callback.CallbackUrl
}

func callbackOfflinePush(ctx context.Context, userIDList []string, msg *common.MsgData, offlinePushUserIDList *[]string) error {
	if !config.Config.Callback.CallbackOfflinePush.Enable {
		return nil
	}
	req := &callbackstruct.CallbackBeforePushReq{
		UserStatusBatchCallbackReq: callbackstruct.UserStatusBatchCallbackReq{
			UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackOfflinePushCommand,
				OperationID:     tracelog.GetOperationID(ctx),
				PlatformID:      int(msg.SenderPlatformID),
				Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
			},
			UserIDList: userIDList,
		},
		OfflinePushInfo: msg.OfflinePushInfo,
		ClientMsgID:     msg.ClientMsgID,
		SendID:          msg.SendID,
		GroupID:         msg.GroupID,
		ContentType:     msg.ContentType,
		SessionType:     msg.SessionType,
		AtUserIDs:       msg.AtUserIDList,
		Content:         utils.GetContent(msg),
	}
	resp := &callbackstruct.CallbackBeforePushResp{}
	err := http.CallBackPostReturn(url(), req, resp, config.Config.Callback.CallbackOfflinePush)
	if err != nil {
		return err
	}
	if len(resp.UserIDList) != 0 {
		*offlinePushUserIDList = resp.UserIDList
	}
	if resp.OfflinePushInfo != nil {
		msg.OfflinePushInfo = resp.OfflinePushInfo
	}
	return nil
}

func callbackOnlinePush(operationID string, userIDList []string, msg *common.MsgData) error {
	if !config.Config.Callback.CallbackOnlinePush.Enable || utils.Contain(msg.SendID, userIDList...) {
		return nil
	}
	req := callbackstruct.CallbackBeforePushReq{
		UserStatusBatchCallbackReq: callbackstruct.UserStatusBatchCallbackReq{
			UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackOnlinePushCommand,
				OperationID:     operationID,
				PlatformID:      int(msg.SenderPlatformID),
				Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
			},
			UserIDList: userIDList,
		},
		ClientMsgID: msg.ClientMsgID,
		SendID:      msg.SendID,
		GroupID:     msg.GroupID,
		ContentType: msg.ContentType,
		SessionType: msg.SessionType,
		AtUserIDs:   msg.AtUserIDList,
		Content:     utils.GetContent(msg),
	}
	resp := &callbackstruct.CallbackBeforePushResp{}
	return http.CallBackPostReturn(url(), req, resp, config.Config.Callback.CallbackOnlinePush)
}

func callbackBeforeSuperGroupOnlinePush(ctx context.Context, groupID string, msg *common.MsgData, pushToUserList *[]string) error {
	if !config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.Enable {
		return nil
	}
	req := callbackstruct.CallbackBeforeSuperGroupOnlinePushReq{
		UserStatusBaseCallback: callbackstruct.UserStatusBaseCallback{
			CallbackCommand: constant.CallbackSuperGroupOnlinePushCommand,
			OperationID:     tracelog.GetOperationID(ctx),
			PlatformID:      int(msg.SenderPlatformID),
			Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
		},
		ClientMsgID:  msg.ClientMsgID,
		SendID:       msg.SendID,
		GroupID:      groupID,
		ContentType:  msg.ContentType,
		SessionType:  msg.SessionType,
		AtUserIDList: msg.AtUserIDList,
		Content:      utils.GetContent(msg),
		Seq:          msg.Seq,
	}
	resp := &callbackstruct.CallbackBeforeSuperGroupOnlinePushResp{}
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackBeforeSuperGroupOnlinePush); err != nil {
		return err
	}
	if len(resp.UserIDList) != 0 {
		*pushToUserList = resp.UserIDList
	}
	return nil
}

//func callbackOfflinePush(operationID string, userIDList []string, msg *common.MsgData, offlinePushUserIDList *[]string) cbApi.CommonCallbackResp {
//	callbackResp := cbapi.CommonCallbackResp{OperationID: operationID}
//	if !config.Config.Callback.CallbackOfflinePush.Enable {
//		return callbackResp
//	}
//	req := cbApi.CallbackBeforePushReq{
//		UserStatusBatchCallbackReq: cbApi.UserStatusBatchCallbackReq{
//			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
//				CallbackCommand: constant.CallbackOfflinePushCommand,
//				OperationID:     operationID,
//				PlatformID:      msg.SenderPlatformID,
//				Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
//			},
//			UserIDList: userIDList,
//		},
//		OfflinePushInfo: msg.OfflinePushInfo,
//		ClientMsgID:     msg.ClientMsgID,
//		SendID:          msg.SendID,
//		GroupID:         msg.GroupID,
//		ContentType:     msg.ContentType,
//		SessionType:     msg.SessionType,
//		AtUserIDList:    msg.AtUserIDList,
//		Content:         callback.GetContent(msg),
//	}
//	resp := &cbApi.CallbackBeforePushResp{CommonCallbackResp: &callbackResp}
//	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackOfflinePushCommand, req, resp, config.Config.Callback.CallbackOfflinePush.CallbackTimeOut); err != nil {
//		callbackResp.ErrCode = http2.StatusInternalServerError
//		callbackResp.ErrMsg = err.Error()
//		if !*config.Config.Callback.CallbackOfflinePush.CallbackFailedContinue {
//			callbackResp.ActionCode = constant.ActionForbidden
//			return callbackResp
//		} else {
//			callbackResp.ActionCode = constant.ActionAllow
//			return callbackResp
//		}
//	}
//	if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow {
//		if len(resp.UserIDList) != 0 {
//			*offlinePushUserIDList = resp.UserIDList
//		}
//		if resp.OfflinePushInfo != nil {
//			msg.OfflinePushInfo = resp.OfflinePushInfo
//		}
//	}
//	log.NewDebug(operationID, utils.GetSelfFuncName(), offlinePushUserIDList, resp.UserIDList)
//	return callbackResp
//}
//
//func callbackOnlinePush(operationID string, userIDList []string, msg *common.MsgData) cbApi.CommonCallbackResp {
//	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
//	if !config.Config.Callback.CallbackOnlinePush.Enable || utils.IsContain(msg.SendID, userIDList) {
//		return callbackResp
//	}
//	req := cbApi.CallbackBeforePushReq{
//		UserStatusBatchCallbackReq: cbApi.UserStatusBatchCallbackReq{
//			UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
//				CallbackCommand: constant.CallbackOnlinePushCommand,
//				OperationID:     operationID,
//				PlatformID:      msg.SenderPlatformID,
//				Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
//			},
//			UserIDList: userIDList,
//		},
//		//OfflinePushInfo: msg.OfflinePushInfo,
//		ClientMsgID:  msg.ClientMsgID,
//		SendID:       msg.SendID,
//		GroupID:      msg.GroupID,
//		ContentType:  msg.ContentType,
//		SessionType:  msg.SessionType,
//		AtUserIDList: msg.AtUserIDList,
//		Content:      callback.GetContent(msg),
//	}
//	resp := &cbApi.CallbackBeforePushResp{CommonCallbackResp: &callbackResp}
//	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackOnlinePushCommand, req, resp, config.Config.Callback.CallbackOnlinePush.CallbackTimeOut); err != nil {
//		callbackResp.ErrCode = http2.StatusInternalServerError
//		callbackResp.ErrMsg = err.Error()
//		if !config.Config.Callback.CallbackOnlinePush.CallbackFailedContinue {
//			callbackResp.ActionCode = constant.ActionForbidden
//			return callbackResp
//		} else {
//			callbackResp.ActionCode = constant.ActionAllow
//			return callbackResp
//		}
//	}
//	if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow {
//		//if resp.OfflinePushInfo != nil {
//		//	msg.OfflinePushInfo = resp.OfflinePushInfo
//		//}
//	}
//	return callbackResp
//}
//
//func callbackBeforeSuperGroupOnlinePush(operationID string, groupID string, msg *common.MsgData, pushToUserList *[]string) cbApi.CommonCallbackResp {
//	log.Debug(operationID, utils.GetSelfFuncName(), groupID, msg.String(), pushToUserList)
//	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
//	if !config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.Enable {
//		return callbackResp
//	}
//	req := cbApi.CallbackBeforeSuperGroupOnlinePushReq{
//		UserStatusBaseCallback: cbApi.UserStatusBaseCallback{
//			CallbackCommand: constant.CallbackSuperGroupOnlinePushCommand,
//			OperationID:     operationID,
//			PlatformID:      msg.SenderPlatformID,
//			Platform:        constant.PlatformIDToName(int(msg.SenderPlatformID)),
//		},
//		//OfflinePushInfo: msg.OfflinePushInfo,
//		ClientMsgID:  msg.ClientMsgID,
//		SendID:       msg.SendID,
//		GroupID:      groupID,
//		ContentType:  msg.ContentType,
//		SessionType:  msg.SessionType,
//		AtUserIDList: msg.AtUserIDList,
//		Content:      callback.GetContent(msg),
//		Seq:          msg.Seq,
//	}
//	resp := &cbApi.CallbackBeforeSuperGroupOnlinePushResp{CommonCallbackResp: &callbackResp}
//	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackSuperGroupOnlinePushCommand, req, resp, config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.CallbackTimeOut); err != nil {
//		callbackResp.ErrCode = http2.StatusInternalServerError
//		callbackResp.ErrMsg = err.Error()
//		if !config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.CallbackFailedContinue {
//			callbackResp.ActionCode = constant.ActionForbidden
//			return callbackResp
//		} else {
//			callbackResp.ActionCode = constant.ActionAllow
//			return callbackResp
//		}
//	}
//	if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow {
//		if len(resp.UserIDList) != 0 {
//			*pushToUserList = resp.UserIDList
//		}
//		//if resp.OfflinePushInfo != nil {
//		//	msg.OfflinePushInfo = resp.OfflinePushInfo
//		//}
//	}
//	log.NewDebug(operationID, utils.GetSelfFuncName(), pushToUserList, resp.UserIDList)
//	return callbackResp
//
//}

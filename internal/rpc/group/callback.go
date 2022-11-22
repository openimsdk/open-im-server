package group

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	http2 "net/http"
)

func callbackBeforeCreateGroup(req *pbGroup.CreateGroupReq) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: req.OperationID}
	if !config.Config.Callback.CallbackBeforeCreateGroup.Enable {
		return callbackResp
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), req.String())
	commonCallbackReq := &cbApi.CallbackBeforeCreateGroupReq{
		CallbackCommand: constant.CallbackBeforeCreateGroupCommand,
		GroupInfo:       *req.GroupInfo,
		InitMemberList:  req.InitMemberList,
	}
	resp := &cbApi.CallbackBeforeCreateGroupResp{
		CommonCallbackResp: &callbackResp,
	}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(req.OperationID, utils.GetSelfFuncName(), commonCallbackReq, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeCreateGroupCommand, commonCallbackReq, resp, config.Config.Callback.CallbackBeforeCreateGroup.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackBeforeCreateGroup.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow {
		if resp.GroupID != nil {
			req.GroupInfo.GroupID = *resp.GroupID
		}
		if resp.GroupName != nil {
			req.GroupInfo.GroupName = *resp.GroupName
		}
		if resp.Notification != nil {
			req.GroupInfo.Notification = *resp.Notification
		}
		if resp.Introduction != nil {
			req.GroupInfo.Introduction = *resp.Introduction
		}
		if resp.FaceURL != nil {
			req.GroupInfo.FaceURL = *resp.FaceURL
		}
		if resp.OwnerUserID != nil {
			req.GroupInfo.OwnerUserID = *resp.OwnerUserID
		}
		if resp.Ex != nil {
			req.GroupInfo.Ex = *resp.Ex
		}
		if resp.Status != nil {
			req.GroupInfo.Status = *resp.Status
		}
		if resp.CreatorUserID != nil {
			req.GroupInfo.CreatorUserID = *resp.CreatorUserID
		}
		if resp.GroupType != nil {
			req.GroupInfo.GroupType = *resp.GroupType
		}
		if resp.NeedVerification != nil {
			req.GroupInfo.NeedVerification = *resp.NeedVerification
		}
		if resp.LookMemberInfo != nil {
			req.GroupInfo.LookMemberInfo = *resp.LookMemberInfo
		}
	}
	return callbackResp
}

func CallbackBeforeMemberJoinGroup(operationID string, groupMember *db.GroupMember, groupEx string) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackBeforeMemberJoinGroup.Enable {
		return callbackResp
	}
	log.NewDebug(operationID, "args: ", *groupMember)
	callbackReq := cbApi.CallbackBeforeMemberJoinGroupReq{
		CallbackCommand: constant.CallbackBeforeMemberJoinGroupCommand,
		GroupID:         groupMember.GroupID,
		UserID:          groupMember.UserID,
		Ex:              groupMember.Ex,
		GroupEx:         groupEx,
	}
	resp := &cbApi.CallbackBeforeMemberJoinGroupResp{
		CommonCallbackResp: &callbackResp,
	}
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeMemberJoinGroupCommand, callbackReq, resp, config.Config.Callback.CallbackBeforeMemberJoinGroup.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackBeforeMemberJoinGroup.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	if resp.MuteEndTime != nil {
		groupMember.MuteEndTime = utils.UnixSecondToTime(*resp.MuteEndTime)
	}
	if resp.FaceURL != nil {
		groupMember.FaceURL = *resp.FaceURL
	}
	if resp.Ex != nil {
		groupMember.Ex = *resp.Ex
	}
	if resp.NickName != nil {
		groupMember.Nickname = *resp.NickName
	}
	if resp.RoleLevel != nil {
		groupMember.RoleLevel = *resp.RoleLevel
	}
	return callbackResp
}

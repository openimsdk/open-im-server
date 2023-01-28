package group

import (
	cbApi "Open_IM/pkg/callback_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	relation "Open_IM/pkg/common/db/mysql"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/trace_log"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func callbackBeforeCreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (err error) {
	defer func() {
		trace_log.SetCtxInfo(ctx, utils.GetFuncName(1), err, "req", req)
	}()
	if !config.Config.Callback.CallbackBeforeCreateGroup.Enable {
		return nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), req.String())
	commonCallbackReq := &cbApi.CallbackBeforeCreateGroupReq{
		CallbackCommand: constant.CallbackBeforeCreateGroupCommand,
		OperationID:     req.OperationID,
		GroupInfo:       *req.GroupInfo,
		InitMemberList:  req.InitMemberList,
	}
	callbackResp := cbApi.CommonCallbackResp{OperationID: req.OperationID}
	resp := &cbApi.CallbackBeforeCreateGroupResp{
		CommonCallbackResp: &callbackResp,
	}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(req.OperationID, utils.GetSelfFuncName(), commonCallbackReq, *resp)
	err = http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeCreateGroupCommand, commonCallbackReq,
		resp, config.Config.Callback.CallbackBeforeCreateGroup)
	if err == nil {
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
	return err
}

func CallbackBeforeMemberJoinGroup(ctx context.Context, operationID string, groupMember *relation.GroupMember, groupEx string) (err error) {
	defer func() {
		trace_log.SetCtxInfo(ctx, utils.GetFuncName(1), err, "groupMember", *groupMember, "groupEx", groupEx)
	}()
	callbackResp := cbApi.CommonCallbackResp{OperationID: operationID}
	if !config.Config.Callback.CallbackBeforeMemberJoinGroup.Enable {
		return nil
	}
	log.NewDebug(operationID, "args: ", *groupMember)
	callbackReq := cbApi.CallbackBeforeMemberJoinGroupReq{
		CallbackCommand: constant.CallbackBeforeMemberJoinGroupCommand,
		OperationID:     operationID,
		GroupID:         groupMember.GroupID,
		UserID:          groupMember.UserID,
		Ex:              groupMember.Ex,
		GroupEx:         groupEx,
	}
	resp := &cbApi.CallbackBeforeMemberJoinGroupResp{
		CommonCallbackResp: &callbackResp,
	}
	err = http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeMemberJoinGroupCommand, callbackReq,
		resp, config.Config.Callback.CallbackBeforeMemberJoinGroup)
	if err == nil {
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
	}
	return err
}

func CallbackBeforeSetGroupMemberInfo(ctx context.Context, req *pbGroup.SetGroupMemberInfoReq) (err error) {
	defer func() {
		trace_log.SetCtxInfo(ctx, utils.GetFuncName(1), err, "req", *req)
	}()
	callbackResp := cbApi.CommonCallbackResp{OperationID: req.OperationID}
	if !config.Config.Callback.CallbackBeforeSetGroupMemberInfo.Enable {
		return nil
	}
	callbackReq := cbApi.CallbackBeforeSetGroupMemberInfoReq{
		CallbackCommand: constant.CallbackBeforeSetGroupMemberInfoCommand,
		OperationID:     req.OperationID,
		GroupID:         req.GroupID,
		UserID:          req.UserID,
	}
	if req.Nickname != nil {
		callbackReq.Nickname = req.Nickname.Value
	}
	if req.FaceURL != nil {
		callbackReq.FaceURL = req.FaceURL.Value
	}
	if req.RoleLevel != nil {
		callbackReq.RoleLevel = req.RoleLevel.Value
	}
	if req.Ex != nil {
		callbackReq.Ex = req.Ex.Value
	}
	resp := &cbApi.CallbackBeforeSetGroupMemberInfoResp{
		CommonCallbackResp: &callbackResp,
	}
	err = http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeSetGroupMemberInfoCommand, callbackReq,
		resp, config.Config.Callback.CallbackBeforeSetGroupMemberInfo.CallbackTimeOut, &config.Config.Callback.CallbackBeforeSetGroupMemberInfo.CallbackFailedContinue)
	if err == nil {
		if resp.FaceURL != nil {
			req.FaceURL = &wrapperspb.StringValue{Value: *resp.FaceURL}
		}
		if resp.Nickname != nil {
			req.Nickname = &wrapperspb.StringValue{Value: *resp.Nickname}
		}
		if resp.RoleLevel != nil {
			req.RoleLevel = &wrapperspb.Int32Value{Value: *resp.RoleLevel}
		}
		if resp.Ex != nil {
			req.Ex = &wrapperspb.StringValue{Value: *resp.Ex}
		}
	}
	return err
}

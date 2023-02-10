package group

import (
	"Open_IM/pkg/apistruct"
	"Open_IM/pkg/callbackstruct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func CallbackBeforeCreateGroup(ctx context.Context, req *group.CreateGroupReq) (err error) {
	if !config.Config.Callback.CallbackBeforeCreateGroup.Enable {
		return nil
	}
	defer func() {
		tracelog.SetCtxInfo(ctx, utils.GetFuncName(1), err, "req", req)
	}()
	operationID := tracelog.GetOperationID(ctx)
	commonCallbackReq := &callbackstruct.CallbackBeforeCreateGroupReq{
		CallbackCommand: constant.CallbackBeforeCreateGroupCommand,
		OperationID:     operationID,
		GroupInfo:       *req.GroupInfo,
	}
	commonCallbackReq.InitMemberList = append(commonCallbackReq.InitMemberList, &apistruct.GroupAddMemberInfo{
		UserID:    req.OwnerUserID,
		RoleLevel: constant.GroupOwner,
	})
	for _, userID := range req.AdminUserIDs {
		commonCallbackReq.InitMemberList = append(commonCallbackReq.InitMemberList, &apistruct.GroupAddMemberInfo{
			UserID:    userID,
			RoleLevel: constant.GroupAdmin,
		})
	}
	for _, userID := range req.AdminUserIDs {
		commonCallbackReq.InitMemberList = append(commonCallbackReq.InitMemberList, &apistruct.GroupAddMemberInfo{
			UserID:    userID,
			RoleLevel: constant.GroupOrdinaryUsers,
		})
	}
	resp := &callbackstruct.CallbackBeforeCreateGroupResp{
		CommonCallbackResp: &callbackstruct.CommonCallbackResp{OperationID: operationID},
	}
	err = http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeCreateGroupCommand, commonCallbackReq, resp, config.Config.Callback.CallbackBeforeCreateGroup)
	if err != nil {
		return err
	}

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
	return nil
}

func CallbackBeforeMemberJoinGroup(ctx context.Context, groupMember *relation.GroupMemberModel, groupEx string) (err error) {
	if !config.Config.Callback.CallbackBeforeMemberJoinGroup.Enable {
		return nil
	}
	defer func() {
		tracelog.SetCtxInfo(ctx, utils.GetFuncName(1), err, "groupMember", *groupMember, "groupEx", groupEx)
	}()
	operationID := tracelog.GetOperationID(ctx)
	callbackResp := callbackstruct.CommonCallbackResp{OperationID: operationID}
	callbackReq := callbackstruct.CallbackBeforeMemberJoinGroupReq{
		CallbackCommand: constant.CallbackBeforeMemberJoinGroupCommand,
		OperationID:     operationID,
		GroupID:         groupMember.GroupID,
		UserID:          groupMember.UserID,
		Ex:              groupMember.Ex,
		GroupEx:         groupEx,
	}
	resp := &callbackstruct.CallbackBeforeMemberJoinGroupResp{
		CommonCallbackResp: &callbackResp,
	}
	err = http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeMemberJoinGroupCommand, callbackReq,
		resp, config.Config.Callback.CallbackBeforeMemberJoinGroup)
	if err != nil {
		return err
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
	return nil
}

func CallbackBeforeSetGroupMemberInfo(ctx context.Context, req *group.SetGroupMemberInfo) (err error) {
	if !config.Config.Callback.CallbackBeforeSetGroupMemberInfo.Enable {
		return nil
	}
	defer func() {
		tracelog.SetCtxInfo(ctx, utils.GetFuncName(1), err, "req", *req)
	}()
	operationID := tracelog.GetOperationID(ctx)
	callbackResp := callbackstruct.CommonCallbackResp{OperationID: operationID}
	callbackReq := callbackstruct.CallbackBeforeSetGroupMemberInfoReq{
		CallbackCommand: constant.CallbackBeforeSetGroupMemberInfoCommand,
		OperationID:     operationID,
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
	resp := &callbackstruct.CallbackBeforeSetGroupMemberInfoResp{
		CommonCallbackResp: &callbackResp,
	}
	err = http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeSetGroupMemberInfoCommand, callbackReq, resp, config.Config.Callback.CallbackBeforeSetGroupMemberInfo)
	if err != nil {
		return err
	}
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
	return err
}

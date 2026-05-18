// Copyright © 2024 OpenIM. All rights reserved.
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

package rtc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/livekit/protocol/auth"
	livekit "github.com/livekit/protocol/livekit"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/rtc"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/proto"
)

// SignalMessageAssemble processes a signal request from the WebSocket gateway
// and assembles the appropriate signal response, sending notifications to peers.
func (s *rtcServer) SignalMessageAssemble(ctx context.Context, req *rtc.SignalMessageAssembleReq) (*rtc.SignalMessageAssembleResp, error) {
	if req.SignalReq == nil {
		return nil, errs.ErrArgs.WrapMsg("signalReq is nil")
	}
	var (
		resp    rtc.SignalResp
		respErr error
	)
	switch payload := req.SignalReq.Payload.(type) {
	case *rtc.SignalReq_Invite:
		log.ZInfo(ctx, "SignalMessageAssemble", "payload", payload.Invite)
		r, err := s.handleInvite(ctx, payload.Invite, req.SignalReq)
		resp.Payload = &rtc.SignalResp_Invite{Invite: r}
		respErr = err
	case *rtc.SignalReq_InviteInGroup:
		r, err := s.handleInviteInGroup(ctx, payload.InviteInGroup, req.SignalReq)
		resp.Payload = &rtc.SignalResp_InviteInGroup{InviteInGroup: r}
		respErr = err
	case *rtc.SignalReq_Cancel:
		r, err := s.handleCancel(ctx, payload.Cancel, req.SignalReq)
		resp.Payload = &rtc.SignalResp_Cancel{Cancel: r}
		respErr = err
	case *rtc.SignalReq_Accept:
		r, err := s.handleAccept(ctx, payload.Accept, req.SignalReq)
		resp.Payload = &rtc.SignalResp_Accept{Accept: r}
		respErr = err
	case *rtc.SignalReq_HungUp:
		r, err := s.handleHungUp(ctx, payload.HungUp, req.SignalReq)
		resp.Payload = &rtc.SignalResp_HungUp{HungUp: r}
		respErr = err
	case *rtc.SignalReq_Reject:
		r, err := s.handleReject(ctx, payload.Reject, req.SignalReq)
		resp.Payload = &rtc.SignalResp_Reject{Reject: r}
		respErr = err
	case *rtc.SignalReq_GetTokenByRoomID:
		r, err := s.handleGetTokenByRoomID(ctx, payload.GetTokenByRoomID)
		resp.Payload = &rtc.SignalResp_GetTokenByRoomID{GetTokenByRoomID: r}
		respErr = err
	default:
		return nil, errs.ErrArgs.WrapMsg("unknown signal payload type")
	}
	if respErr != nil {
		log.ZError(ctx, "SignalMessageAssemble", respErr, "err", respErr.Error())
		return nil, respErr
	}
	return &rtc.SignalMessageAssembleResp{SignalResp: &resp}, nil
}

// handleInvite processes a 1-to-1 call invitation.
func (s *rtcServer) handleInvite(ctx context.Context, req *rtc.SignalInviteReq, signalReq *rtc.SignalReq) (*rtc.SignalInviteResp, error) {
	inv := req.Invitation
	if inv == nil {
		log.ZError(ctx, "handleInvite", errs.ErrArgs, "r", "invitation is nil")
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}
	inv.RoomID = newRoomID()
	inv.InviterUserID = req.UserID
	inv.InitiateTime = time.Now().UnixMilli()

	if len(inv.InviteeUserIDList) == 0 {
		return nil, errs.ErrArgs.WrapMsg("no invitees", "inviteeUserIDList", inv.InviteeUserIDList)
	}

	for _, inviteeID := range inv.InviteeUserIDList {
		allowed, err := s.isCallAllowed(ctx, req.UserID, inviteeID)
		if err != nil {
			log.ZError(ctx, "handleInvite: isCallAllowed failed", err, "inviteeID", inviteeID)
			return nil, err
		}
		if !allowed {
			return nil, errs.ErrNoPermission.WrapMsg("the invitee does not accept calls from you", "inviteeID", inviteeID)
		}
	}

	// 检测哪些被叫用户正忙（已在通话中），记录到 BusyLineUserIDList
	busyUserIDs, err := s.db.GetBusyUserIDs(ctx, inv.InviteeUserIDList)
	if err != nil {
		log.ZWarn(ctx, "handleInvite: GetBusyUserIDs failed (non-fatal)", err)
	}
	busySet := make(map[string]struct{}, len(busyUserIDs))
	for _, uid := range busyUserIDs {
		busySet[uid] = struct{}{}
	}
	inv.BusyLineUserIDList = busyUserIDs

	if len(inv.InviteeUserIDList) == len(busyUserIDs) {
		return nil, errs.ErrNoPermission.WrapMsg("all invitees are busy", "inviteeUserIDList", inv.InviteeUserIDList)
	}

	// 从主叫用户资料获取铃声 URL，注入到邀请信息中，被叫方收到后播放主叫方铃声
	if inviterInfo, err := s.userClient.GetUserInfo(ctx, req.UserID); err == nil && inviterInfo.CallRingtoneURL != "" {
		inv.CallerRingtoneURL = inviterInfo.CallRingtoneURL
	}

	// 查询被叫方铃声 URL，供主叫方在等待时播放
	var calleeRingtoneURL string
	if len(inv.InviteeUserIDList) > 0 {
		if inviteeInfo, err := s.userClient.GetUserInfo(ctx, inv.InviteeUserIDList[0]); err == nil {
			calleeRingtoneURL = inviteeInfo.CallRingtoneURL
		}
	}

	if _, err := s.roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{Name: inv.RoomID}); err != nil {
		log.ZError(ctx, "handleInvite", err, "r", err.Error())
		return nil, errs.WrapMsg(err, "LiveKit CreateRoom failed", "roomID", inv.RoomID)
	}

	token, err := s.genToken(inv.RoomID, req.UserID)
	if err != nil {
		if _, delErr := s.roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{Room: inv.RoomID}); delErr != nil {
			log.ZWarn(ctx, "handleInvite: rollback DeleteRoom failed", delErr, "roomID", inv.RoomID)
		}
		return nil, err
	}

	if err := s.db.CreateInvitation(ctx, invitationToModel(inv, req.OfflinePushInfo)); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.ZWarn(ctx, "handleInvite: duplicate invitation (idempotent retry)", err, "roomID", inv.RoomID)
		} else {
			if _, delErr := s.roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{Room: inv.RoomID}); delErr != nil {
				log.ZWarn(ctx, "handleInvite: rollback DeleteRoom failed", delErr, "roomID", inv.RoomID)
			}
			return nil, errs.WrapMsg(err, "CreateInvitation failed", "roomID", inv.RoomID)
		}
	}

	content, err := marshalSignalReq(signalReq)
	if err != nil {
		return nil, err
	}

	for _, inviteeID := range inv.InviteeUserIDList {
		if _, busy := busySet[inviteeID]; busy {
			log.ZInfo(ctx, "handleInvite: skip busy invitee", "inviteeID", inviteeID)
			continue
		}
		log.ZInfo(ctx, "sendSignalingNotification to invitee", "sendID", req.UserID, "recvID", inviteeID)
		if err := s.sendSignalingNotification(ctx, req.UserID, inviteeID, int32(constant.SingleChatType), "", req.OfflinePushInfo, content); err != nil {
			log.ZError(ctx, "sendSignalingNotification to invitee failed", err, "inviteeID", inviteeID)
			return nil, errs.WrapMsg(err, "failed to notify invitee", "inviteeID", inviteeID)
		}
	}

	log.ZDebug(ctx, "handleInvite", "token", token, "roomID", inv.RoomID, "liveURL", s.config.RpcConfig.LiveKit.ExternalAddress)
	return &rtc.SignalInviteResp{
		Token:              token,
		RoomID:             inv.RoomID,
		LiveURL:            s.config.RpcConfig.LiveKit.ExternalAddress,
		BusyLineUserIDList: busyUserIDs,
		CalleeRingtoneURL:  calleeRingtoneURL,
	}, nil
}

// handleInviteInGroup processes a group call invitation.
func (s *rtcServer) handleInviteInGroup(ctx context.Context, req *rtc.SignalInviteInGroupReq, signalReq *rtc.SignalReq) (*rtc.SignalInviteInGroupResp, error) {
	inv := req.Invitation
	if inv == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}
	if inv.GroupID == "" {
		return nil, errs.ErrArgs.WrapMsg("groupID is empty")
	}

	inv.RoomID = newRoomID()
	inv.InviterUserID = req.UserID
	inv.InitiateTime = time.Now().UnixMilli()

	// 检测哪些被叫用户正忙（已在通话中），记录到 BusyLineUserIDList
	busyUserIDs, err := s.db.GetBusyUserIDs(ctx, inv.InviteeUserIDList)
	if err != nil {
		log.ZWarn(ctx, "handleInviteInGroup: GetBusyUserIDs failed (non-fatal)", err)
	}
	busySet := make(map[string]struct{}, len(busyUserIDs))
	for _, uid := range busyUserIDs {
		busySet[uid] = struct{}{}
	}
	inv.BusyLineUserIDList = busyUserIDs

	// 从主叫用户资料获取铃声 URL，注入到邀请信息中，被叫方收到后播放主叫方铃声
	if inviterInfo, err := s.userClient.GetUserInfo(ctx, req.UserID); err == nil && inviterInfo.CallRingtoneURL != "" {
		inv.CallerRingtoneURL = inviterInfo.CallRingtoneURL
	}

	// 查询第一位被叫的铃声 URL，供主叫方在等待时播放
	var calleeRingtoneURL string
	if len(inv.InviteeUserIDList) > 0 {
		if inviteeInfo, err := s.userClient.GetUserInfo(ctx, inv.InviteeUserIDList[0]); err == nil {
			calleeRingtoneURL = inviteeInfo.CallRingtoneURL
		}
	}

	if _, err := s.roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{Name: inv.RoomID}); err != nil {
		return nil, errs.WrapMsg(err, "LiveKit CreateRoom failed", "roomID", inv.RoomID)
	}

	token, err := s.genToken(inv.RoomID, req.UserID)
	if err != nil {
		if _, delErr := s.roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{Room: inv.RoomID}); delErr != nil {
			log.ZWarn(ctx, "handleInviteInGroup: rollback DeleteRoom failed", delErr, "roomID", inv.RoomID)
		}
		return nil, err
	}

	if err := s.db.CreateInvitation(ctx, invitationToModel(inv, req.OfflinePushInfo)); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			if _, delErr := s.roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{Room: inv.RoomID}); delErr != nil {
				log.ZWarn(ctx, "handleInviteInGroup: rollback DeleteRoom failed", delErr, "roomID", inv.RoomID)
			}
			return nil, errs.WrapMsg(err, "CreateInvitation failed", "roomID", inv.RoomID)
		}
		log.ZWarn(ctx, "handleInviteInGroup: duplicate invitation (idempotent retry)", err, "roomID", inv.RoomID)
	}

	content, err := marshalSignalReq(signalReq)
	if err != nil {
		return nil, err
	}
	for _, inviteeID := range inv.InviteeUserIDList {
		allowed, err := s.isCallAllowed(ctx, req.UserID, inviteeID)
		if err != nil {
			log.ZWarn(ctx, "handleInviteInGroup: isCallAllowed failed, skipping invitee", err, "inviteeID", inviteeID)
			continue
		}
		if !allowed {
			log.ZInfo(ctx, "handleInviteInGroup: skipping invitee (call setting blocked)", "inviteeID", inviteeID)
			continue
		}
		if _, busy := busySet[inviteeID]; busy {
			log.ZInfo(ctx, "handleInviteInGroup: skip busy invitee", "inviteeID", inviteeID)
			continue
		}
		if err := s.sendSignalingNotification(ctx, req.UserID, inviteeID, int32(constant.ReadGroupChatType), inv.GroupID, req.OfflinePushInfo, content); err != nil {
			log.ZWarn(ctx, "sendSignalingNotification to group invitee failed", err, "inviteeID", inviteeID)
		}
	}

	resp := &rtc.SignalInviteInGroupResp{
		Token:              token,
		RoomID:             inv.RoomID,
		LiveURL:            s.config.RpcConfig.LiveKit.ExternalAddress,
		BusyLineUserIDList: busyUserIDs,
		CalleeRingtoneURL:  calleeRingtoneURL,
	}

	log.ZDebug(ctx, "handleInviteInGroup", "req", req, "resp", resp)

	return resp, nil
}

// isCallAllowed 判断 inviterID 是否被允许向 inviteeID 发起音视频通话。
// 规则：
//   - CallAcceptSettingPublic(0)  → 所有人均可
//   - CallAcceptSettingFriends(1) → 仅当 inviterID 在 inviteeID 好友列表中
//   - CallAcceptSettingNobody(2)  → 任何人均不可
func (s *rtcServer) isCallAllowed(ctx context.Context, inviterID, inviteeID string) (bool, error) {
	userInfo, err := s.userClient.GetUserInfo(ctx, inviteeID)
	if err != nil {
		return false, err
	}
	switch userInfo.CallAcceptSetting {
	case model.CallAcceptSettingNobody:
		return false, nil
	case model.CallAcceptSettingFriends:
		isFriend, err := s.relationClient.IsFriend(ctx, inviteeID, inviterID)
		if err != nil {
			return false, err
		}
		return isFriend, nil
	default: // CallAcceptSettingPublic
		return true, nil
	}
}

func (s *rtcServer) handleAccept(ctx context.Context, req *rtc.SignalAcceptReq, signalReq *rtc.SignalReq) (*rtc.SignalAcceptResp, error) {
	if req.Invitation == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}

	// 从 DB 获取权威邀请数据，验证邀请存在且 userID 在被邀请人列表中
	dbInv, err := s.db.GetInvitationByRoomID(ctx, req.Invitation.RoomID)
	if err != nil {
		return nil, errs.WrapMsg(err, "invitation not found or expired", "roomID", req.Invitation.RoomID)
	}
	if !datautil.Contain(req.UserID, dbInv.InviteeUserIDList...) {
		return nil, errs.ErrNoPermission.WrapMsg("user not in invitee list", "userID", req.UserID)
	}

	token, err := s.genToken(dbInv.RoomID, req.UserID)
	if err != nil {
		return nil, err
	}

	sessionType := int32(constant.SingleChatType)
	if dbInv.GroupID != "" {
		sessionType = int32(constant.ReadGroupChatType)
	}

	content, err := marshalSignalReq(signalReq)
	if err != nil {
		return nil, err
	}

	if err := s.sendSignalingNotification(ctx, req.UserID, dbInv.InviterUserID, sessionType, dbInv.GroupID, req.OfflinePushInfo, content); err != nil {
		log.ZWarn(ctx, "sendSignalingNotification accept to inviter failed", err, "inviterID", dbInv.InviterUserID)
	}

	// 接受邀请后不删除 invitation：通话仍在进行，双方应被标记为忙线（BusyLineUserIDList）。
	// invitation 的清理由以下路径负责：
	//   - 主动挂断：handleHungUp → DeleteInvitation
	//   - 主叫取消：handleCancel → DeleteInvitation
	//   - 被叫拒绝：handleReject → DeleteInvitation / RemoveInvitee
	//   - 异常中断：MongoDB TTL 索引（expire_at 字段）自动清理

	return &rtc.SignalAcceptResp{
		Token:   token,
		RoomID:  dbInv.RoomID,
		LiveURL: s.config.RpcConfig.LiveKit.ExternalAddress,
	}, nil
}

// handleReject processes a call rejection.
func (s *rtcServer) handleReject(ctx context.Context, req *rtc.SignalRejectReq, signalReq *rtc.SignalReq) (*rtc.SignalRejectResp, error) {
	if req.Invitation == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}

	dbInv, err := s.db.GetInvitationByRoomID(ctx, req.Invitation.RoomID)
	if err != nil {
		return nil, errs.WrapMsg(err, "invitation not found or expired", "roomID", req.Invitation.RoomID)
	}
	if !datautil.Contain(req.UserID, dbInv.InviteeUserIDList...) {
		return nil, errs.ErrNoPermission.WrapMsg("user not in invitee list", "userID", req.UserID)
	}

	sessionType := int32(constant.SingleChatType)
	if dbInv.GroupID != "" {
		sessionType = int32(constant.ReadGroupChatType)
	}
	content, err := marshalSignalReq(signalReq)
	if err != nil {
		return nil, err
	}
	if err := s.sendSignalingNotification(ctx, req.UserID, dbInv.InviterUserID, sessionType, dbInv.GroupID, req.OfflinePushInfo, content); err != nil {
		log.ZWarn(ctx, "sendSignalingNotification reject to inviter failed", err, "inviterID", dbInv.InviterUserID)
	}

	if dbInv.GroupID != "" {
		if err := s.db.RemoveInvitee(ctx, dbInv.RoomID, req.UserID); err != nil {
			log.ZWarn(ctx, "RemoveInvitee failed", err, "roomID", dbInv.RoomID, "userID", req.UserID)
		}
	} else {
		if err := s.db.DeleteInvitation(ctx, dbInv.RoomID); err != nil {
			log.ZWarn(ctx, "DeleteInvitation failed", err, "roomID", dbInv.RoomID)
		}
	}

	return &rtc.SignalRejectResp{}, nil
}

// handleCancel processes a call cancellation.
func (s *rtcServer) handleCancel(ctx context.Context, req *rtc.SignalCancelReq, signalReq *rtc.SignalReq) (*rtc.SignalCancelResp, error) {
	if req.Invitation == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}

	dbInv, err := s.db.GetInvitationByRoomID(ctx, req.Invitation.RoomID)
	if err != nil {
		return nil, errs.WrapMsg(err, "invitation not found or expired", "roomID", req.Invitation.RoomID)
	}
	if req.UserID != dbInv.InviterUserID {
		return nil, errs.ErrNoPermission.WrapMsg("only the inviter can cancel", "userID", req.UserID, "inviterUserID", dbInv.InviterUserID)
	}

	sessionType := int32(constant.SingleChatType)
	if dbInv.GroupID != "" {
		sessionType = int32(constant.ReadGroupChatType)
	}
	content, err := marshalSignalReq(signalReq)
	if err != nil {
		return nil, err
	}
	for _, inviteeID := range dbInv.InviteeUserIDList {
		if err := s.sendSignalingNotification(ctx, req.UserID, inviteeID, sessionType, dbInv.GroupID, req.OfflinePushInfo, content); err != nil {
			log.ZWarn(ctx, "sendSignalingNotification cancel to invitee failed", err, "inviteeID", inviteeID)
		}
	}

	if err := s.db.DeleteInvitation(ctx, dbInv.RoomID); err != nil {
		log.ZWarn(ctx, "DeleteInvitation failed", err, "roomID", dbInv.RoomID)
	}

	return &rtc.SignalCancelResp{}, nil
}

// handleHungUp processes a call hang-up.
func (s *rtcServer) handleHungUp(ctx context.Context, req *rtc.SignalHungUpReq, signalReq *rtc.SignalReq) (*rtc.SignalHungUpResp, error) {
	if req.Invitation == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}

	dbInv, err := s.db.GetInvitationByRoomID(ctx, req.Invitation.RoomID)
	if err != nil {
		return nil, errs.WrapMsg(err, "invitation not found or expired", "roomID", req.Invitation.RoomID)
	}
	if req.UserID != dbInv.InviterUserID && !datautil.Contain(req.UserID, dbInv.InviteeUserIDList...) {
		return nil, errs.ErrNoPermission.WrapMsg("user is not a participant of this call", "userID", req.UserID)
	}

	sessionType := int32(constant.SingleChatType)
	if dbInv.GroupID != "" {
		sessionType = int32(constant.ReadGroupChatType)
	}
	content, err := marshalSignalReq(signalReq)
	if err != nil {
		return nil, err
	}
	// 使用 DB 中的参与者列表，不信任客户端传入的 InviteeUserIDList
	for _, peerID := range hungUpPeerIDsFromDB(dbInv, req.UserID) {
		if err := s.sendSignalingNotification(ctx, req.UserID, peerID, sessionType, dbInv.GroupID, req.OfflinePushInfo, content); err != nil {
			log.ZWarn(ctx, "sendSignalingNotification hungUp to peer failed", err, "peerID", peerID)
		}
	}

	// Terminate the LiveKit room
	if _, err := s.roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{Room: dbInv.RoomID}); err != nil {
		log.ZWarn(ctx, "LiveKit DeleteRoom failed", err, "roomID", dbInv.RoomID)
	}

	if err := s.db.DeleteInvitation(ctx, dbInv.RoomID); err != nil {
		log.ZWarn(ctx, "DeleteInvitation failed", err, "roomID", dbInv.RoomID)
	}

	return &rtc.SignalHungUpResp{}, nil
}

// handleGetTokenByRoomID returns a LiveKit token for an existing room.
func (s *rtcServer) handleGetTokenByRoomID(ctx context.Context, req *rtc.SignalGetTokenByRoomIDReq) (*rtc.SignalGetTokenByRoomIDResp, error) {
	dbInv, err := s.db.GetInvitationByRoomID(ctx, req.RoomID)
	if err != nil {
		return nil, errs.WrapMsg(err, "room not found or expired", "roomID", req.RoomID)
	}
	if req.UserID != dbInv.InviterUserID && !datautil.Contain(req.UserID, dbInv.InviteeUserIDList...) {
		return nil, errs.ErrNoPermission.WrapMsg("user is not a participant of this room", "userID", req.UserID)
	}

	token, err := s.genToken(req.RoomID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &rtc.SignalGetTokenByRoomIDResp{
		Token:   token,
		LiveURL: s.config.RpcConfig.LiveKit.ExternalAddress,
	}, nil
}

// SignalGetRoomByGroupID returns room information for a group.
func (s *rtcServer) SignalGetRoomByGroupID(ctx context.Context, req *rtc.SignalGetRoomByGroupIDReq) (*rtc.SignalGetRoomByGroupIDResp, error) {
	inv, err := s.db.GetInvitationByGroupID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	return &rtc.SignalGetRoomByGroupIDResp{
		Invitation: modelToInvitationInfo(inv),
		RoomID:     inv.RoomID,
	}, nil
}

// SignalGetTokenByRoomID returns a token for joining a room directly (HTTP API path).
// Fix P0(安全): 同 handleGetTokenByRoomID，添加参与者身份校验。
func (s *rtcServer) SignalGetTokenByRoomID(ctx context.Context, req *rtc.SignalGetTokenByRoomIDReq) (*rtc.SignalGetTokenByRoomIDResp, error) {
	dbInv, err := s.db.GetInvitationByRoomID(ctx, req.RoomID)
	if err != nil {
		return nil, errs.WrapMsg(err, "room not found or expired", "roomID", req.RoomID)
	}
	if req.UserID != dbInv.InviterUserID && !datautil.Contain(req.UserID, dbInv.InviteeUserIDList...) {
		return nil, errs.ErrNoPermission.WrapMsg("user is not a participant of this room", "userID", req.UserID)
	}

	token, err := s.genToken(req.RoomID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &rtc.SignalGetTokenByRoomIDResp{
		Token:   token,
		LiveURL: s.config.RpcConfig.LiveKit.ExternalAddress,
	}, nil
}

// SignalGetRooms returns room info for a list of room IDs.
func (s *rtcServer) SignalGetRooms(ctx context.Context, req *rtc.SignalGetRoomsReq) (*rtc.SignalGetRoomsResp, error) {
	if len(req.RoomIDs) == 0 {
		return &rtc.SignalGetRoomsResp{}, nil
	}
	invs, err := s.db.GetInvitationsByRoomIDs(ctx, req.RoomIDs)
	if err != nil {
		return nil, err
	}
	roomList := make([]*rtc.SignalGetRoomByGroupIDResp, 0, len(invs))
	for _, inv := range invs {
		roomList = append(roomList, &rtc.SignalGetRoomByGroupIDResp{
			Invitation: modelToInvitationInfo(inv),
			RoomID:     inv.RoomID,
		})
	}
	return &rtc.SignalGetRoomsResp{RoomList: roomList}, nil
}

// GetSignalInvitationInfo retrieves a pending invitation by roomID.
func (s *rtcServer) GetSignalInvitationInfo(ctx context.Context, req *rtc.GetSignalInvitationInfoReq) (*rtc.GetSignalInvitationInfoResp, error) {
	inv, err := s.db.GetInvitationByRoomID(ctx, req.RoomID)
	if err != nil {
		return nil, err
	}
	return &rtc.GetSignalInvitationInfoResp{
		InvitationInfo: modelToInvitationInfo(inv),
		OfflinePushInfo: &sdkws.OfflinePushInfo{
			Title: inv.OfflinePushTitle,
			Desc:  inv.OfflinePushDesc,
			Ex:    inv.OfflinePushEx,
		},
	}, nil
}

// GetSignalInvitationInfoStartApp retrieves a pending invitation for a user when the app starts.
func (s *rtcServer) GetSignalInvitationInfoStartApp(ctx context.Context, req *rtc.GetSignalInvitationInfoStartAppReq) (*rtc.GetSignalInvitationInfoStartAppResp, error) {
	inv, err := s.db.GetInvitationByInviteeUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &rtc.GetSignalInvitationInfoStartAppResp{
		Invitation: modelToInvitationInfo(inv),
		OfflinePushInfo: &sdkws.OfflinePushInfo{
			Title: inv.OfflinePushTitle,
			Desc:  inv.OfflinePushDesc,
			Ex:    inv.OfflinePushEx,
		},
	}, nil
}

// SignalSendCustomSignal forwards a custom signal to all participants in a room.
func (s *rtcServer) SignalSendCustomSignal(ctx context.Context, req *rtc.SignalSendCustomSignalReq) (*rtc.SignalSendCustomSignalResp, error) {
	inv, err := s.db.GetInvitationByRoomID(ctx, req.RoomID)
	if err != nil {
		log.ZWarn(ctx, "GetInvitationByRoomID failed for custom signal", err, "roomID", req.RoomID)
		return &rtc.SignalSendCustomSignalResp{}, nil
	}
	opUserID := mcontext.GetOpUserID(ctx)
	// Fix P3: 处理 json.Marshal 错误
	content, err := json.Marshal(map[string]any{
		"roomID":     req.RoomID,
		"customInfo": req.CustomInfo,
	})
	if err != nil {
		return nil, errs.WrapMsg(err, "marshal custom signal content failed")
	}
	recipients := make([]string, 0, len(inv.InviteeUserIDList)+1)
	recipients = append(recipients, inv.InviteeUserIDList...)
	recipients = append(recipients, inv.InviterUserID)
	for _, uid := range recipients {
		if uid == opUserID {
			continue
		}
		if err := s.sendCustomSignalNotification(ctx, opUserID, uid, int32(constant.SingleChatType), content); err != nil {
			log.ZWarn(ctx, "sendCustomSignalNotification failed", err, "to", uid)
		}
	}
	return &rtc.SignalSendCustomSignalResp{}, nil
}

// GetSignalInvitationRecords returns paginated call history.
func (s *rtcServer) GetSignalInvitationRecords(ctx context.Context, req *rtc.GetSignalInvitationRecordsReq) (*rtc.GetSignalInvitationRecordsResp, error) {
	total, records, err := s.db.SearchRecords(ctx, req.SendID, req.RecvID, req.SessionType, req.StartTime, req.EndTime, req.Pagination)
	if err != nil {
		return nil, err
	}
	signalRecords := datautil.Slice(records, func(r *model.SignalRecord) *rtc.SignalRecord {
		return &rtc.SignalRecord{
			RoomID:              r.RoomID,
			SID:                 r.SID,
			FileName:            r.FileName,
			MediaType:           r.MediaType,
			SessionType:         r.SessionType,
			InviterUserID:       r.InviterUserID,
			InviterUserNickname: r.InviterUserNickname,
			GroupID:             r.GroupID,
			GroupName:           r.GroupName,
			CreateTime:          r.CreateTime,
			EndTime:             r.EndTime,
			Size:                r.FileSize,
			FileURL:             r.FileURL,
		}
	})
	return &rtc.GetSignalInvitationRecordsResp{
		Total:         int32(total),
		SignalRecords: signalRecords,
	}, nil
}

// DeleteSignalRecords removes call history records by their SIDs.
func (s *rtcServer) DeleteSignalRecords(ctx context.Context, req *rtc.DeleteSignalRecordsReq) (*rtc.DeleteSignalRecordsResp, error) {
	if err := s.db.DeleteRecords(ctx, req.SIDs); err != nil {
		return nil, err
	}
	return &rtc.DeleteSignalRecordsResp{}, nil
}

// ---- helpers ----

// genToken generates a LiveKit access token for the given room and identity.
func (s *rtcServer) genToken(roomID, userID string) (string, error) {
	lk := s.config.RpcConfig.LiveKit
	at := auth.NewAccessToken(lk.APIKey, lk.APISecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomID,
	}
	at.SetVideoGrant(grant).
		SetIdentity(userID).
		SetValidFor(s.tokenExpiry)
	return at.ToJWT()
}

// signalingMsgOptions 返回信令通知消息应设置的 Options。
//
// Fix P2+P2(安全): 原代码传 make(map[string]bool) 空 map，导致：
//  1. IsNotificationByMsg 将信令消息误判为普通聊天消息，触发黑名单/好友关系等权限拦截
//  2. IsHistory/IsPersistent 默认为 true，信令消息被写入历史记录占用存储
//  3. IsUnreadCount/IsConversationUpdate 默认 true，污染未读数和会话列表
//
// 信令消息应走 Notification 通道（对话 ID 前缀 "n_"），绕过聊天消息权限校验，
// 且不写历史、不计未读、不更新会话。离线推送根据 offlinePushInfo 控制，此处不强制关闭。
func signalingMsgOptions() map[string]bool {
	opts := make(map[string]bool, 8)
	// IsNotNotification=false 表示"这是通知消息"，让 IsNotificationByMsg 返回 true
	// 从而跳过 modifyMessageByUserMessageReceiveOpt 中的黑名单/好友关系等校验
	datautil.SetSwitchFromOptions(opts, constant.IsNotNotification, false)
	datautil.SetSwitchFromOptions(opts, constant.IsSendMsg, false)
	datautil.SetSwitchFromOptions(opts, constant.IsHistory, false)
	datautil.SetSwitchFromOptions(opts, constant.IsPersistent, false)
	datautil.SetSwitchFromOptions(opts, constant.IsUnreadCount, false)
	datautil.SetSwitchFromOptions(opts, constant.IsConversationUpdate, false)
	datautil.SetSwitchFromOptions(opts, constant.IsSenderConversationUpdate, false)
	datautil.SetSwitchFromOptions(opts, constant.IsSenderSync, false)
	return opts
}

// sendSignalingNotification sends a SignalingNotification message to a user via the msg service.
// groupID 在 SessionType 为群类型（如 ReadGroupChatType）时必须非空，否则 msg 服务群聊校验会失败。
func (s *rtcServer) sendSignalingNotification(ctx context.Context, sendID, recvID string, sessionType int32, groupID string, offlinePush *sdkws.OfflinePushInfo, content []byte) error {
	now := time.Now().UnixMilli()
	msgData := &sdkws.MsgData{
		SendID:      sendID,
		RecvID:      recvID,
		SessionType: sessionType,
		GroupID:     groupID,
		ContentType: int32(constant.SignalingNotification),
		MsgFrom:     int32(constant.SysMsgType),
		Content:     content,
		CreateTime:  now,
		SendTime:    now,
		ServerMsgID: uuid.New().String(),
		ClientMsgID: uuid.New().String(),
		Options:     signalingMsgOptions(),
	}
	if offlinePush != nil {
		msgData.OfflinePushInfo = offlinePush
	}

	_, err := s.msgClient.MsgClient.SendMsg(ctx, &pbmsg.SendMsgReq{MsgData: msgData})
	if err != nil {
		log.ZError(ctx, "sendSignalingNotification", err, "msgdata", msgData)
		return err
	}
	log.ZInfo(ctx, "sendSignalingNotification", "msgData", msgData)

	return nil
}

// sendCustomSignalNotification sends a CustomSignalNotification (1605) to a user.
func (s *rtcServer) sendCustomSignalNotification(ctx context.Context, sendID, recvID string, sessionType int32, content []byte) error {
	now := time.Now().UnixMilli()
	msgData := &sdkws.MsgData{
		SendID:      sendID,
		RecvID:      recvID,
		SessionType: sessionType,
		ContentType: int32(constant.CustomSignalNotification),
		MsgFrom:     int32(constant.SysMsgType),
		Content:     content,
		CreateTime:  now,
		SendTime:    now,
		ServerMsgID: uuid.New().String(),
		ClientMsgID: uuid.New().String(),
		Options:     signalingMsgOptions(),
	}
	_, err := s.msgClient.MsgClient.SendMsg(ctx, &pbmsg.SendMsgReq{MsgData: msgData})
	return err
}

// marshalSignalReq serializes a SignalReq to bytes.
// Fix P2: 原代码使用 _ 吞掉错误，序列化失败时返回 nil，导致被叫收到空 Content 消息，来电通知丢失。
func marshalSignalReq(req *rtc.SignalReq) ([]byte, error) {
	b, err := proto.Marshal(req)
	if err != nil {
		return nil, errs.WrapMsg(err, "marshal SignalReq failed")
	}
	return b, nil
}

// newRoomID generates a unique room ID.
func newRoomID() string {
	return fmt.Sprintf("room-%s", uuid.New().String())
}

// invitationToModel converts a proto InvitationInfo to the database model.
func invitationToModel(inv *rtc.InvitationInfo, push *sdkws.OfflinePushInfo) *model.SignalInvitation {
	now := time.Now()
	m := &model.SignalInvitation{
		RoomID:             inv.RoomID,
		InviterUserID:      inv.InviterUserID,
		InviteeUserIDList:  inv.InviteeUserIDList,
		CustomData:         inv.CustomData,
		GroupID:            inv.GroupID,
		Timeout:            inv.Timeout,
		MediaType:          inv.MediaType,
		PlatformID:         inv.PlatformID,
		SessionType:        inv.SessionType,
		InitiateTime:       inv.InitiateTime,
		BusyLineUserIDList: inv.BusyLineUserIDList,
		CreateTime:         now.UnixMilli(),
		ExpireAt:           now.Add(time.Duration(inv.Timeout+30) * time.Second),
	}
	if push != nil {
		m.OfflinePushTitle = push.Title
		m.OfflinePushDesc = push.Desc
		m.OfflinePushEx = push.Ex
	}
	return m
}

// modelToInvitationInfo converts a database model to proto InvitationInfo.
func modelToInvitationInfo(m *model.SignalInvitation) *rtc.InvitationInfo {
	if m == nil {
		return nil
	}
	return &rtc.InvitationInfo{
		InviterUserID:      m.InviterUserID,
		InviteeUserIDList:  m.InviteeUserIDList,
		CustomData:         m.CustomData,
		GroupID:            m.GroupID,
		RoomID:             m.RoomID,
		Timeout:            m.Timeout,
		MediaType:          m.MediaType,
		PlatformID:         m.PlatformID,
		SessionType:        m.SessionType,
		InitiateTime:       m.InitiateTime,
		BusyLineUserIDList: m.BusyLineUserIDList,
	}
}

// hungUpPeerIDsFromDB returns IDs that should receive hang-up notification, based on authoritative DB data.
func hungUpPeerIDsFromDB(inv *model.SignalInvitation, callerID string) []string {
	if callerID == inv.InviterUserID {
		return inv.InviteeUserIDList
	}
	return []string{inv.InviterUserID}
}

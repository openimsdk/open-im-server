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
		return nil, respErr
	}
	return &rtc.SignalMessageAssembleResp{SignalResp: &resp}, nil
}

// handleInvite processes a 1-to-1 call invitation.
func (s *rtcServer) handleInvite(ctx context.Context, req *rtc.SignalInviteReq, signalReq *rtc.SignalReq) (*rtc.SignalInviteResp, error) {
	inv := req.Invitation
	if inv == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}
	if inv.RoomID == "" {
		inv.RoomID = newRoomID()
	}
	inv.InviterUserID = req.UserID
	inv.InitiateTime = time.Now().UnixMilli()

	if _, err := s.roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{Name: inv.RoomID}); err != nil {
		return nil, errs.WrapMsg(err, "LiveKit CreateRoom failed", "roomID", inv.RoomID)
	}

	token, err := s.genToken(inv.RoomID, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.db.CreateInvitation(ctx, invitationToModel(inv, req.OfflinePushInfo)); err != nil {
		log.ZWarn(ctx, "CreateInvitation failed", err, "roomID", inv.RoomID)
	}

	content := marshalSignalReq(signalReq)
	for _, inviteeID := range inv.InviteeUserIDList {
		if err := s.sendSignalingNotification(ctx, req.UserID, inviteeID, int32(constant.SingleChatType), req.OfflinePushInfo, content); err != nil {
			log.ZWarn(ctx, "sendSignalingNotification to invitee failed", err, "inviteeID", inviteeID)
		}
	}

	return &rtc.SignalInviteResp{
		Token:   token,
		RoomID:  inv.RoomID,
		LiveURL: s.config.RpcConfig.LiveKit.ExternalAddress,
	}, nil
}

// handleInviteInGroup processes a group call invitation.
func (s *rtcServer) handleInviteInGroup(ctx context.Context, req *rtc.SignalInviteInGroupReq, signalReq *rtc.SignalReq) (*rtc.SignalInviteInGroupResp, error) {
	inv := req.Invitation
	if inv == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}
	if inv.RoomID == "" {
		inv.RoomID = newRoomID()
	}
	inv.InviterUserID = req.UserID
	inv.InitiateTime = time.Now().UnixMilli()

	if _, err := s.roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{Name: inv.RoomID}); err != nil {
		return nil, errs.WrapMsg(err, "LiveKit CreateRoom failed", "roomID", inv.RoomID)
	}

	token, err := s.genToken(inv.RoomID, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.db.CreateInvitation(ctx, invitationToModel(inv, req.OfflinePushInfo)); err != nil {
		log.ZWarn(ctx, "CreateInvitation failed", err, "roomID", inv.RoomID)
	}

	content := marshalSignalReq(signalReq)
	for _, inviteeID := range inv.InviteeUserIDList {
		if err := s.sendSignalingNotification(ctx, req.UserID, inviteeID, int32(constant.ReadGroupChatType), req.OfflinePushInfo, content); err != nil {
			log.ZWarn(ctx, "sendSignalingNotification to group invitee failed", err, "inviteeID", inviteeID)
		}
	}

	return &rtc.SignalInviteInGroupResp{
		Token:   token,
		RoomID:  inv.RoomID,
		LiveURL: s.config.RpcConfig.LiveKit.ExternalAddress,
	}, nil
}

// handleAccept processes a call acceptance.
func (s *rtcServer) handleAccept(ctx context.Context, req *rtc.SignalAcceptReq, signalReq *rtc.SignalReq) (*rtc.SignalAcceptResp, error) {
	inv := req.Invitation
	if inv == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}

	token, err := s.genToken(inv.RoomID, req.UserID)
	if err != nil {
		return nil, err
	}

	sessionType := int32(constant.SingleChatType)
	if inv.GroupID != "" {
		sessionType = int32(constant.ReadGroupChatType)
	}
	content := marshalSignalReq(signalReq)
	if err := s.sendSignalingNotification(ctx, req.UserID, inv.InviterUserID, sessionType, req.OfflinePushInfo, content); err != nil {
		log.ZWarn(ctx, "sendSignalingNotification accept to inviter failed", err, "inviterID", inv.InviterUserID)
	}

	return &rtc.SignalAcceptResp{
		Token:   token,
		RoomID:  inv.RoomID,
		LiveURL: s.config.RpcConfig.LiveKit.ExternalAddress,
	}, nil
}

// handleReject processes a call rejection.
func (s *rtcServer) handleReject(ctx context.Context, req *rtc.SignalRejectReq, signalReq *rtc.SignalReq) (*rtc.SignalRejectResp, error) {
	inv := req.Invitation
	if inv == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}

	sessionType := int32(constant.SingleChatType)
	if inv.GroupID != "" {
		sessionType = int32(constant.ReadGroupChatType)
	}
	content := marshalSignalReq(signalReq)
	if err := s.sendSignalingNotification(ctx, req.UserID, inv.InviterUserID, sessionType, req.OfflinePushInfo, content); err != nil {
		log.ZWarn(ctx, "sendSignalingNotification reject to inviter failed", err, "inviterID", inv.InviterUserID)
	}

	if inv.GroupID != "" {
		if err := s.db.RemoveInvitee(ctx, inv.RoomID, req.UserID); err != nil {
			log.ZWarn(ctx, "RemoveInvitee failed", err, "roomID", inv.RoomID, "userID", req.UserID)
		}
	} else {
		if err := s.db.DeleteInvitation(ctx, inv.RoomID); err != nil {
			log.ZWarn(ctx, "DeleteInvitation failed", err, "roomID", inv.RoomID)
		}
	}

	return &rtc.SignalRejectResp{}, nil
}

// handleCancel processes a call cancellation.
func (s *rtcServer) handleCancel(ctx context.Context, req *rtc.SignalCancelReq, signalReq *rtc.SignalReq) (*rtc.SignalCancelResp, error) {
	inv := req.Invitation
	if inv == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}

	sessionType := int32(constant.SingleChatType)
	if inv.GroupID != "" {
		sessionType = int32(constant.ReadGroupChatType)
	}
	content := marshalSignalReq(signalReq)
	for _, inviteeID := range inv.InviteeUserIDList {
		if err := s.sendSignalingNotification(ctx, req.UserID, inviteeID, sessionType, req.OfflinePushInfo, content); err != nil {
			log.ZWarn(ctx, "sendSignalingNotification cancel to invitee failed", err, "inviteeID", inviteeID)
		}
	}

	if err := s.db.DeleteInvitation(ctx, inv.RoomID); err != nil {
		log.ZWarn(ctx, "DeleteInvitation failed", err, "roomID", inv.RoomID)
	}

	return &rtc.SignalCancelResp{}, nil
}

// handleHungUp processes a call hang-up.
func (s *rtcServer) handleHungUp(ctx context.Context, req *rtc.SignalHungUpReq, signalReq *rtc.SignalReq) (*rtc.SignalHungUpResp, error) {
	inv := req.Invitation
	if inv == nil {
		return nil, errs.ErrArgs.WrapMsg("invitation is nil")
	}

	sessionType := int32(constant.SingleChatType)
	if inv.GroupID != "" {
		sessionType = int32(constant.ReadGroupChatType)
	}
	content := marshalSignalReq(signalReq)
	for _, peerID := range hungUpPeerIDs(inv, req.UserID) {
		if err := s.sendSignalingNotification(ctx, req.UserID, peerID, sessionType, req.OfflinePushInfo, content); err != nil {
			log.ZWarn(ctx, "sendSignalingNotification hungUp to peer failed", err, "peerID", peerID)
		}
	}

	// Terminate the LiveKit room
	if _, err := s.roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{Room: inv.RoomID}); err != nil {
		log.ZWarn(ctx, "LiveKit DeleteRoom failed", err, "roomID", inv.RoomID)
	}

	if err := s.db.DeleteInvitation(ctx, inv.RoomID); err != nil {
		log.ZWarn(ctx, "DeleteInvitation failed", err, "roomID", inv.RoomID)
	}

	return &rtc.SignalHungUpResp{}, nil
}

// handleGetTokenByRoomID returns a LiveKit token for an existing room.
func (s *rtcServer) handleGetTokenByRoomID(ctx context.Context, req *rtc.SignalGetTokenByRoomIDReq) (*rtc.SignalGetTokenByRoomIDResp, error) {
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
func (s *rtcServer) SignalGetTokenByRoomID(ctx context.Context, req *rtc.SignalGetTokenByRoomIDReq) (*rtc.SignalGetTokenByRoomIDResp, error) {
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
	content, _ := json.Marshal(map[string]any{
		"roomID":     req.RoomID,
		"customInfo": req.CustomInfo,
	})
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

// sendSignalingNotification sends a SignalingNotification message to a user via the msg service.
func (s *rtcServer) sendSignalingNotification(ctx context.Context, sendID, recvID string, sessionType int32, offlinePush *sdkws.OfflinePushInfo, content []byte) error {
	now := time.Now().UnixMilli()
	msgData := &sdkws.MsgData{
		SendID:      sendID,
		RecvID:      recvID,
		SessionType: sessionType,
		ContentType: int32(constant.SignalingNotification),
		MsgFrom:     int32(constant.SysMsgType),
		Content:     content,
		CreateTime:  now,
		SendTime:    now,
		ServerMsgID: uuid.New().String(),
		ClientMsgID: uuid.New().String(),
		Options:     make(map[string]bool),
	}
	if offlinePush != nil {
		msgData.OfflinePushInfo = offlinePush
	}
	_, err := s.msgClient.MsgClient.SendMsg(ctx, &pbmsg.SendMsgReq{MsgData: msgData})
	return err
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
		Options:     make(map[string]bool),
	}
	_, err := s.msgClient.MsgClient.SendMsg(ctx, &pbmsg.SendMsgReq{MsgData: msgData})
	return err
}

func marshalSignalReq(req *rtc.SignalReq) []byte {
	b, _ := proto.Marshal(req)
	return b
}

// newRoomID generates a unique room ID.
func newRoomID() string {
	return fmt.Sprintf("room-%s", uuid.New().String())
}

// invitationToModel converts a proto InvitationInfo to the database model.
func invitationToModel(inv *rtc.InvitationInfo, push *sdkws.OfflinePushInfo) *model.SignalInvitation {
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
		CreateTime:         time.Now().UnixMilli(),
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

// hungUpPeerIDs returns the IDs that should receive hang-up notification.
func hungUpPeerIDs(inv *rtc.InvitationInfo, callerID string) []string {
	if callerID == inv.InviterUserID {
		return inv.InviteeUserIDList
	}
	return []string{inv.InviterUserID}
}

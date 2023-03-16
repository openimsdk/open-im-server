package third

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
)

func (t *thirdServer) GetSignalInvitationInfo(ctx context.Context, req *third.GetSignalInvitationInfoReq) (resp *third.GetSignalInvitationInfoResp, err error) {
	signalReq, err := t.thirdDatabase.GetSignalInvitationInfoByClientMsgID(ctx, req.ClientMsgID)
	if err != nil {
		return nil, err
	}
	resp = &third.GetSignalInvitationInfoResp{}
	resp.InvitationInfo = signalReq.Invitation
	resp.OfflinePushInfo = signalReq.OfflinePushInfo
	return resp, nil
}

func (t *thirdServer) GetSignalInvitationInfoStartApp(ctx context.Context, req *third.GetSignalInvitationInfoStartAppReq) (resp *third.GetSignalInvitationInfoStartAppResp, err error) {
	signalReq, err := t.thirdDatabase.GetAvailableSignalInvitationInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	resp = &third.GetSignalInvitationInfoStartAppResp{}
	resp.InvitationInfo = signalReq.Invitation
	resp.OfflinePushInfo = signalReq.OfflinePushInfo
	return resp, nil
}

package rtc

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbRtc "Open_IM/pkg/proto/rtc"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/livekit/protocol/auth"
	//lksdk "github.com/livekit/server-sdk-go"
)

type RtcService struct {
}

func (r *RtcService) GetJoinToken(_ context.Context, req *pbRtc.GetJoinTokenReq) (resp *pbRtc.GetJoinTokenResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbRtc.GetJoinTokenResp{}
	canPublish := true
	canSubscribe := true
	at := auth.NewAccessToken(req.ApiKey, req.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         req.Room,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	}
	byte, err := json.Marshal(req.MetaData)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "json marshal failed", err.Error())
		resp.CommonResp = &pbRtc.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}
		return
	}
	at.AddGrant(grant).
		SetIdentity(req.Identity).
		// optional
		SetName("participant-name").
		SetValidFor(time.Hour).SetMetadata(string(byte))
	jwt, err := at.ToJWT()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "toJwt failed", err.Error(), "jwt: ", jwt)
	}
	resp.Jwt = jwt
	resp.CommonResp = &pbRtc.CommonResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, err
}

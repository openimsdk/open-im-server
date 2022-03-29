package open_im_media

import (
	pbRtc "Open_IM/pkg/proto/rtc"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:11300"
)

//var roomClient *lksdk.RoomServiceClient

type Media struct {
}

func NewMedia() *Media {
	return &Media{}
}

func (m *Media) GetJoinToken(room, identity string, operationID string, data *open_im_sdk.ParticipantMetaData) (string, string, error) {
	var newData pbRtc.ParticipantMetaData
	copier.Copy(&newData, data)
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		return "", "", err
	}
	defer conn.Close()
	c := pbRtc.NewRtcServiceClient(conn)
	req := &pbRtc.GetJoinTokenReq{Room: room, OperationID: operationID, Identity: identity, MetaData: &newData}
	resp, err := c.GetJoinToken(context.Background(), req)
	if err != nil {
		return "", "", err
	}
	if resp.CommonResp.ErrCode != 0 {
		return "", "", errors.New(resp.CommonResp.ErrMsg)
	}
	return resp.Jwt, resp.LiveURL, nil
	//at := auth.NewAccessToken(m.ApiKey, m.ApiSecret)
	//grant := &auth.VideoGrant{
	//	RoomJoin: true,
	//	Room:     room,
	//}
	//at.AddGrant(grant).
	//	SetIdentity(identity).
	//	SetValidFor(time.Hour)
	//
	//return at.ToJWT()
}

func init() {
	//roomClient = lksdk.NewRoomServiceClient(MediaAddress, ApiKey, ApiSecret)
}

func (m *Media) CreateRoom(roomName string) (error, error) {
	return nil, nil
	//return roomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
	//	Name:         roomName,
	//	EmptyTimeout: 60 * 3,
	//})

}

package open_im_media

import (
	pbRtc "Open_IM/pkg/proto/rtc"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"
	"errors"
	"google.golang.org/grpc"
)

const (
	MediaAddress = "ws://43.128.5.63:7880"
	ApiKey       = "APIGPW3gnFTzqHH"
	ApiSecret    = "23ztfSqsfQ8hKkHzHTl3Z4bvaxro0snjk5jwbp5p6Q3"
	// Address gRPC服务地址
	Address = "127.0.0.1:11300"
)

//var roomClient *lksdk.RoomServiceClient

type Media struct {
	MediaAddress string
	ApiKey       string
	ApiSecret    string
}

func NewMedia() *Media {
	return &Media{MediaAddress: MediaAddress,
		ApiKey:    ApiKey,
		ApiSecret: ApiSecret}
}
func (m *Media) GetUrl() string {
	return m.MediaAddress
}

func (m *Media) GetJoinToken(room, identity string, operationID string, data *open_im_sdk.ParticipantMetaData) (string, error) {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pbRtc.NewRtcServiceClient(conn)
	req := &pbRtc.GetJoinTokenReq{ApiKey: m.ApiKey, ApiSecret: m.ApiSecret, Room: room, OperationID: operationID, Identity: identity, MetaData: data}
	resp, err := c.GetJoinToken(context.Background(), req)
	if err != nil {
		return "", err
	}
	if resp.CommonResp.ErrCode != 0 {
		return "", errors.New(resp.CommonResp.ErrMsg)
	}
	return resp.Jwt, nil
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

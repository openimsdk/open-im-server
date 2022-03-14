package open_im_media

const (
	MediaAddress = "ws://43.128.5.63:7880"
	ApiKey       = "APIGPW3gnFTzqHH"
	ApiSecret    = "23ztfSqsfQ8hKkHzHTl3Z4bvaxro0snjk5jwbp5p6Q3"
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
func (m *Media) GetJoinToken(room, identity string) (string, error) {
	return identity, nil
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

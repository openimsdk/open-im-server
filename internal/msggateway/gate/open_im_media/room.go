package open_im_media

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:11300"
)

type Media struct {
}

func NewMedia() *Media {
	return &Media{}
}

func init() {

}

func (m *Media) CreateRoom(roomName string) (error, error) {
	return nil, nil

}

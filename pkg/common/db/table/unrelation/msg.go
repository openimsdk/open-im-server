package unrelation

import "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"

type MsgModel struct {
	SendID           string                 `bson:"send_id"`
	RecvID           string                 `bson:"recv_id"`
	GroupID          string                 `bson:"group_id"`
	ClientMsgID      string                 `bson:"client_msg_id"` // 客户端消息ID
	ServerMsgID      string                 `bson:"server_msg_id"` // 服务端消息ID
	SenderPlatformID int32                  `bson:"sender_platform_id"`
	SenderNickname   string                 `bson:"sender_nickname"`
	SenderFaceURL    string                 `bson:"sender_face_url"`
	SessionType      int32                  `bson:"session_type"`
	MsgFrom          int32                  `bson:"msg_from"`
	ContentType      int32                  `bson:"contentType"`
	Content          []byte                 `bson:"content"`
	Seq              int64                  `bson:"seq"`
	SendTime         int64                  `bson:"sendTime"`
	CreateTime       int64                  `bson:"createTime"`
	Status           int32                  `bson:"status"`
	Options          map[string]bool        `bson:"options"`
	OfflinePushInfo  *sdkws.OfflinePushInfo `bson:"offlinePushInfo"`
	AtUserIDList     []string               `bson:"atUserIDList"`
	MsgDataList      []byte                 `bson:"msgDataList"`
	AttachedInfo     string                 `bson:"attachedInfo"`
	Ex               string                 `bson:"ex"`
}

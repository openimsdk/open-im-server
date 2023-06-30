package apistruct

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type DelMsgReq struct {
	UserID      string   `json:"userID,omitempty" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty" binding:"required"`
	OperationID string   `json:"operationID,omitempty" binding:"required"`
}

type DelMsgResp struct {
}

type CleanUpMsgReq struct {
	UserID      string `json:"userID"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type CleanUpMsgResp struct {
}

type DelSuperGroupMsgReq struct {
	UserID      string   `json:"userID" binding:"required"`
	GroupID     string   `json:"groupID" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty"`
	IsAllDelete bool     `json:"isAllDelete"`
	OperationID string   `json:"operationID" binding:"required"`
}

type DelSuperGroupMsgResp struct {
}

type MsgDeleteNotificationElem struct {
	GroupID     string   `json:"groupID"`
	IsAllDelete bool     `json:"isAllDelete"`
	SeqList     []uint32 `json:"seqList"`
}

type SetMsgMinSeqReq struct {
	UserID      string `json:"userID"  binding:"required"`
	GroupID     string `json:"groupID"`
	MinSeq      uint32 `json:"minSeq"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type SetMsgMinSeqResp struct {
}

type ModifyMessageReactionExtensionsReq struct {
	OperationID           string                     `json:"operationID" binding:"required"`
	conversationID        string                     `json:"conversationID"  binding:"required"`
	SessionType           int32                      `json:"sessionType" binding:"required"`
	ReactionExtensionList map[string]*sdkws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID           string                     `json:"clientMsgID" binding:"required"`
	Ex                    *string                    `json:"ex"`
	AttachedInfo          *string                    `json:"attachedInfo"`
	IsReact               bool                       `json:"isReact"`
	IsExternalExtensions  bool                       `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                      `json:"msgFirstModifyTime"`
}

type ModifyMessageReactionExtensionsResp struct {
	Data struct {
		ResultKeyValue     []*msg.KeyValueResp `json:"result"`
		MsgFirstModifyTime int64               `json:"msgFirstModifyTime"`
		IsReact            bool                `json:"isReact"`
	} `json:"data"`
}

//type OperateMessageListReactionExtensionsReq struct {
//	OperationID            string                                                        `json:"operationID" binding:"required"`
//	conversationID               string                                                        `json:"conversationID"  binding:"required"`
//	SessionType            string                                                        `json:"sessionType" binding:"required"`
//	MessageReactionKeyList []*msg.GetMessageListReactionExtensionsReq_MessageReactionKey `json:"messageReactionKeyList" binding:"required"`
//}

type OperateMessageListReactionExtensionsResp struct {
	Data struct {
		SuccessList []*msg.ExtendMsgResp `json:"successList"`
		FailedList  []*msg.ExtendMsgResp `json:"failedList"`
	} `json:"data"`
}

type SetMessageReactionExtensionsCallbackReq ModifyMessageReactionExtensionsReq

type SetMessageReactionExtensionsCallbackResp ModifyMessageReactionExtensionsResp

//type GetMessageListReactionExtensionsReq OperateMessageListReactionExtensionsReq

type GetMessageListReactionExtensionsResp struct {
	Data []*msg.SingleMessageExtensionResult `json:"data"`
}

type AddMessageReactionExtensionsReq ModifyMessageReactionExtensionsReq

type AddMessageReactionExtensionsResp ModifyMessageReactionExtensionsResp

type DeleteMessageReactionExtensionsReq struct {
	OperationID           string            `json:"operationID" binding:"required"`
	conversationID        string            `json:"conversationID" binding:"required"`
	SessionType           int32             `json:"sessionType" binding:"required"`
	ClientMsgID           string            `json:"clientMsgID" binding:"required"`
	IsExternalExtensions  bool              `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64             `json:"msgFirstModifyTime" binding:"required"`
	ReactionExtensionList []*sdkws.KeyValue `json:"reactionExtensionList" binding:"required"`
}

type DeleteMessageReactionExtensionsResp struct {
	Data []*msg.KeyValueResp
}

type PictureBaseInfo struct {
	UUID   string `mapstructure:"uuid"`
	Type   string `mapstructure:"type" `
	Size   int64  `mapstructure:"size" `
	Width  int32  `mapstructure:"width" `
	Height int32  `mapstructure:"height"`
	Url    string `mapstructure:"url" `
}

type PictureElem struct {
	SourcePath      string          `mapstructure:"sourcePath"`
	SourcePicture   PictureBaseInfo `mapstructure:"sourcePicture"`
	BigPicture      PictureBaseInfo `mapstructure:"bigPicture" `
	SnapshotPicture PictureBaseInfo `mapstructure:"snapshotPicture"`
}
type SoundElem struct {
	UUID      string `mapstructure:"uuid"`
	SoundPath string `mapstructure:"soundPath"`
	SourceURL string `mapstructure:"sourceUrl"`
	DataSize  int64  `mapstructure:"dataSize"`
	Duration  int64  `mapstructure:"duration"`
}
type VideoElem struct {
	VideoPath      string `mapstructure:"videoPath"`
	VideoUUID      string `mapstructure:"videoUUID"`
	VideoURL       string `mapstructure:"videoUrl"`
	VideoType      string `mapstructure:"videoType"`
	VideoSize      int64  `mapstructure:"videoSize"`
	Duration       int64  `mapstructure:"duration"`
	SnapshotPath   string `mapstructure:"snapshotPath"`
	SnapshotUUID   string `mapstructure:"snapshotUUID"`
	SnapshotSize   int64  `mapstructure:"snapshotSize"`
	SnapshotURL    string `mapstructure:"snapshotUrl"`
	SnapshotWidth  int32  `mapstructure:"snapshotWidth"`
	SnapshotHeight int32  `mapstructure:"snapshotHeight"`
}
type FileElem struct {
	FilePath  string `mapstructure:"filePath"`
	UUID      string `mapstructure:"uuid"`
	SourceURL string `mapstructure:"sourceUrl"`
	FileName  string `mapstructure:"fileName"`
	FileSize  int64  `mapstructure:"fileSize"`
}
type AtElem struct {
	Text       string   `mapstructure:"text"`
	AtUserList []string `mapstructure:"atUserList"`
	IsAtSelf   bool     `mapstructure:"isAtSelf"`
}
type LocationElem struct {
	Description string  `mapstructure:"description"`
	Longitude   float64 `mapstructure:"longitude"`
	Latitude    float64 `mapstructure:"latitude"`
}
type CustomElem struct {
	Data        string `mapstructure:"data" validate:"required"`
	Description string `mapstructure:"description"`
	Extension   string `mapstructure:"extension"`
}
type TextElem struct {
	Text string `mapstructure:"text" validate:"required"`
}

type RevokeElem struct {
	RevokeMsgClientID string `mapstructure:"revokeMsgClientID" validate:"required"`
}
type OANotificationElem struct {
	NotificationName    string      `mapstructure:"notificationName" json:"notificationName" validate:"required"`
	NotificationFaceURL string      `mapstructure:"notificationFaceURL" json:"notificationFaceURL"`
	NotificationType    int32       `mapstructure:"notificationType" json:"notificationType" validate:"required"`
	Text                string      `mapstructure:"text" json:"text" validate:"required"`
	Url                 string      `mapstructure:"url" json:"url"`
	MixType             int32       `mapstructure:"mixType" json:"mixType"`
	PictureElem         PictureElem `mapstructure:"pictureElem" json:"pictureElem"`
	SoundElem           SoundElem   `mapstructure:"soundElem" json:"soundElem"`
	VideoElem           VideoElem   `mapstructure:"videoElem" json:"videoElem"`
	FileElem            FileElem    `mapstructure:"fileElem" json:"fileElem"`
	Ex                  string      `mapstructure:"ex" json:"ex"`
}
type MessageRevoked struct {
	RevokerID       string `mapstructure:"revokerID" json:"revokerID" validate:"required"`
	RevokerRole     int32  `mapstructure:"revokerRole" json:"revokerRole" validate:"required"`
	ClientMsgID     string `mapstructure:"clientMsgID" json:"clientMsgID" validate:"required"`
	RevokerNickname string `mapstructure:"revokerNickname" json:"revokerNickname"`
	SessionType     int32  `mapstructure:"sessionType" json:"sessionType" validate:"required"`
	Seq             uint32 `mapstructure:"seq" json:"seq" validate:"required"`
}

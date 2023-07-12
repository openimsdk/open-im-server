// Copyright © 2023 OpenIM. All rights reserved.
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

package apistruct

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

// define a  DelMsgReq struct
type DelMsgReq struct {
	UserID      string   `json:"userID,omitempty"      binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty"     binding:"required"`
	OperationID string   `json:"operationID,omitempty" binding:"required"`
}

// define a DelMsgResp struct
type DelMsgResp struct {
}

// define a CleanUpMsgReq struct
type CleanUpMsgReq struct {
	UserID      string `json:"userID"      binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

// define a CleanUpMsgResp struct
type CleanUpMsgResp struct {
}

// define a DelSuperGroupMsgReq struct
type DelSuperGroupMsgReq struct {
	UserID      string   `json:"userID"            binding:"required"`
	GroupID     string   `json:"groupID"           binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty"`
	IsAllDelete bool     `json:"isAllDelete"`
	OperationID string   `json:"operationID"       binding:"required"`
}

// define a DelSuperGroupMsgResp struct
type DelSuperGroupMsgResp struct {
}

// difine a MsgDeleteNotificationElem struct
type MsgDeleteNotificationElem struct {
	GroupID     string   `json:"groupID"`
	IsAllDelete bool     `json:"isAllDelete"`
	SeqList     []uint32 `json:"seqList"`
}

// define a SetMsgMinSeqReq struct
type SetMsgMinSeqReq struct {
	UserID      string `json:"userID"      binding:"required"`
	GroupID     string `json:"groupID"`
	MinSeq      uint32 `json:"minSeq"      binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

// define a SetMsgMinSeqResp struct
type SetMsgMinSeqResp struct {
}

// define s ModifyMessageReactionExtensionsReq struct
type ModifyMessageReactionExtensionsReq struct {
	OperationID           string                     `json:"operationID"                     binding:"required"`
	conversationID        string                     `json:"conversationID"                  binding:"required"`
	SessionType           int32                      `json:"sessionType"                     binding:"required"`
	ReactionExtensionList map[string]*sdkws.KeyValue `json:"reactionExtensionList,omitempty" binding:"required"`
	ClientMsgID           string                     `json:"clientMsgID"                     binding:"required"`
	Ex                    *string                    `json:"ex"`
	AttachedInfo          *string                    `json:"attachedInfo"`
	IsReact               bool                       `json:"isReact"`
	IsExternalExtensions  bool                       `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64                      `json:"msgFirstModifyTime"`
}

// define a ModifyMessageReactionExtensionsResp struct
type ModifyMessageReactionExtensionsResp struct {
	Data struct {
		ResultKeyValue     []*msg.KeyValueResp `json:"result"`
		MsgFirstModifyTime int64               `json:"msgFirstModifyTime"`
		IsReact            bool                `json:"isReact"`
	} `json:"data"`
}

//type OperateMessageListReactionExtensionsReq struct {
// 	OperationID            string                                                        `json:"operationID"
// binding:"required"` 	conversationID               string
// `json:"conversationID"  binding:"required"` 	SessionType            string
//             `json:"sessionType" binding:"required"` 	MessageReactionKeyList
// []*msg.GetMessageListReactionExtensionsReq_MessageReactionKey `json:"messageReactionKeyList" binding:"required"`
//}

type OperateMessageListReactionExtensionsResp struct {
	Data struct {
		SuccessList []*msg.ExtendMsgResp `json:"successList"`
		FailedList  []*msg.ExtendMsgResp `json:"failedList"`
	} `json:"data"`
}

// renamed SetMessageReactionExtensionsCallbackResp
type SetMessageReactionExtensionsCallbackReq ModifyMessageReactionExtensionsReq

// renamed SetMessageReactionExtensionsCallbackResp
type SetMessageReactionExtensionsCallbackResp ModifyMessageReactionExtensionsResp

// type GetMessageListReactionExtensionsReq OperateMessageListReactionExtensionsReq
type GetMessageListReactionExtensionsResp struct {
	Data []*msg.SingleMessageExtensionResult `json:"data"`
}

// AddMessageReactionExtensionsReq struct
type AddMessageReactionExtensionsReq ModifyMessageReactionExtensionsReq

// AddMessageReactionExtensionsResp struct
type AddMessageReactionExtensionsResp ModifyMessageReactionExtensionsResp

// DeleteMessageReactionExtensionsReq struct
type DeleteMessageReactionExtensionsReq struct {
	OperationID           string            `json:"operationID"           binding:"required"`
	conversationID        string            `json:"conversationID"        binding:"required"`
	SessionType           int32             `json:"sessionType"           binding:"required"`
	ClientMsgID           string            `json:"clientMsgID"           binding:"required"`
	IsExternalExtensions  bool              `json:"isExternalExtensions"`
	MsgFirstModifyTime    int64             `json:"msgFirstModifyTime"    binding:"required"`
	ReactionExtensionList []*sdkws.KeyValue `json:"reactionExtensionList" binding:"required"`
}

// DeleteMessageReactionExtensionsResp struct
type DeleteMessageReactionExtensionsResp struct {
	Data []*msg.KeyValueResp
}

// define a picture base info struct
type PictureBaseInfo struct {
	UUID   string `mapstructure:"uuid"`
	Type   string `mapstructure:"type"`
	Size   int64  `mapstructure:"size"`
	Width  int32  `mapstructure:"width"`
	Height int32  `mapstructure:"height"`
	URL    string `mapstructure:"url"`
}

// define a picture elem struct
type PictureElem struct {
	SourcePath      string          `mapstructure:"sourcePath"`
	SourcePicture   PictureBaseInfo `mapstructure:"sourcePicture"`
	BigPicture      PictureBaseInfo `mapstructure:"bigPicture"`
	SnapshotPicture PictureBaseInfo `mapstructure:"snapshotPicture"`
}

// define a sound elem struct
type SoundElem struct {
	UUID      string `mapstructure:"uuid"`
	SoundPath string `mapstructure:"soundPath"`
	SourceURL string `mapstructure:"sourceUrl"`
	DataSize  int64  `mapstructure:"dataSize"`
	Duration  int64  `mapstructure:"duration"`
}

// define a video element struct
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

// define a file elem struct
type FileElem struct {
	FilePath  string `mapstructure:"filePath"`
	UUID      string `mapstructure:"uuid"`
	SourceURL string `mapstructure:"sourceUrl"`
	FileName  string `mapstructure:"fileName"`
	FileSize  int64  `mapstructure:"fileSize"`
}

// define a atelem struct
type AtElem struct {
	Text       string   `mapstructure:"text"`
	AtUserList []string `mapstructure:"atUserList"`
	IsAtSelf   bool     `mapstructure:"isAtSelf"`
}

// define a locatinelem struct
type LocationElem struct {
	Description string  `mapstructure:"description"`
	Longitude   float64 `mapstructure:"longitude"`
	Latitude    float64 `mapstructure:"latitude"`
}

// define a customelem struct
type CustomElem struct {
	Data        string `mapstructure:"data"        validate:"required"`
	Description string `mapstructure:"description"`
	Extension   string `mapstructure:"extension"`
}

// define a textelem struct
type TextElem struct {
	Text string `mapstructure:"text" validate:"required"`
}

// define a revoke elem struct
type RevokeElem struct {
	RevokeMsgClientID string `mapstructure:"revokeMsgClientID" validate:"required"`
}

// define a OANotificationElem struct
type OANotificationElem struct {
	NotificationName    string      `mapstructure:"notificationName"    json:"notificationName"    validate:"required"`
	NotificationFaceURL string      `mapstructure:"notificationFaceURL" json:"notificationFaceURL"`
	NotificationType    int32       `mapstructure:"notificationType"    json:"notificationType"    validate:"required"`
	Text                string      `mapstructure:"text"                json:"text"                validate:"required"`
	UrL                 string      `mapstructure:"url"                 json:"url"`
	MixType             int32       `mapstructure:"mixType"             json:"mixType"`
	PictureElem         PictureElem `mapstructure:"pictureElem"         json:"pictureElem"`
	SoundElem           SoundElem   `mapstructure:"soundElem"           json:"soundElem"`
	VideoElem           VideoElem   `mapstructure:"videoElem"           json:"videoElem"`
	FileElem            FileElem    `mapstructure:"fileElem"            json:"fileElem"`
	Ex                  string      `mapstructure:"ex"                  json:"ex"`
}

// define a message revoked struct
type MessageRevoked struct {
	RevokerID       string `mapstructure:"revokerID"       json:"revokerID"       validate:"required"`
	RevokerRole     int32  `mapstructure:"revokerRole"     json:"revokerRole"     validate:"required"`
	ClientMsgID     string `mapstructure:"clientMsgID"     json:"clientMsgID"     validate:"required"`
	RevokerNickname string `mapstructure:"revokerNickname" json:"revokerNickname"`
	SessionType     int32  `mapstructure:"sessionType"     json:"sessionType"     validate:"required"`
	Seq             uint32 `mapstructure:"seq"             json:"seq"             validate:"required"`
}

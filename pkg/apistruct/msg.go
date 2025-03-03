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

import "github.com/openimsdk/protocol/sdkws"

type PictureBaseInfo struct {
	UUID   string `mapstructure:"uuid"`
	Type   string `mapstructure:"type"   validate:"required"`
	Size   int64  `mapstructure:"size"`
	Width  int32  `mapstructure:"width"  validate:"required"`
	Height int32  `mapstructure:"height" validate:"required"`
	Url    string `mapstructure:"url"    validate:"required"`
}

type PictureElem struct {
	SourcePath      string          `mapstructure:"sourcePath"`
	SourcePicture   PictureBaseInfo `mapstructure:"sourcePicture"   validate:"required"`
	BigPicture      PictureBaseInfo `mapstructure:"bigPicture"      validate:"required"`
	SnapshotPicture PictureBaseInfo `mapstructure:"snapshotPicture" validate:"required"`
}

type SoundElem struct {
	UUID      string `mapstructure:"uuid"`
	SoundPath string `mapstructure:"soundPath"`
	SourceURL string `mapstructure:"sourceUrl" validate:"required"`
	DataSize  int64  `mapstructure:"dataSize"`
	Duration  int64  `mapstructure:"duration"  validate:"required,min=1"`
}

type VideoElem struct {
	VideoPath      string `mapstructure:"videoPath"`
	VideoUUID      string `mapstructure:"videoUUID"`
	VideoURL       string `mapstructure:"videoUrl"       validate:"required"`
	VideoType      string `mapstructure:"videoType"      validate:"required"`
	VideoSize      int64  `mapstructure:"videoSize"      validate:"required"`
	Duration       int64  `mapstructure:"duration"       validate:"required"`
	SnapshotPath   string `mapstructure:"snapshotPath"`
	SnapshotUUID   string `mapstructure:"snapshotUUID"`
	SnapshotSize   int64  `mapstructure:"snapshotSize"`
	SnapshotURL    string `mapstructure:"snapshotUrl"    validate:"required"`
	SnapshotWidth  int32  `mapstructure:"snapshotWidth"  validate:"required"`
	SnapshotHeight int32  `mapstructure:"snapshotHeight" validate:"required"`
}

type FileElem struct {
	FilePath  string `mapstructure:"filePath"`
	UUID      string `mapstructure:"uuid"`
	SourceURL string `mapstructure:"sourceUrl" validate:"required"`
	FileName  string `mapstructure:"fileName"  validate:"required"`
	FileSize  int64  `mapstructure:"fileSize"  validate:"required"`
}
type AtElem struct {
	Text       string   `mapstructure:"text"`
	AtUserList []string `mapstructure:"atUserList" validate:"required,max=1000"`
	IsAtSelf   bool     `mapstructure:"isAtSelf"`
}
type LocationElem struct {
	Description string  `mapstructure:"description"`
	Longitude   float64 `mapstructure:"longitude"   validate:"required"`
	Latitude    float64 `mapstructure:"latitude"    validate:"required"`
}

type CustomElem struct {
	Data        string `mapstructure:"data"        validate:"required"`
	Description string `mapstructure:"description"`
	Extension   string `mapstructure:"extension"`
}

type TextElem struct {
	Content string `json:"content" validate:"required"`
}

type MarkdownTextElem struct {
	Content string `mapstructure:"content" validate:"required"`
}

type StreamMsgElem struct {
	Type    string `mapstructure:"type" validate:"required"`
	Content string `mapstructure:"content" validate:"required"`
}

type RevokeElem struct {
	RevokeMsgClientID string `mapstructure:"revokeMsgClientID" validate:"required"`
}

type QuoteElem struct {
	Text         string     `json:"text,omitempty"`
	QuoteMessage *MsgStruct `json:"quoteMessage,omitempty"`
}

type OANotificationElem struct {
	NotificationName    string       `mapstructure:"notificationName"    json:"notificationName"    validate:"required"`
	NotificationFaceURL string       `mapstructure:"notificationFaceURL" json:"notificationFaceURL"`
	NotificationType    int32        `mapstructure:"notificationType"    json:"notificationType"    validate:"required"`
	Text                string       `mapstructure:"text"                json:"text"                validate:"required"`
	Url                 string       `mapstructure:"url"                 json:"url"`
	MixType             int32        `mapstructure:"mixType"             json:"mixType"             validate:"gte=0,lte=5"`
	PictureElem         *PictureElem `mapstructure:"pictureElem"         json:"pictureElem"`
	SoundElem           *SoundElem   `mapstructure:"soundElem"           json:"soundElem"`
	VideoElem           *VideoElem   `mapstructure:"videoElem"           json:"videoElem"`
	FileElem            *FileElem    `mapstructure:"fileElem"            json:"fileElem"`
	Ex                  string       `mapstructure:"ex"                  json:"ex"`
}

type MessageRevoked struct {
	RevokerID       string `mapstructure:"revokerID"       json:"revokerID"       validate:"required"`
	RevokerRole     int32  `mapstructure:"revokerRole"     json:"revokerRole"     validate:"required"`
	ClientMsgID     string `mapstructure:"clientMsgID"     json:"clientMsgID"     validate:"required"`
	RevokerNickname string `mapstructure:"revokerNickname" json:"revokerNickname"`
	SessionType     int32  `mapstructure:"sessionType"     json:"sessionType"     validate:"required"`
	Seq             uint32 `mapstructure:"seq"             json:"seq"             validate:"required"`
}

type MsgStruct struct {
	ClientMsgID          string                 `json:"clientMsgID,omitempty"`
	ServerMsgID          string                 `json:"serverMsgID,omitempty"`
	CreateTime           int64                  `json:"createTime"`
	SendTime             int64                  `json:"sendTime"`
	SessionType          int32                  `json:"sessionType"`
	SendID               string                 `json:"sendID,omitempty"`
	RecvID               string                 `json:"recvID,omitempty"`
	MsgFrom              int32                  `json:"msgFrom"`
	ContentType          int32                  `json:"contentType"`
	SenderPlatformID     int32                  `json:"senderPlatformID"`
	SenderNickname       string                 `json:"senderNickname,omitempty"`
	SenderFaceURL        string                 `json:"senderFaceUrl,omitempty"`
	GroupID              string                 `json:"groupID,omitempty"`
	Content              string                 `json:"content,omitempty"`
	Seq                  int64                  `json:"seq"`
	IsRead               bool                   `json:"isRead"`
	Status               int32                  `json:"status"`
	IsReact              bool                   `json:"isReact,omitempty"`
	IsExternalExtensions bool                   `json:"isExternalExtensions,omitempty"`
	OfflinePush          *sdkws.OfflinePushInfo `json:"offlinePush,omitempty"`
	AttachedInfo         string                 `json:"attachedInfo,omitempty"`
	Ex                   string                 `json:"ex,omitempty"`
	LocalEx              string                 `json:"localEx,omitempty"`
	TextElem             *TextElem              `json:"textElem,omitempty"`
	PictureElem          *PictureElem           `json:"pictureElem,omitempty"`
	SoundElem            *SoundElem             `json:"soundElem,omitempty"`
	VideoElem            *VideoElem             `json:"videoElem,omitempty"`
	FileElem             *FileElem              `json:"fileElem,omitempty"`
	AtTextElem           *AtElem                `json:"atTextElem,omitempty"`
	LocationElem         *LocationElem          `json:"locationElem,omitempty"`
	CustomElem           *CustomElem            `json:"customElem,omitempty"`
	QuoteElem            *QuoteElem             `json:"quoteElem,omitempty"`
}

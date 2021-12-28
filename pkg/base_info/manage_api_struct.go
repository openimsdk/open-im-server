package base_info

import open_im_sdk "Open_IM/pkg/proto/sdk_ws"

type paramsManagementSendMsg struct {
	OperationID      string                       `json:"operationID" binding:"required"`
	SendID           string                       `json:"sendID" binding:"required"`
	RecvID           string                       `json:"recvID" `
	GroupID          string                       `json:"groupID" `
	SenderNickName   string                       `json:"senderNickName" `
	SenderFaceURL    string                       `json:"senderFaceURL" `
	SenderPlatformID int32                        `json:"senderPlatformID"`
	ForceList        []string                     `json:"forceList" `
	Content          map[string]interface{}       `json:"content" binding:"required"`
	ContentType      int32                        `json:"contentType" binding:"required"`
	SessionType      int32                        `json:"sessionType" binding:"required"`
	IsOnlineOnly     bool                         `json:"isOnlineOnly"`
	OfflinePushInfo  *open_im_sdk.OfflinePushInfo `json:"offlinePushInfo"`
}

type PictureBaseInfo struct {
	UUID   string `mapstructure:"uuid"`
	Type   string `mapstructure:"type" validate:"required"`
	Size   int64  `mapstructure:"size" validate:"required"`
	Width  int32  `mapstructure:"width" validate:"required"`
	Height int32  `mapstructure:"height" validate:"required"`
	Url    string `mapstructure:"url" validate:"required"`
}

type PictureElem struct {
	SourcePath      string          `mapstructure:"sourcePath"`
	SourcePicture   PictureBaseInfo `mapstructure:"sourcePicture" validate:"required"`
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

type paramsDeleteUsers struct {
	OperationID   string   `json:"operationID" binding:"required"`
	DeleteUidList []string `json:"deleteUidList" binding:"required"`
}
type paramsGetAllUsersUid struct {
	OperationID string `json:"operationID" binding:"required"`
}
type paramsGetUsersOnlineStatus struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
}
type paramsAccountCheck struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=100"`
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

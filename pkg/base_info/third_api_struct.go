package base_info

import "mime/multipart"

type MinioStorageCredentialReq struct {
	OperationID string `json:"operationID"`
}

type MiniostorageCredentialResp struct {
	SecretAccessKey string `json:"secretAccessKey"`
	AccessKeyID     string `json:"accessKeyID"`
	SessionToken    string `json:"sessionToken"`
	BucketName      string `json:"bucketName"`
	StsEndpointURL  string `json:"stsEndpointURL"`
}

type MinioUploadFileReq struct {
	OperationID string `form:"operationID" binding:"required"`
	FileType    int    `form:"fileType" binding:"required"`
}

type MinioUploadFileResp struct {
	URL             string `json:"URL"`
	NewName         string `json:"newName"`
	SnapshotURL     string `json:"snapshotURL,omitempty"`
	SnapshotNewName string `json:"snapshotName,omitempty"`
}

type UploadUpdateAppReq struct {
	OperationID string                `form:"operationID" binding:"required"`
	Type        int                   `form:"type" binding:"required"`
	Version     string                `form:"version"  binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
	Yaml        *multipart.FileHeader `form:"yaml"`
	ForceUpdate bool                  `form:"forceUpdate"`
	UpdateLog   string                `form:"updateLog" binding:"required"`
}

type UploadUpdateAppResp struct {
	CommResp
}

type GetDownloadURLReq struct {
	OperationID string `json:"operationID" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	Version     string `json:"version" binding:"required"`
}

type GetDownloadURLResp struct {
	CommResp
	Data struct {
		HasNewVersion bool   `json:"hasNewVersion"`
		ForceUpdate   bool   `json:"forceUpdate"`
		FileURL       string `json:"fileURL"`
		YamlURL       string `json:"yamlURL"`
		Version       string `json:"version"`
		UpdateLog     string `json:"update_log"`
	} `json:"data"`
}

type GetRTCInvitationInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ClientMsgID string `json:"clientMsgID" binding:"required"`
}

type GetRTCInvitationInfoResp struct {
	CommResp
	Data struct {
		OpUserID   string `json:"opUserID"`
		Invitation struct {
			InviterUserID     string   `json:"inviterUserID"`
			InviteeUserIDList []string `json:"inviteeUserIDList"`
			GroupID           string   `json:"groupID"`
			RoomID            string   `json:"roomID"`
			Timeout           int32    `json:"timeout"`
			MediaType         string   `json:"mediaType"`
			SessionType       int32    `json:"sessionType"`
		} `json:"invitation"`
		OfflinePushInfo struct{} `json:"offlinePushInfo"`
	} `json:"data"`
}

type GetRTCInvitationInfoStartAppReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetRTCInvitationInfoStartAppResp struct {
	GetRTCInvitationInfoResp
}

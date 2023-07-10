package cont

import "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3"

type InitiateUploadResult struct {
	UploadID string             `json:"uploadID"` // 上传ID
	PartSize int64              `json:"partSize"` // 分片大小
	Sign     *s3.AuthSignResult `json:"sign"`     // 分片信息
}

type UploadResult struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
	Key  string `json:"key"`
}

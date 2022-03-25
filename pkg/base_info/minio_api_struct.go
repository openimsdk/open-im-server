package base_info

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

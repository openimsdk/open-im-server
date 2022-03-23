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
	OperationID string `json:"operationID"`
	FileType    int    `json:"fileType"`
}

type MinioUploadFileResp struct {
	URL             string `json:"URL"`
	NewName         string `json:"newName"`
	SnapshotURL     string `json:"snapshotURL" binding:"omitempty"`
	SnapshotNewName string `json:"snapshotName" binding:"omitempty"`
}

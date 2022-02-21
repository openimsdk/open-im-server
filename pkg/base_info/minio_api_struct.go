package base_info

type MinioStorageCredentialReq struct {
	OperationID string `json:"operationID"`
}

type MiniostorageCredentialResp struct {
	SecretAccessKey string `json:"secretAccessKey"`
	AccessKeyID string `json:"accessKeyID"`
	SessionToken string `json:"sessionToken"`
	BucketName string `json:"bucketName"`
	StsEndpointURL string `json:"stsEndpointURL"`
}

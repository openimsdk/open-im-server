package base_info

type OSSCredentialReq struct {
	OperationID string `json:"operationID"`
	Filename    string `json:"filename"`
	FileType    string `json:"file_type"`
}

type OSSCredentialRespData struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Token           string `json:"token"`
	Bucket          string `json:"bucket"`
	FinalHost       string `json:"final_host"`
}

type OSSCredentialResp struct {
	CommResp
	OssData OSSCredentialRespData  `json:"-"`
	Data    map[string]interface{} `json:"data"`
}

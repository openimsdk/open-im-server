package base_info

type SetClientInitConfigReq struct {
	OperationID     string  `json:"operationID"  binding:"required"`
	DiscoverPageURL *string `json:"discoverPageURL"`
}

type SetClientInitConfigResp struct {
	CommResp
}

type GetClientInitConfigReq struct {
	OperationID string `json:"operationID"  binding:"required"`
}

type GetClientInitConfigResp struct {
	CommResp
	Data struct {
		DiscoverPageURL string `json:"discoverPageURL"`
	} `json:"data"`
}

package apistruct

type GetConfigReq struct {
	ConfigName string `json:"config_name"`
}

type GetConfigListReq struct {
}

type GetConfigListResp struct {
	Environment string   `json:"environment"`
	Version     string   `json:"version"`
	ConfigNames []string `json:"config_names"`
}

package apistruct

type GetConfigReq struct {
	ConfigName string `json:"config_name"`
}

type GetConfigListResp struct {
	Environment string   `json:"environment"`
	Version     string   `json:"version"`
	ConfigNames []string `json:"config_names"`
}

type SetConfigReq struct {
	ConfigName string `json:"config_name"`
	Data       []byte `json:"data"`
}

package apistruct

type GetConfigReq struct {
	ConfigName string `json:"configName"`
}

type GetConfigListResp struct {
	Environment string   `json:"environment"`
	Version     string   `json:"version"`
	ConfigNames []string `json:"configNames"`
}

type SetConfigReq struct {
	ConfigName string `json:"configName"`
	Data       string `json:"data"`
}

type SetConfigsReq struct {
	Configs []SetConfigReq `json:"configs"`
}

type SetEnableConfigManagerReq struct {
	Enable bool `json:"enable"`
}

type GetEnableConfigManagerResp struct {
	Enable bool `json:"enable"`
}

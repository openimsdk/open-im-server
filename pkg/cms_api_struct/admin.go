package cms_api_struct

import (
	apiStruct "Open_IM/pkg/base_info"
)

type AdminLoginRequest struct {
	AdminName string `json:"admin_name" binding:"required"`
	Secret    string `json:"secret" binding:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type UploadUpdateAppReq struct {
	OperationID string `form:"operationID" binding:"required"`
	Type        int    `form:"type" binding:"required"`
	Version     string `form:"version"  binding:"required"`
	//File        *multipart.FileHeader `form:"file" binding:"required"`
	//Yaml        *multipart.FileHeader `form:"yaml" binding:"required"`
	ForceUpdate bool `form:"forceUpdate"  binding:"required"`
}

type UploadUpdateAppResp struct {
	apiStruct.CommResp
}

type GetDownloadURLReq struct {
	OperationID string `json:"operationID" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	Version     string `json:"version" binding:"required"`
}

type GetDownloadURLResp struct {
	apiStruct.CommResp
	Data struct {
		HasNewVersion bool   `json:"hasNewVersion"`
		ForceUpdate   bool   `json:"forceUpdate"`
		FileURL       string `json:"fileURL"`
		YamlURL       string `json:"yamlURL"`
	} `json:"data"`
}

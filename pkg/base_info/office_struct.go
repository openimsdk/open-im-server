package base_info

import (
	pbOffice "Open_IM/pkg/proto/office"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
)

type GetUserTagsReq struct {
	pbOffice.GetUserTagsReq
}

type GetUserTagsResp struct {
	CommResp
	Data struct {
		Tags []*pbOffice.Tag `json:"tags"`
	} `json:"data"`
}

type CreateTagReq struct {
	pbOffice.CreateTagReq
}

type CreateTagResp struct {
	CommResp
}

type DeleteTagReq struct {
	pbOffice.DeleteTagReq
}

type DeleteTagResp struct {
	CommResp
}

type SetTagReq struct {
	pbOffice.SetTagReq
}

type SetTagResp struct {
	CommResp
}

type SendMsg2TagReq struct {
	pbOffice.SendMsg2TagReq
}

type SendMsg2TagResp struct {
	CommResp
}

type GetTagSendLogsReq struct {
	server_api_params.RequestPagination
	UserID      string `json:"userID"`
	OperationID string `json:"operationID"`
}

type GetTagSendLogsResp struct {
	CommResp
	Data struct {
		Logs        []*pbOffice.TagSendLog `json:"logs"`
		CurrentPage int32                  `json:"currentPage"`
		ShowNumber  int32                  `json:"showNumber"`
	} `json:"data"`
}

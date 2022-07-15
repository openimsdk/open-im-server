package base_info

import (
	pbOffice "Open_IM/pkg/proto/office"
)

type GetUserTagsReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetUserTagsResp struct {
	CommResp
	Data struct {
		Tags []*pbOffice.Tag `json:"tags"`
	} `json:"data"`
}

type CreateTagReq struct {
	TagName     string   `json:"tagName" binding:"required"`    // tag名称
	UserIDList  []string `json:"userIDList" binding:"required"` // 用户ID列表
	OperationID string   `json:"operationID" binding:"required"`
}

type CreateTagResp struct {
	CommResp
}

type DeleteTagReq struct {
	TagID       string `json:"tagID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type DeleteTagResp struct {
	CommResp
}

type SetTagReq struct {
	TagID              string   `json:"tagID" binding:"required"`
	NewName            string   `json:"newName"`
	IncreaseUserIDList []string `json:"increaseUserIDList"`
	ReduceUserIDList   []string `json:"reduceUserIDList"`
	OperationID        string   `json:"operationID" binding:"required"`
}

type SetTagResp struct {
	CommResp
}

type SendMsg2TagReq struct {
	TagList   []string `json:"tagList"`
	UserList  []string `json:"userList"`
	GroupList []string `json:"groupList"`

	SenderPlatformID int32  `json:"senderPlatformID" binding:"required"`
	Content          string `json:"content" binding:"required"`
	OperationID      string `json:"operationID" binding:"required"`
}

type SendMsg2TagResp struct {
	CommResp
}

type GetTagSendLogsReq struct {
	PageNumber  int32  `json:"pageNumber" binding:"required"`
	ShowNumber  int32  `json:"showNumber" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type GetTagSendLogsResp struct {
	CommResp
	Data struct {
		Logs        []*pbOffice.TagSendLog `json:"logs"`
		CurrentPage int32                  `json:"currentPage"`
		ShowNumber  int32                  `json:"showNumber"`
	} `json:"data"`
}

type GetUserTagByIDReq struct {
	TagID       string `json:"tagID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type GetUserTagByIDResp struct {
	CommResp
	Data struct {
		Tag *pbOffice.Tag `json:"tag"`
	} `json:"data"`
}

package base_info

type DelMsgReq struct {
	UserID      string   `json:"userID,omitempty" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty" binding:"required"`
	OperationID string   `json:"operationID,omitempty" binding:"required"`
}

type DelMsgResp struct {
	CommResp
}

type CleanUpMsgReq struct {
	UserID      string `json:"userID"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type CleanUpMsgResp struct {
	CommResp
}
type DelSuperGroupMsgReq struct {
	UserID      string   `json:"userID,omitempty" binding:"required"`
	GroupID     string   `json:"groupID,omitempty" binding:"required"`
	SeqList     []uint32 `json:"seqList,omitempty"`
	IsAllDelete bool     `json:"isAllDelete"`
	OperationID string   `json:"operationID,omitempty" binding:"required"`
}

type DelSuperGroupMsgResp struct {
	CommResp
}
type MsgDeleteNotificationElem struct {
	GroupID     string   `json:"groupID"`
	IsAllDelete bool     `json:"isAllDelete"`
	SeqList     []uint32 `json:"seqList"`
}

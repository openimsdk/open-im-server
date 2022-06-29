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

package base_info

type DelMsgReq struct {
	OpUserID    string   `json:"opUserID,omitempty"`
	UserID      string   `json:"userID,omitempty"`
	SeqList     []uint32 `json:"seqList,omitempty"`
	OperationID string   `json:"operationID,omitempty"`
}

type DelMsgResp struct {
	CommResp
}

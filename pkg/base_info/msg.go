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

//UserID               string   `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
//	GroupID              string   `protobuf:"bytes,2,opt,name=groupID" json:"groupID,omitempty"`
//	MinSeq               uint32   `protobuf:"varint,3,opt,name=minSeq" json:"minSeq,omitempty"`
//	OperationID          string   `protobuf:"bytes,4,opt,name=operationID" json:"operationID,omitempty"`
//	OpUserID             string   `protobuf:"bytes,5,opt,name=opUserID" json:"opUserID,omitempty"`
type SetMsgMinSeqReq struct {
	UserID      string `json:"userID"  binding:"required"`
	GroupID     string `json:"groupID"`
	MinSeq      uint32 `json:"minSeq"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}
type SetMsgMinSeqResp struct {
	CommResp
}

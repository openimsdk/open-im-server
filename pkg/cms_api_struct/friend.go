package cms_api_struct

type GetFriendsReq struct {
	OperationID    string `json:"operationID"`
	UserID         string `json:"userID"`
	FriendUserName string `json:"friendUserName"`
	FriendUserID   string `json:"friendUserID"`
	RequestPagination
}

type FriendInfo struct {
	OwnerUserID    string `json:"ownerUserID"`
	Remark         string `json:"remark"`
	CreateTime     uint32 `json:"createTime"`
	UserID         string `json:"userID"`
	Nickname       string `json:"nickName"`
	AddSource      int32  `json:"addSource"`
	OperatorUserID string `json:"operatorUserID"`
}

type GetFriendsResp struct {
	ResponsePagination
	FriendInfoList []*FriendInfo `json:"friendInfoList"`
	FriendNums     int32         `json:"friendNums"`
}

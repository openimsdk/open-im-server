package cms_api_struct

type GroupResponse struct {
	GroupOwnerName         string `json:"GroupOwnerName"`
	GroupOwnerID           string `json:"GroupOwnerID"`
	GroupID                string `json:"groupID"`
	GroupName              string `json:"groupName"`
	Notification           string `json:"notification"`
	Introduction           string `json:"introduction"`
	FaceURL                string `json:"faceURL"`
	OwnerUserID            string `json:"ownerUserID"`
	CreateTime             uint32 `json:"createTime"`
	MemberCount            uint32 `json:"memberCount"`
	Ex                     string `json:"ex"`
	Status                 int32  `json:"status"`
	CreatorUserID          string `json:"creatorUserID"`
	GroupType              int32  `json:"groupType"`
	NeedVerification       int32  `json:"needVerification"`
	LookMemberInfo         int32  `json:"lookMemberInfo"`
	ApplyMemberFriend      int32  `json:"applyMemberFriend"`
	NotificationUpdateTime uint32 `json:"notificationUpdateTime"`
	NotificationUserID     string `json:"notificationUserID"`
}

type GetGroupsRequest struct {
	RequestPagination
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID"`
	GroupName   string `json:"groupName"`
}

type GetGroupsResponse struct {
	Groups    []GroupResponse `json:"groups"`
	GroupNums int             `json:"groupNums"`
	ResponsePagination
}

type GetGroupMembersRequest struct {
	GroupID     string `form:"groupID" binding:"required"`
	UserName    string `form:"userName"`
	OperationID string `json:"operationID" binding:"required"`
	RequestPagination
}

type GroupMemberResponse struct {
	GroupID        string `json:"groupID"`
	UserID         string `json:"userID"`
	RoleLevel      int32  `json:"roleLevel"`
	JoinTime       int32  `json:"joinTime"`
	Nickname       string `json:"nickname"`
	FaceURL        string `json:"faceURL"`
	AppMangerLevel int32  `json:"appMangerLevel"` //if >0
	JoinSource     int32  `json:"joinSource"`
	OperatorUserID string `json:"operatorUserID"`
	Ex             string `json:"ex"`
	MuteEndTime    uint32 `json:"muteEndTime"`
	InviterUserID  string `json:"inviterUserID"`
}

type GetGroupMembersResponse struct {
	GroupMembers []GroupMemberResponse `json:"groupMembers"`
	ResponsePagination
	MemberNums int `json:"memberNums"`
}

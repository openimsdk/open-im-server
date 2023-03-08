package apistruct

//type ParamsCommFriend struct {
//	OperationID string `json:"operationID" binding:"required"`
//	ToUserID    string `json:"toUserID" binding:"required"`
//	FromUserID  string `json:"fromUserID" binding:"required"`
//}
//
//type AddBlacklistReq struct {
//	ParamsCommFriend
//}
//type AddBlacklistResp struct {
//
//}
//
//type ImportFriendReq struct {
//	FriendUserIDList []string `json:"friendUserIDList" binding:"required"`
//	OperationID      string   `json:"operationID" binding:"required"`
//	FromUserID       string   `json:"fromUserID" binding:"required"`
//}
//type UserIDResult struct {
//	UserID string `json:"userID"`
//	Result int32  `json:"result"`
//}
//type ImportFriendResp struct {
//
//	UserIDResultList []UserIDResult `json:"data"`
//}
//
//type AddFriendReq struct {
//	ParamsCommFriend
//	ReqMsg string `json:"reqMsg"`
//}
//type AddFriendResp struct {
//
//}
//
//type AddFriendResponseReq struct {
//	ParamsCommFriend
//	Flag      int32  `json:"flag" binding:"required,oneof=-1 0 1"`
//	HandleMsg string `json:"handleMsg"`
//}
//type AddFriendResponseResp struct {
//
//}
//
//type DeleteFriendReq struct {
//	ParamsCommFriend
//}
//type DeleteFriendResp struct {
//
//}
//
//type GetBlackListReq struct {
//	OperationID string `json:"operationID" binding:"required"`
//	FromUserID  string `json:"fromUserID" binding:"required"`
//}
//type GetBlackListResp struct {
//
//	BlackUserInfoList []*sdkws.PublicUserInfo `json:"-"`
//	Map              []map[string]interface{}      `json:"data" swaggerignore:"true"`
//}
//
////type PublicUserInfo struct {
////	UserID   string `json:"userID"`
////	Nickname string `json:"nickname"`
////	FaceUrl  string `json:"faceUrl"`
////	Gender   int32  `json:"gender"`
////}
//
//type SetFriendRemarkReq struct {
//	ParamsCommFriend
//	Remark string `json:"remark"`
//}
//type SetFriendRemarkResp struct {
//
//}
//
//type RemoveBlacklistReq struct {
//	ParamsCommFriend
//}
//type RemoveBlacklistResp struct {
//
//}
//
//type IsFriendReq struct {
//	ParamsCommFriend
//}
//type Response struct {
//	Friend bool `json:"isFriend"`
//}
//type IsFriendResp struct {
//
//	Response Response `json:"data"`
//}
//
//type GetFriendsInfoReq struct {
//	ParamsCommFriend
//}
//type GetFriendsInfoResp struct {
//
//	FriendInfoList []*sdkws.FriendInfo `json:"-"`
//	Map           []map[string]interface{}  `json:"data" swaggerignore:"true"`
//}
//
//type GetFriendListReq struct {
//	OperationID string `json:"operationID" binding:"required"`
//	FromUserID  string `json:"fromUserID" binding:"required"`
//}
//type GetFriendListResp struct {
//
//	FriendInfoList []*sdkws.FriendInfo `json:"-"`
//	Map           []map[string]interface{}  `json:"data" swaggerignore:"true"`
//}
//
//type GetFriendApplyListReq struct {
//	OperationID string `json:"operationID" binding:"required"`
//	FromUserID  string `json:"fromUserID" binding:"required"`
//}
//type GetFriendApplyListResp struct {
//
//	FriendRequestList []*sdkws.FriendRequest `json:"-"`
//	Map              []map[string]interface{}     `json:"data" swaggerignore:"true"`
//}
//
//type GetSelfApplyListReq struct {
//	OperationID string `json:"operationID" binding:"required"`
//	FromUserID  string `json:"fromUserID" binding:"required"`
//}
//type GetSelfApplyListResp struct {
//
//	FriendRequestList []*sdkws.FriendRequest `json:"-"`
//	Map              []map[string]interface{}     `json:"data" swaggerignore:"true"`
//}

type FriendInfo struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
	FaceURL  string `json:"faceURL"`
	Gender   int32  `json:"gender"`
	Ex       string `json:"ex"`
}

type PublicUserInfo struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
	FaceURL  string `json:"faceURL"`
	Gender   int32  `json:"gender"`
	Ex       string `json:"ex"`
}

type FriendRequest struct {
	FromUserID    string `json:"fromUserID"`
	FromNickname  string `json:"fromNickname"`
	FromFaceURL   string `json:"fromFaceURL"`
	FromGender    int32  `json:"fromGender"`
	ToUserID      string `json:"toUserID"`
	ToNickname    string `json:"toNickname"`
	ToFaceURL     string `json:"toFaceURL"`
	ToGender      int32  `json:"toGender"`
	HandleResult  int32  `json:"handleResult"`
	ReqMsg        string `json:"reqMsg"`
	CreateTime    uint32 `json:"createTime"`
	HandlerUserID string `json:"handlerUserID"`
	HandleMsg     string `json:"handleMsg"`
	HandleTime    uint32 `json:"handleTime"`
	Ex            string `json:"ex"`
}

type AddBlacklistReq struct {
	ToUserID   string `json:"toUserID" binding:"required"`
	FromUserID string `json:"fromUserID" binding:"required"`
}
type AddBlacklistResp struct {
}

type ImportFriendReq struct {
	FriendUserIDList []string `json:"friendUserIDList" binding:"required"`
	FromUserID       string   `json:"fromUserID" binding:"required"`
}

type ImportFriendResp struct {
	//
}

type AddFriendReq struct {
	ToUserID   string `json:"toUserID" binding:"required"`
	FromUserID string `json:"fromUserID" binding:"required"`
	ReqMsg     string `json:"reqMsg"`
}
type AddFriendResp struct {
	//
}

type AddFriendResponseReq struct {
	ToUserID     string `json:"toUserID" binding:"required"`
	FromUserID   string `json:"fromUserID" binding:"required"`
	HandleResult int32  `json:"flag" binding:"required,oneof=-1 0 1"`
	HandleMsg    string `json:"handleMsg"`
}
type AddFriendResponseResp struct {
}

type DeleteFriendReq struct {
	ToUserID   string `json:"toUserID" binding:"required"`
	FromUserID string `json:"fromUserID" binding:"required"`
}
type DeleteFriendResp struct {
}

type GetBlackListReq struct {
	FromUserID string `json:"fromUserID" binding:"required"`
}
type GetBlackListResp struct {
	BlackUserInfoList []PublicUserInfo `json:"blackUserInfoList"`
}

type SetFriendRemarkReq struct {
	ToUserID   string `json:"toUserID" binding:"required"`
	FromUserID string `json:"fromUserID" binding:"required"`
	Remark     string `json:"remark"`
}
type SetFriendRemarkResp struct {
}

type RemoveBlacklistReq struct {
	ToUserID   string `json:"toUserID" binding:"required"`
	FromUserID string `json:"fromUserID" binding:"required"`
}
type RemoveBlacklistResp struct {
}

type IsFriendReq struct {
	ToUserID   string `json:"toUserID" binding:"required"`
	FromUserID string `json:"fromUserID" binding:"required"`
}
type Response struct {
	Friend bool `json:"isFriend"`
}
type IsFriendResp struct {
	Response Response `json:"data"`
}

type GetFriendListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetFriendListResp struct {
	OwnerUserID    string `json:"ownerUserID"`
	Remark         string `json:"remark"`
	CreateTime     uint32 `json:"createTime"`
	AddSource      int32  `json:"addSource"`
	OperatorUserID string `json:"operatorUserID"`
	Ex             string `json:"ex"`
	//FriendUser           *UserInfo // TODO
}

type GetFriendApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}

type GetFriendApplyListResp struct {
	FriendRequestList []FriendRequest `json:"friendRequestList"`
}

type GetSelfApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetSelfApplyListResp struct {
	FriendRequestList []FriendRequest `json:"friendRequestList"`
}

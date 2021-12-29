package base_info

import open_im_sdk "Open_IM/pkg/proto/sdk_ws"

type ParamsCommFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	ToUserID    string `json:"toUserID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}

type AddBlacklistReq struct {
	ParamsCommFriend
}
type AddBlacklistResp struct {
	CommResp
}

type ImportFriendReq struct {
	FriendUserIDList []string `json:"friendUserIDList" binding:"required"`
	OperationID      string   `json:"operationID" binding:"required"`
	FromUserID       string   `json:"fromUserID" binding:"required"`
}
type ImportFriendResp struct {
	CommResp
	Data []string `json:"data"`
}

type AddFriendReq struct {
	ParamsCommFriend
	ReqMsg string `json:"reqMsg"`
}
type AddFriendResp struct {
	CommResp
}

type AddFriendResponseReq struct {
	ParamsCommFriend
	Flag      int32  `json:"flag" binding:"required"`
	HandleMsg string `json:"handleMsg"`
}
type AddFriendResponseResp struct {
	CommResp
}

type DeleteFriendReq struct {
	ParamsCommFriend
}
type DeleteFriendResp struct {
	CommResp
}

type GetBlackListReq struct {
	ParamsCommFriend
}
type GetBlackListResp struct {
	CommResp
	BlackUserInfoList []*BlackUserInfo `json:"data"`
}

//type PublicUserInfo struct {
//	UserID   string `json:"userID"`
//	Nickname string `json:"nickname"`
//	FaceUrl  string `json:"faceUrl"`
//	Gender   int32  `json:"gender"`
//}

type BlackUserInfo struct {
	open_im_sdk.PublicUserInfo
}

type SetFriendCommentReq struct {
	ParamsCommFriend
	Remark string `json:"remark" binding:"required"`
}
type SetFriendCommentResp struct {
	CommResp
}

type RemoveBlackListReq struct {
	ParamsCommFriend
}
type RemoveBlackListResp struct {
	CommResp
}

type IsFriendReq struct {
	ParamsCommFriend
}
type IsFriendResp struct {
	CommResp
	Response bool `json:"response"`
}

type GetFriendsInfoReq struct {
	ParamsCommFriend
}
type GetFriendsInfoResp struct {
	CommResp
	FriendInfoList []*open_im_sdk.FriendInfo `json:"data"`
}

type GetFriendListReq struct {
	ParamsCommFriend
}
type GetFriendListResp struct {
	CommResp
	FriendInfoList []*open_im_sdk.FriendInfo `json:"data"`
}

type GetFriendApplyListReq struct {
	ParamsCommFriend
}
type GetFriendApplyListResp struct {
	CommResp
	FriendRequestList []*open_im_sdk.FriendRequest `json:"data"`
}

type GetSelfApplyListReq struct {
	ParamsCommFriend
}
type GetSelfApplyListResp struct {
	CommResp
	FriendRequestList []*open_im_sdk.FriendRequest `json:"data"`
}

package base_info

import open_im_sdk "Open_IM/pkg/proto/sdk_ws"

type paramsCommFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	ToUserID    string `json:"toUserID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}

type AddBlacklistReq struct {
	paramsCommFriend
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
	paramsCommFriend
	ReqMsg string `json:"reqMsg"`
}
type AddFriendResp struct {
	CommResp
}

type AddFriendResponseReq struct {
	paramsCommFriend
	Flag      int32  `json:"flag" binding:"required"`
	HandleMsg string `json:"handleMsg"`
}
type AddFriendResponseResp struct {
	CommResp
}

type DeleteFriendReq struct {
	paramsCommFriend
}
type DeleteFriendResp struct {
	CommResp
}

type GetBlackListReq struct {
	paramsCommFriend
}
type GetBlackListResp struct {
	CommResp
	BlackUserInfoList []*blackUserInfo `json:"data"`
}

//type PublicUserInfo struct {
//	UserID   string `json:"userID"`
//	Nickname string `json:"nickname"`
//	FaceUrl  string `json:"faceUrl"`
//	Gender   int32  `json:"gender"`
//}

type blackUserInfo struct {
	open_im_sdk.PublicUserInfo
}

type SetFriendCommentReq struct {
	paramsCommFriend
	Remark string `json:"remark" binding:"required"`
}
type SetFriendCommentResp struct {
	CommResp
}

type RemoveBlackListReq struct {
	paramsCommFriend
}
type RemoveBlackListResp struct {
	CommResp
}

type IsFriendReq struct {
	paramsCommFriend
}
type IsFriendResp struct {
	CommResp
	Response bool `json:"response"`
}

type GetFriendsInfoReq struct {
	paramsCommFriend
}
type GetFriendsInfoResp struct {
	CommResp
	FriendInfoList []*open_im_sdk.FriendInfo `json:"data"`
}

type GetFriendListReq struct {
	paramsCommFriend
}
type GetFriendListResp struct {
	CommResp
	FriendInfoList []*open_im_sdk.FriendInfo `json:"data"`
}

type GetFriendApplyListReq struct {
	paramsCommFriend
}
type GetFriendApplyListResp struct {
	CommResp
	FriendRequestList open_im_sdk.FriendRequest `json:"data"`
}

type GetSelfApplyListReq struct {
	paramsCommFriend
}
type GetSelfApplyListResp struct {
	CommResp
	FriendRequestList open_im_sdk.FriendRequest `json:"data"`
}

package base_info

import "Open_IM/pkg/proto/office"

type CreateOneWorkMomentReq struct {
	office.CreateOneWorkMomentReq
}

type CreateOneWorkMomentResp struct {
	CommResp
	Data struct{} `json:"data"`
}

type DeleteOneWorkMomentReq struct {
	WorkMomentID string `json:"workMomentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}

type DeleteOneWorkMomentResp struct {
	CommResp
	Data struct{} `json:"data"`
}

type LikeOneWorkMomentReq struct {
	WorkMomentID string `json:"workMomentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}

type LikeOneWorkMomentResp struct {
	CommResp
	Data struct{} `json:"data"`
}

type CommentOneWorkMomentReq struct {
	WorkMomentID string `json:"workMomentID" binding:"required"`
	ReplyUserID  string `json:"replyUserID"`
	Content      string `json:"content"  binding:"required"`
	OperationID  string `json:"operationID"  binding:"required"`
}

type CommentOneWorkMomentResp struct {
	CommResp
	Data struct{} `json:"data"`
}

type DeleteCommentReq struct {
	WorkMomentID string `json:"workMomentID" binding:"required"`
	ContentID    string `json:"contentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}

type DeleteCommentResp struct {
	CommResp
	Data struct{} `json:"data"`
}

type WorkMomentsUserCommonReq struct {
	PageNumber  int32  `json:"pageNumber" binding:"required"`
	ShowNumber  int32  `json:"showNumber" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type GetWorkMomentByIDReq struct {
	WorkMomentID string `json:"workMomentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}

type WorkMoment struct {
	WorkMomentID       string            `json:"workMomentID"`
	UserID             string            `json:"userID"`
	Content            string            `json:"content"`
	LikeUserList       []*WorkMomentUser `json:"likeUsers"`
	Comments           []*Comment        `json:"comments"`
	FaceURL            string            `json:"faceURL"`
	UserName           string            `json:"userName"`
	AtUserList         []*WorkMomentUser `json:"atUsers"`
	PermissionUserList []*WorkMomentUser `json:"permissionUsers"`
	CreateTime         int32             `json:"createTime"`
	Permission         int32             `json:"permission"`
}

type WorkMomentUser struct {
	UserID   string `json:"userID"`
	UserName string `json:"userName"`
}

type Comment struct {
	UserID        string `json:"userID"`
	UserName      string `json:"userName"`
	ReplyUserID   string `json:"replyUserID"`
	ReplyUserName string `json:"replyUserName"`
	ContentID     string `json:"contentID"`
	Content       string `json:"content"`
	CreateTime    int32  `json:"createTime"`
}

type GetWorkMomentByIDResp struct {
	CommResp
	Data struct {
		WorkMoment *WorkMoment `json:"workMoment"`
	} `json:"data"`
}

type GetUserWorkMomentsReq struct {
	WorkMomentsUserCommonReq
	UserID string `json:"userID"`
}

type GetUserWorkMomentsResp struct {
	CommResp
	Data struct {
		WorkMoments []*WorkMoment `json:"workMoments"`
		CurrentPage int32         `json:"currentPage"`
		ShowNumber  int32         `json:"showNumber"`
	} `json:"data"`
}

type GetUserFriendWorkMomentsReq struct {
	WorkMomentsUserCommonReq
}

type GetUserFriendWorkMomentsResp struct {
	CommResp
	Data struct {
		WorkMoments []*WorkMoment `json:"workMoments"`
		CurrentPage int32         `json:"currentPage"`
		ShowNumber  int32         `json:"showNumber"`
	} `json:"data"`
}

type SetUserWorkMomentsLevelReq struct {
	office.SetUserWorkMomentsLevelReq
}

type SetUserWorkMomentsLevelResp struct {
	CommResp
	Data struct{} `json:"data"`
}

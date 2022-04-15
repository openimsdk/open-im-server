package base_info

import "Open_IM/pkg/proto/office"

type CreateOneWorkMomentReq struct {
	office.CreateOneWorkMomentReq
}

type CreateOneWorkMomentResp struct {
	CommResp
}

type DeleteOneWorkMomentReq struct {
	office.DeleteOneWorkMomentReq
}

type DeleteOneWorkMomentResp struct {
	CommResp
}

type LikeOneWorkMomentReq struct {
	office.LikeOneWorkMomentReq
}

type LikeOneWorkMomentResp struct {
	CommResp
}

type CommentOneWorkMomentReq struct {
	office.CommentOneWorkMomentReq
}

type CommentOneWorkMomentResp struct {
	CommResp
}

type WorkMomentsUserCommonReq struct {
	PageNumber  int32  `json:"pageNumber" binding:"required"`
	ShowNumber  int32  `json:"showNumber" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"UserID" binding:"required"`
}

type GetUserWorkMomentsReq struct {
	WorkMomentsUserCommonReq
}

type GetUserWorkMomentsResp struct {
	CommResp
	Data struct {
		WorkMoments []*office.WorkMoment `json:"workMoments"`
		CurrentPage int32                `json:"currentPage"`
		ShowNumber  int32                `json:"showNumber"`
	} `json:"data"`
}

type GetUserFriendWorkMomentsReq struct {
	WorkMomentsUserCommonReq
}

type GetUserFriendWorkMomentsResp struct {
	CommResp
	Data struct {
		WorkMoments []*office.WorkMoment `json:"workMoments"`
		CurrentPage int32                `json:"currentPage"`
		ShowNumber  int32                `json:"showNumber"`
	} `json:"data"`
}

type GetUserWorkMomentsCommentsMsgReq struct {
	WorkMomentsUserCommonReq
}

type GetUserWorkMomentsCommentsMsgResp struct {
	CommResp
	Data struct {
		CommentsMsg    []*office.CommentsMsg `json:"comments"`
		CurrentPage int32             `json:"currentPage"`
		ShowNumber  int32             `json:"showNumber"`
	} `json:"data"`
}

type SetUserWorkMomentsLevelReq struct {
	office.SetUserWorkMomentsLevelReq
}

type SetUserWorkMomentsLevelResp struct {
	CommResp
}

type ClearUserWorkMomentsCommentsMsgReq struct {
	office.ClearUserWorkMomentsCommentsMsgReq
}

type ClearUserWorkMomentsCommentsMsgResp struct {
	CommResp
}
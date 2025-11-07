package callbackstruct

type CallbackBeforeAddFriendReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
	Ex              string `json:"ex"`
}

type CallbackBeforeAddFriendResp struct {
	CommonCallbackResp
}

type CallBackAddFriendReplyBeforeReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"toUserID"`
}

type CallBackAddFriendReplyBeforeResp struct {
	CommonCallbackResp
}

type CallbackBeforeSetFriendRemarkReq struct {
	CallbackCommand `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	FriendUserID    string `json:"friendUserID"`
	Remark          string `json:"remark"`
}

type CallbackBeforeSetFriendRemarkResp struct {
	CommonCallbackResp
	Remark string `json:"remark"`
}

type CallbackAfterSetFriendRemarkReq struct {
	CallbackCommand `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	FriendUserID    string `json:"friendUserID"`
	Remark          string `json:"remark"`
}

type CallbackAfterSetFriendRemarkResp struct {
	CommonCallbackResp
}
type CallbackAfterAddFriendReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"`
}

type CallbackAfterAddFriendResp struct {
	CommonCallbackResp
}
type CallbackBeforeAddBlackReq struct {
	CallbackCommand `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID" `
	BlackUserID     string `json:"blackUserID"`
}

type CallbackBeforeAddBlackResp struct {
	CommonCallbackResp
}

type CallbackBeforeAddFriendAgreeReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"blackUserID"`
	HandleResult    int32  `json:"HandleResult"`
	HandleMsg       string `json:"HandleMsg"`
}

type CallbackBeforeAddFriendAgreeResp struct {
	CommonCallbackResp
}

type CallbackAfterAddFriendAgreeReq struct {
	CallbackCommand `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"blackUserID"`
	HandleResult    int32  `json:"HandleResult"`
	HandleMsg       string `json:"HandleMsg"`
}

type CallbackAfterAddFriendAgreeResp struct {
	CommonCallbackResp
}

type CallbackAfterDeleteFriendReq struct {
	CallbackCommand `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID" `
	FriendUserID    string `json:"friendUserID"`
}
type CallbackAfterDeleteFriendResp struct {
	CommonCallbackResp
}

type CallbackBeforeImportFriendsReq struct {
	CallbackCommand `json:"callbackCommand"`
	OwnerUserID     string   `json:"ownerUserID" `
	FriendUserIDs   []string `json:"friendUserIDs"`
}
type CallbackBeforeImportFriendsResp struct {
	CommonCallbackResp
	FriendUserIDs []string `json:"friendUserIDs"`
}
type CallbackAfterImportFriendsReq struct {
	CallbackCommand `json:"callbackCommand"`
	OwnerUserID     string   `json:"ownerUserID" `
	FriendUserIDs   []string `json:"friendUserIDs"`
}
type CallbackAfterImportFriendsResp struct {
	CommonCallbackResp
}

type CallbackAfterRemoveBlackReq struct {
	CallbackCommand `json:"callbackCommand"`
	OwnerUserID     string `json:"ownerUserID"`
	BlackUserID     string `json:"blackUserID"`
}
type CallbackAfterRemoveBlackResp struct {
	CommonCallbackResp
}

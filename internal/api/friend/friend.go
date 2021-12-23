package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsCommFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	ToUserID    string `json:"toUserID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}

func AddBlacklist(c *gin.Context) {
	params := paramsCommFriend{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		log.NewError("0", "BindJSON failed ", err.Error())
		return
	}
	req := &pbFriend.AddBlacklistReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "AddBlacklist args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.AddBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddBlacklist failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add blacklist rpc server failed"})
		return
	}

	resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.CommID.OperationID, "AddBlacklist api return ", resp)
}

type paramsImportFriendReq struct {
	FriendUserIDList []string `json:"friendUserIDList" binding:"required"`
	OperationID      string   `json:"operationID" binding:"required"`
	Token            string   `json:"token"`
	FromUserID       string   `json:"fromUserID" binding:"required"`
	OpUserID         string   `json:"opUserID" binding:"required"`
}

func ImportFriend(c *gin.Context) {
	params := paramsImportFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		log.NewError("0", "BindJSON failed ", err.Error())
		return
	}

	req := &pbFriend.ImportFriendReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}

	log.NewInfo(req.OperationID, "ImportFriend args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.ImportFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "ImportFriend failed", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "cImportFriend failed " + err.Error()})
		return
	}

	failedUidList := make([]string, 0)
	for _, v := range RpcResp.FailedUidList {
		failedUidList = append(failedUidList, v)
	}
	resp := gin.H{"errCode": RpcResp.CommonResp.ErrCode, "errMsg": RpcResp.CommonResp.ErrMsg, "failedUidList": failedUidList}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.OperationID, "AddBlacklist api return ", resp)
}

type paramsAddFriend struct {
	paramsCommFriend
	ReqMessage string `json:"reqMessage"`
}

func AddFriend(c *gin.Context) {
	params := paramsAddFriend{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.AddFriendReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	req.ReqMessage = params.ReqMessage
	log.NewInfo("AddFriend args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.AddFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call AddFriend rpc server failed"})
		return
	}

	resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.CommID.OperationID, "AddFriend api return ", resp)
}

type paramsAddFriendResponse struct {
	paramsCommFriend
	Flag int32 `json:"flag" binding:"required"`
}

func AddFriendResponse(c *gin.Context) {
	params := paramsAddFriendResponse{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.AddFriendResponseReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	req.Flag = params.Flag

	log.NewInfo(req.CommID.OperationID, "AddFriendResponse args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.AddFriendResponse(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddFriendResponse failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add_friend_response rpc server failed"})
		return
	}

	resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.CommID.OperationID, "AddFriendResponse api return ", resp)
}

type paramsDeleteFriend struct {
	paramsCommFriend
}

func DeleteFriend(c *gin.Context) {
	params := paramsDeleteFriend{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.DeleteFriendReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "DeleteFriend args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.DeleteFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "DeleteFriend failed ", err, req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call delete_friend rpc server failed"})
		return
	}

	resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.CommID.OperationID, "AddFriendResponse api return ", resp)
}

type paramsGetBlackList struct {
	paramsCommFriend
}

type PublicUserInfo struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
	FaceUrl  string `json:"faceUrl"`
	Gender   int32  `json:"gender"`
}

type blackUserInfo struct {
	PublicUserInfo
}

func GetBlacklist(c *gin.Context) {
	params := paramsGetBlackList{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetBlacklistReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}

	log.NewInfo(req.CommID.OperationID, "GetBlacklist args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.GetBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetBlacklist failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get blacklist rpc server failed"})
		return
	}

	if RpcResp.ErrCode == 0 {
		userBlackList := make([]blackUserInfo, 0)
		for _, friend := range RpcResp.Data {
			var b blackUserInfo
			utils.CopyStructFields(&b, friend)

			userBlackList = append(userBlackList, b)
		}
		resp := gin.H{
			"errCode": RpcResp.ErrCode,
			"errMsg":  RpcResp.ErrMsg,
			"data":    userBlackList,
		}
		c.JSON(http.StatusOK, resp)
		log.NewInfo(req.CommID.OperationID, "GetBlacklist api return ", resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
		c.JSON(http.StatusOK, resp)
		log.NewError(req.CommID.OperationID, "GetBlacklist api return ", resp)
	}
}

type paramsSetFriendComment struct {
	paramsCommFriend
	remark string `json:"remark" binding:"required"`
}

func SetFriendComment(c *gin.Context) {
	params := paramsSetFriendComment{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.SetFriendCommentReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	req.Remark = params.remark

	log.NewInfo(req.CommID.OperationID, "SetFriendComment args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.SetFriendComment(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "SetFriendComment failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call set friend comment rpc server failed"})
		return
	}
	resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.CommID.OperationID, "SetFriendComment api return ", resp)
}

type paramsRemoveBlackList struct {
	paramsCommFriend
}

func RemoveBlacklist(c *gin.Context) {
	params := paramsRemoveBlackList{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.RemoveBlacklistReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}

	log.NewInfo(req.CommID.OperationID, "RemoveBlacklist args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.RemoveBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "RemoveBlacklist failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call remove blacklist rpc server failed"})
		return
	}

	resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.CommID.OperationID, "SetFriendComment api return ", resp)
}

type paramsIsFriend struct {
	paramsCommFriend
}

func IsFriend(c *gin.Context) {
	params := paramsIsFriend{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.IsFriendReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "IsFriend args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.IsFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "IsFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add friend rpc server failed"})
		return
	}
	resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg, "isFriend": RpcResp.ShipType}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
}

type paramsSearchFriend struct {
	paramsCommFriend
}

func GetFriendsInfo(c *gin.Context) {
	params := paramsSearchFriend{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetFriendsInfoReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "GetFriendsInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.GetFriendsInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendsInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call search friend rpc server failed"})
		return
	}

	if RpcResp.ErrCode == 0 {
		var fi friendInfo
		utils.CopyStructFields(&fi, RpcResp.Data.FriendUser)
		utils.CopyStructFields(&fi, RpcResp.Data)

		resp := gin.H{
			"errCode": RpcResp.ErrCode,
			"errMsg":  RpcResp.ErrMsg,
			"data":    fi,
		}
		log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
		c.JSON(http.StatusOK, resp)
	} else {
		resp := gin.H{
			"errCode": RpcResp.ErrCode,
			"errMsg":  RpcResp.ErrMsg,
		}
		log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
		c.JSON(http.StatusOK, resp)
	}
}

type paramsGetFriendList struct {
	paramsCommFriend
}

type friendInfo struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
	FaceUrl  string `json:"faceUrl"`
	Gender   int32  `json:"gender"`
	Mobile   string `json:"mobile"`
	Birth    string `json:"birth"`
	Email    string `json:"email"`
	Ext      string `json:"ext"`
	Remark   string `json:"remark"`
	IsBlack  int32  `json:"isBlack"`
}

func GetFriendList(c *gin.Context) {
	params := paramsGetFriendList{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetFriendListReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "GetFriendList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.GetFriendList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get friend list rpc server failed"})
		return
	}

	if RpcResp.ErrCode == 0 {
		friendsInfo := make([]friendInfo, 0)
		for _, friend := range RpcResp.Data {

			var fi friendInfo
			utils.CopyStructFields(&fi, friend.FriendUser)
			utils.CopyStructFields(&fi, RpcResp.Data)
			friendsInfo = append(friendsInfo, fi)
		}
		resp := gin.H{
			"errCode": RpcResp.ErrCode,
			"errMsg":  RpcResp.ErrMsg,
			"data":    friendsInfo,
		}
		c.JSON(http.StatusOK, resp)
		log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
		c.JSON(http.StatusOK, resp)
		log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	}

}

type paramsGetApplyList struct {
	paramsCommFriend
}

type FriendApplicationUserInfo struct {
	PublicUserInfo
	ApplyTime  int64  `json:"applyTime"`
	ReqMessage string `json:"reqMessage`
	Flag       int32  `json:"flag"`
}

func GetFriendApplyList(c *gin.Context) {
	params := paramsGetApplyList{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetFriendApplyReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "GetFriendApplyList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)

	RpcResp, err := client.GetFriendApplyList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendApplyList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get friend apply list rpc server failed"})
		return
	}

	if RpcResp.ErrCode == 0 {
		userInfoList := make([]FriendApplicationUserInfo, 0)
		for _, applyUserinfo := range RpcResp.Data {
			var un FriendApplicationUserInfo
			utils.CopyStructFields(&un, applyUserinfo.UserInfo)
			utils.CopyStructFields(&un, applyUserinfo)
			userInfoList = append(userInfoList, un)
		}
		resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg, "data": userInfoList}
		c.JSON(http.StatusOK, resp)
		log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
		c.JSON(http.StatusOK, resp)
		log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	}
}

func GetSelfApplyList(c *gin.Context) {
	params := paramsGetApplyList{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetFriendApplyReq{}
	utils.CopyStructFields(req.CommID, params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "GetSelfApplyList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	RpcResp, err := client.GetSelfApplyList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetSelfApplyList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get self apply list rpc server failed"})
		return
	}

	if RpcResp.ErrCode == 0 {
		userInfoList := make([]FriendApplicationUserInfo, 0)
		for _, applyUserinfo := range RpcResp.Data {
			var un FriendApplicationUserInfo
			utils.CopyStructFields(&un, applyUserinfo.UserInfo)
			utils.CopyStructFields(&un, applyUserinfo)
			userInfoList = append(userInfoList, un)
		}
		resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg, "data": userInfoList}
		c.JSON(http.StatusOK, resp)
		log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrCode, "errMsg": RpcResp.ErrMsg}
		c.JSON(http.StatusOK, resp)
		log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	}
}

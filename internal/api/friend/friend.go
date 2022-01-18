package friend

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/friend"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AddBlack(c *gin.Context) {
	params := api.AddBlacklistReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AddBlacklistReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "AddBlacklist args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.AddBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddBlacklist failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add blacklist rpc server failed"})
		return
	}
	resp := api.AddBlacklistResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "AddBlacklist api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func ImportFriend(c *gin.Context) {
	params := api.ImportFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.ImportFriendReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "ImportFriend args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.ImportFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "ImportFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "ImportFriend failed "})
		return
	}
	resp := api.ImportFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	if resp.ErrCode == 0 {
		for _, v := range RpcResp.UserIDResultList {
			resp.UserIDResultList = append(resp.UserIDResultList, api.UserIDResult{UserID: v.UserID, Result: v.Result})
		}
	}
	if len(resp.UserIDResultList) == 0 {
		resp.UserIDResultList = []api.UserIDResult{}
	}
	log.NewInfo(req.OperationID, "ImportFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func AddFriend(c *gin.Context) {
	params := api.AddFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AddFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.ReqMsg = params.ReqMsg
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "AddFriend args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.AddFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call AddFriend rpc server failed"})
		return
	}

	resp := api.AddFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "AddFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func AddFriendResponse(c *gin.Context) {
	params := api.AddFriendResponseReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AddFriendResponseReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.HandleMsg = params.HandleMsg
	req.HandleResult = params.Flag
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	utils.CopyStructFields(req, &params)
	log.NewInfo(req.CommID.OperationID, "AddFriendResponse args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.AddFriendResponse(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddFriendResponse failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add_friend_response rpc server failed"})
		return
	}

	resp := api.AddFriendResponseResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "AddFriendResponse api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func DeleteFriend(c *gin.Context) {
	params := api.DeleteFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.DeleteFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "DeleteFriend args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.DeleteFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "DeleteFriend failed ", err, req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call delete_friend rpc server failed"})
		return
	}

	resp := api.DeleteFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "DeleteFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetBlacklist(c *gin.Context) {
	params := api.GetBlackListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetBlacklistReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "GetBlacklist args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetBlacklist failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get blacklist rpc server failed"})
		return
	}

	resp := api.GetBlackListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	for _, v := range RpcResp.BlackUserInfoList {
		black := open_im_sdk.PublicUserInfo{}
		utils.CopyStructFields(&black, v)
		resp.BlackUserInfoList = append(resp.BlackUserInfoList, &black)
	}
	resp.Data = jsonData.JsonDataList(resp.BlackUserInfoList)
	log.NewInfo(req.CommID.OperationID, "GetBlacklist api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func SetFriendRemark(c *gin.Context) {
	params := api.SetFriendRemarkReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SetFriendRemarkReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.Remark = params.Remark
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "SetFriendComment args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.SetFriendRemark(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "SetFriendComment failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call set friend comment rpc server failed"})
		return
	}
	resp := api.SetFriendRemarkResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}

	log.NewInfo(req.CommID.OperationID, "SetFriendComment api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func RemoveBlack(c *gin.Context) {
	params := api.RemoveBlackListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.RemoveBlacklistReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}

	log.NewInfo(req.CommID.OperationID, "RemoveBlacklist args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.RemoveBlacklist(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "RemoveBlacklist failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call remove blacklist rpc server failed"})
		return
	}
	resp := api.RemoveBlackListResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "RemoveBlacklist api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func IsFriend(c *gin.Context) {
	params := api.IsFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.IsFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "IsFriend args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.IsFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "IsFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add friend rpc server failed"})
		return
	}
	resp := api.IsFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	resp.Response.Friend = RpcResp.Response

	log.NewInfo(req.CommID.OperationID, "IsFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}

//
//func GetFriendsInfo(c *gin.Context) {
//	params := api.GetFriendsInfoReq{}
//	if err := c.BindJSON(&params); err != nil {
//		log.NewError("0", "BindJSON failed ", err.Error())
//		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
//		return
//	}
//	req := &rpc.GetFriendsInfoReq{}
//	utils.CopyStructFields(req.CommID, params)
//	var ok bool
//	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
//	if !ok {
//		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
//		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
//		return
//	}
//	log.NewInfo(req.CommID.OperationID, "GetFriendsInfo args ", req.String())
//
//	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
//	client := rpc.NewFriendClient(etcdConn)
//	RpcResp, err := client.GetFriendsInfo(context.Background(), req)
//	if err != nil {
//		log.NewError(req.CommID.OperationID, "GetFriendsInfo failed ", err.Error(), req.String())
//		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call search friend rpc server failed"})
//		return
//	}
//
//	resp := api.GetFriendsInfoResp{CommResp:api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
//	utils.CopyStructFields(&resp, RpcResp)
//	c.JSON(http.StatusOK, resp)
//	log.NewInfo(req.CommID.OperationID, "GetFriendsInfo api return ", resp)
//}

func GetFriendList(c *gin.Context) {
	params := api.GetFriendListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetFriendListReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "GetFriendList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetFriendList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get friend list rpc server failed"})
		return
	}

	resp := api.GetFriendListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, FriendInfoList: RpcResp.FriendInfoList}
	resp.Data = jsonData.JsonDataList(resp.FriendInfoList)
	log.NewInfo(req.CommID.OperationID, "GetFriendList api return ", resp)
	c.JSON(http.StatusOK, resp)
	//c.JSON(http.StatusOK, resp)
}

func GetFriendApplyList(c *gin.Context) {
	params := api.GetFriendApplyListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetFriendApplyListReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "GetFriendApplyList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)

	RpcResp, err := client.GetFriendApplyList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendApplyList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get friend apply list rpc server failed"})
		return
	}

	resp := api.GetFriendApplyListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, FriendRequestList: RpcResp.FriendRequestList}
	resp.Data = jsonData.JsonDataList(resp.FriendRequestList)
	log.NewInfo(req.CommID.OperationID, "GetFriendApplyList api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetSelfFriendApplyList(c *gin.Context) {
	params := api.GetSelfApplyListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetSelfApplyListReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params)
	var ok bool
	ok, req.CommID.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.CommID.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.CommID.OperationID, "GetSelfApplyList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.GetSelfApplyList(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetSelfApplyList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get self apply list rpc server failed"})
		return
	}
	resp := api.GetSelfApplyListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, FriendRequestList: RpcResp.FriendRequestList}
	resp.Data = jsonData.JsonDataList(resp.FriendRequestList)
	log.NewInfo(req.CommID.OperationID, "GetSelfApplyList api return ", resp)
	c.JSON(http.StatusOK, resp)
}

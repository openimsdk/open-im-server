package group

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetJoinedSuperGroupList(c *gin.Context) {
	req := api.GetJoinedSuperGroupListReq{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, opUserID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb := rpc.GetJoinedSuperGroupListReq{OperationID: req.OperationID, OpUserID: opUserID, UserID: req.FromUserID}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	rpcResp, err := client.GetJoinedSuperGroupList(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, "InviteUserToGroup failed ", err.Error(), reqPb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	GroupListResp := api.GetJoinedSuperGroupListResp{GetJoinedGroupListResp: api.GetJoinedGroupListResp{CommResp: api.CommResp{ErrCode: rpcResp.CommonResp.ErrCode, ErrMsg: rpcResp.CommonResp.ErrMsg}, GroupInfoList: rpcResp.GroupList}}
	GroupListResp.Data = jsonData.JsonDataList(GroupListResp.GroupInfoList)
	log.NewInfo(req.OperationID, "GetJoinedSuperGroupList api return ", GroupListResp)
	c.JSON(http.StatusOK, GroupListResp)
}

func GetSuperGroupsInfo(c *gin.Context) {
	req := api.GetSuperGroupsInfoReq{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, opUserID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb := rpc.GetSuperGroupsInfoReq{OperationID: req.OperationID, OpUserID: opUserID, GroupIDList: req.GroupIDList}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	rpcResp, err := client.GetSuperGroupsInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, "InviteUserToGroup failed ", err.Error(), reqPb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.GetSuperGroupsInfoResp{GetGroupInfoResp: api.GetGroupInfoResp{CommResp: api.CommResp{ErrCode: rpcResp.CommonResp.ErrCode, ErrMsg: rpcResp.CommonResp.ErrMsg}, GroupInfoList: rpcResp.GroupInfoList}}
	resp.Data = jsonData.JsonDataList(resp.GroupInfoList)
	log.NewInfo(req.OperationID, "GetGroupsInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

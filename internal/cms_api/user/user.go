package user

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdminCms "Open_IM/pkg/proto/admin_cms"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	pb "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
	var (
		req   cms_api_struct.AddUserRequest
		reqPb pb.AddUserReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.AddUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg})
}

func BlockUser(c *gin.Context) {
	var (
		req   cms_api_struct.BlockUserRequest
		reqPb pb.BlockUserReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.BlockUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg})
}

func UnblockUser(c *gin.Context) {
	var (
		req   cms_api_struct.UnblockUserRequest
		resp  cms_api_struct.UnBlockUserResponse
		reqPb pb.UnBlockUserReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.UnBlockUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), err.Error(), reqPb.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg})
}

func GetBlockUsers(c *gin.Context) {
	var (
		req    cms_api_struct.GetBlockUsersRequest
		resp   cms_api_struct.GetBlockUsersResponse
		reqPb  pb.GetBlockUsersReq
		respPb *pb.GetBlockUsersResp
	)
	reqPb.Pagination = &commonPb.RequestPagination{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	utils.CopyStructFields(&reqPb.Pagination, &req)
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "blockUsers", reqPb.Pagination, req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetBlockUsers(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetBlockUsers rpc", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	for _, v := range respPb.BlockUsers {
		resp.BlockUsers = append(resp.BlockUsers, cms_api_struct.BlockUser{
			UserResponse: cms_api_struct.UserResponse{
				UserID:      v.UserInfo.UserID,
				FaceURL:     v.UserInfo.FaceURL,
				Nickname:    v.UserInfo.Nickname,
				PhoneNumber: v.UserInfo.PhoneNumber,
				Email:       v.UserInfo.Email,
				Gender:      int(v.UserInfo.Gender),
			},
			BeginDisableTime: v.BeginDisableTime,
			EndDisableTime:   v.EndDisableTime,
		})
	}
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.UserNums = respPb.UserNums
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp})
}

func GetUserIDByEmailAndPhoneNumber(c *gin.Context) {
	var (
		req    cms_api_struct.GetUserIDByEmailAndPhoneNumberRequest
		resp   cms_api_struct.GetUserIDByEmailAndPhoneNumberResponse
		reqPb  pbAdminCms.GetUserIDByEmailAndPhoneNumberReq
		respPb *pbAdminCms.GetUserIDByEmailAndPhoneNumberResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	reqPb.Email = req.Email
	reqPb.PhoneNumber = req.PhoneNumber
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbAdminCms.NewAdminCMSClient(etcdConn)
	respPb, err := client.GetUserIDByEmailAndPhoneNumber(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.UserIDList = respPb.UserIDList
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp})
}

package user

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	pb "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserById(c *gin.Context) {
	var (
		req   cms_api_struct.GetUserRequest
		resp  cms_api_struct.GetUserResponse
		reqPb pb.GetUserByIdReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetUserById(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	if respPb.User.UserId == "" {
		openIMHttp.RespHttp200(c, constant.OK, nil)
		return
	}
	utils.CopyStructFields(&resp, respPb.User)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetUsersByName(c *gin.Context) {
	var (
		req cms_api_struct.GetUsersByNameRequest
		resp cms_api_struct.GetUsersByNameResponse
		reqPb pb.GetUsersByNameReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.UserName = req.UserName
	reqPb.Pagination = &commonPb.RequestPagination{
		PageNumber: int32(req.PageNumber),
		ShowNumber: int32(req.ShowNumber),
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetUsersByName(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	utils.CopyStructFields(&resp.Users, respPb.Users)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.UserNums = respPb.UserNums
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetUsers(c *gin.Context) {
	var (
		req   cms_api_struct.GetUsersRequest
		resp  cms_api_struct.GetUsersResponse
		reqPb pb.GetUsersReq
	)
	reqPb.Pagination = &commonPb.RequestPagination{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb.Pagination, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetUsers(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err, resp)
		return
	}
	utils.CopyStructFields(&resp.Users, respPb.User)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.UserNums = respPb.UserNums
	openIMHttp.RespHttp200(c, constant.OK, resp)

}

func ResignUser(c *gin.Context) {
	var (
		req   cms_api_struct.ResignUserRequest
		resp  cms_api_struct.ResignUserResponse
		reqPb pb.ResignUserReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	fmt.Println(reqPb.UserId)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	_, err := client.ResignUser(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err, resp)
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func AlterUser(c *gin.Context) {
	var (
		req   cms_api_struct.AlterUserRequest
		resp  cms_api_struct.AlterUserResponse
		reqPb pb.AlterUserReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	_, err := client.AlterUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError("0", "microserver failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func AddUser(c *gin.Context) {
	var (
		req   cms_api_struct.AddUserRequest
		reqPb pb.AddUserReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	_, err := client.AddUser(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func BlockUser(c *gin.Context) {
	var (
		req   cms_api_struct.BlockUserRequest
		resp  cms_api_struct.BlockUserResponse
		reqPb pb.BlockUserReq
	)
	if err := c.BindJSON(&req); err != nil {
		fmt.Println(err)
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	fmt.Println(reqPb)
	_, err := client.BlockUser(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err, resp)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func UnblockUser(c *gin.Context) {
	var (
		req   cms_api_struct.UnblockUserRequest
		resp  cms_api_struct.UnBlockUserResponse
		reqPb pb.UnBlockUserReq
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	_, err := client.UnBlockUser(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err, resp)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetBlockUsers(c *gin.Context) {
	var (
		req    cms_api_struct.GetBlockUsersRequest
		resp   cms_api_struct.GetBlockUsersResponse
		reqPb  pb.GetBlockUsersReq
		respPb *pb.GetBlockUsersResp
	)
	reqPb.Pagination = &commonPb.RequestPagination{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	utils.CopyStructFields(&reqPb.Pagination, &req)
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "blockUsers", reqPb.Pagination, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetBlockUsers(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetBlockUsers rpc", err.Error())
		openIMHttp.RespHttp200(c, err, resp)
		return
	}
	for _, v := range respPb.BlockUsers {
		resp.BlockUsers = append(resp.BlockUsers, cms_api_struct.BlockUser{
			UserResponse: cms_api_struct.UserResponse{
				UserId:       v.User.UserId,
				ProfilePhoto: v.User.ProfilePhoto,
				Nickname:     v.User.Nickname,
				IsBlock:      v.User.IsBlock,
				CreateTime:   v.User.CreateTime,
			},
			BeginDisableTime: v.BeginDisableTime,
			EndDisableTime:   v.EndDisableTime,
		})
	}
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.UserNums = respPb.UserNums
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetBlockUserById(c *gin.Context) {
	var (
		req   cms_api_struct.GetBlockUserRequest
		resp  cms_api_struct.GetBlockUserResponse
		reqPb pb.GetBlockUserByIdReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.UserId = req.UserId
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetBlockUserById(context.Background(), &reqPb)
	if err != nil {
		log.NewError("0", "GetBlockUserById rpc failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.EndDisableTime = respPb.BlockUser.EndDisableTime
	resp.BeginDisableTime = respPb.BlockUser.BeginDisableTime
	utils.CopyStructFields(&resp, respPb.BlockUser.User)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func DeleteUser(c *gin.Context) {
	var (
		req   cms_api_struct.DeleteUserRequest
		reqPb pb.DeleteUserReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.UserId = req.UserId
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	_, err := client.DeleteUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError("0", "DeleteUser rpc failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

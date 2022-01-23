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

func GetUser(c *gin.Context) {
	var (
		req    cms_api_struct.GetUserRequest
		resp   cms_api_struct.GetUserResponse
		reqPb  pb.GetUserReq
		respPb *pb.GetUserResp
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, err.(constant.ErrInfo), nil)
		return
	}
	// resp.UserId = resp.UserId
	// resp.Nickname = resp.UserId
	// resp.ProfilePhoto = resp.ProfilePhoto
	// resp.UserResponse =
	utils.CopyStructFields(&resp, respPb)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetUsers(c *gin.Context) {
	var (
		req    cms_api_struct.GetUsersRequest
		resp   cms_api_struct.GetUsersResponse
		reqPb  pb.GetUsersReq
		respPb *pb.GetUsersResp
	)
	reqPb.Pagination = &commonPb.RequestPagination{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb.Pagination, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetUsers(context.Background(), &reqPb)
	for _, v := range respPb.User {
		resp.Users = append(resp.Users, &cms_api_struct.UserResponse{
			ProfilePhoto: v.ProfilePhoto,
			Nickname:     v.Nickname,
			UserId:       v.UserID,
			CreateTime:   v.CreateTime,
		})
	}
	if err != nil {
		openIMHttp.RespHttp200(c, err.(constant.ErrInfo), resp)
	}
	fmt.Println(resp)
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
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	_, err := client.ResignUser(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, constant.ErrDB, resp)
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func AlterUser(c *gin.Context) {
	var (
		req    cms_api_struct.AlterUserRequest
		resp   cms_api_struct.AlterUserResponse
		reqPb  pb.AlterUserReq
		respPb *pb.AlterUserResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.AlterUser(context.Background(), &reqPb)
	fmt.Println(respPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err.(constant.ErrInfo), resp)
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func AddUser(c *gin.Context) {
	var (
		req    cms_api_struct.AddUserRequest
		resp   cms_api_struct.AddUserResponse
		reqPb  pb.AddUserReq
		respPb *pb.AddUserResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.AddUser(context.Background(), &reqPb)
	fmt.Println(respPb)
	if err != nil {

	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func BlockUser(c *gin.Context) {
	var (
		req    cms_api_struct.BlockUserRequest
		resp   cms_api_struct.BlockUserResponse
		reqPb  pb.BlockUserReq
		respPb *pb.BlockUserResp
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	utils.CopyStructFields(&reqPb, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.BlockUser(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err.(constant.ErrInfo), resp)
	}
	fmt.Println(respPb)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func UnblockUser(c *gin.Context) {
	var (
		req    cms_api_struct.UnblockUserRequest
		resp   cms_api_struct.UnBlockUserResponse
		reqPb  pb.UnBlockUserReq
		respPb *pb.UnBlockUserResp
	)
	utils.CopyStructFields(&reqPb, req)
	if err := c.ShouldBind(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.UnBlockUser(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err.(constant.ErrInfo), resp)
	}
	fmt.Println(respPb)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetBlockUsers(c *gin.Context) {
	var (
		req    cms_api_struct.GetBlockUsersRequest
		resp   cms_api_struct.GetOrganizationsResponse
		reqPb  pb.GetBlockUsersReq
		respPb *pb.GetBlockUsersResp
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetBlockUsers(context.Background(), &reqPb)
	if err != nil {
		openIMHttp.RespHttp200(c, err.(constant.ErrInfo), resp)
	}
	fmt.Println(respPb)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

package statistics

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pb "Open_IM/pkg/proto/statistics"
	"Open_IM/pkg/utils"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetMessagesStatistics(c *gin.Context) {
	var (
		req   cms_api_struct.GetMessageStatisticsRequest
		resp  cms_api_struct.GetMessageStatisticsResponse
		reqPb pb.GetMessageStatisticsReq
	)
	reqPb.StatisticsReq = &pb.StatisticsReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb.StatisticsReq, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImStatisticsName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetMessageStatistics(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetMessageStatistics failed", err.Error())
		openIMHttp.RespHttp200(c, err, resp)
		return
	}
	// utils.CopyStructFields(&resp, respPb)
	resp.GroupMessageNum = int(respPb.GroupMessageNum)
	resp.PrivateMessageNum = int(respPb.PrivateMessageNum)
	for _, v := range respPb.PrivateMessageNumList {
		resp.PrivateMessageNumList = append(resp.PrivateMessageNumList, struct {
			Date       string "json:\"date\""
			MessageNum int    "json:\"message_num\""
		}{
			Date:       v.Date,
			MessageNum: int(v.Num),
		})
	}
	for _, v := range respPb.GroupMessageNumList {
		resp.GroupMessageNumList = append(resp.GroupMessageNumList, struct {
			Date       string "json:\"date\""
			MessageNum int    "json:\"message_num\""
		}{
			Date:       v.Date,
			MessageNum: int(v.Num),
		})
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetUserStatistics(c *gin.Context) {
	var (
		req   cms_api_struct.GetUserStatisticsRequest
		resp  cms_api_struct.GetUserStatisticsResponse
		reqPb pb.GetUserStatisticsReq
	)
	reqPb.StatisticsReq = &pb.StatisticsReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb.StatisticsReq, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImStatisticsName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetUserStatistics(context.Background(), &reqPb)
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "GetUserStatistics failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	// utils.CopyStructFields(&resp, respPb)
	resp.ActiveUserNum = int(respPb.ActiveUserNum)
	resp.IncreaseUserNum = int(respPb.IncreaseUserNum)
	resp.TotalUserNum = int(respPb.TotalUserNum)
	for _, v := range respPb.ActiveUserNumList {
		resp.ActiveUserNumList = append(resp.ActiveUserNumList, struct {
			Date          string "json:\"date\""
			ActiveUserNum int    "json:\"active_user_num\""
		}{
			Date:          v.Date,
			ActiveUserNum: int(v.Num),
		})
	}
	for _, v := range respPb.IncreaseUserNumList {
		resp.IncreaseUserNumList = append(resp.IncreaseUserNumList, struct {
			Date            string "json:\"date\""
			IncreaseUserNum int    "json:\"increase_user_num\""
		}{
			Date:            v.Date,
			IncreaseUserNum: int(v.Num),
		})
	}
	for _, v := range respPb.TotalUserNumList {
		resp.TotalUserNumList = append(resp.TotalUserNumList, struct {
			Date         string "json:\"date\""
			TotalUserNum int    "json:\"total_user_num\""
		}{
			Date:         v.Date,
			TotalUserNum: int(v.Num),
		})
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetGroupStatistics(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupStatisticsRequest
		resp  cms_api_struct.GetGroupStatisticsResponse
		reqPb pb.GetGroupStatisticsReq
	)
	reqPb.StatisticsReq = &pb.StatisticsReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb.StatisticsReq, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImStatisticsName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetGroupStatistics(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetGroupStatistics failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	// utils.CopyStructFields(&resp, respPb)
	resp.IncreaseGroupNum = int(respPb.GetIncreaseGroupNum())
	resp.TotalGroupNum = int(respPb.GetTotalGroupNum())
	for _, v := range respPb.IncreaseGroupNumList {
		resp.IncreaseGroupNumList = append(resp.IncreaseGroupNumList,
			struct {
				Date             string "json:\"date\""
				IncreaseGroupNum int    "json:\"increase_group_num\""
			}{
				Date:             v.Date,
				IncreaseGroupNum: int(v.Num),
			})
	}
	for _, v := range respPb.TotalGroupNumList {
		resp.TotalGroupNumList = append(resp.TotalGroupNumList,
			struct {
				Date          string "json:\"date\""
				TotalGroupNum int    "json:\"total_group_num\""
			}{
				Date:          v.Date,
				TotalGroupNum: int(v.Num),
			})

	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetActiveUser(c *gin.Context) {
	var (
		req   cms_api_struct.GetActiveUserRequest
		resp  cms_api_struct.GetActiveUserResponse
		reqPb pb.GetActiveUserReq
	)
	reqPb.StatisticsReq = &pb.StatisticsReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb.StatisticsReq, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImStatisticsName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetActiveUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetActiveUser failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	utils.CopyStructFields(&resp.ActiveUserList, respPb.Users)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetActiveGroup(c *gin.Context) {
	var (
		req   cms_api_struct.GetActiveGroupRequest
		resp  cms_api_struct.GetActiveGroupResponse
		reqPb pb.GetActiveGroupReq
	)
	reqPb.StatisticsReq = &pb.StatisticsReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb.StatisticsReq, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImStatisticsName)
	client := pb.NewUserClient(etcdConn)
	respPb, err := client.GetActiveGroup(context.Background(), &reqPb)
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "GetActiveGroup failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	for _, group := range respPb.Groups {
		resp.ActiveGroupList = append(resp.ActiveGroupList, struct {
			GroupName  string "json:\"group_name\""
			GroupId    string "json:\"group_id\""
			MessageNum int    "json:\"message_num\""
		}{
			GroupName:  group.GroupName,
			GroupId:    group.GroupId,
			MessageNum: int(group.MessageNum),
		})
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

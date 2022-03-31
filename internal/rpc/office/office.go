package office

import (
	"Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbOffice "Open_IM/pkg/proto/office"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type officeServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewOfficeServer(port int) *officeServer {
	log.NewPrivateLog("office")
	return &officeServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImOfficeName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *officeServer) Run() {
	log.NewInfo("0", "officeServer rpc start ")
	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("0", "Listen failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("0", "listen network success, ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbOffice.RegisterOfficeServiceServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *officeServer) GetUserTags(_ context.Context, req *pbOffice.GetUserTagsReq) (resp *pbOffice.GetUserTagsResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req ", req.String())
	resp = &pbOffice.GetUserTagsResp{
		CommonResp: &pbOffice.CommonResp{},
		Tags:       []*pbOffice.Tag{},
	}
	tags, err := db.DB.GetUserTags(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "tags: ", tags)
	for _, v := range tags {
		tag := &pbOffice.Tag{
			TagID:   v.TagID,
			TagName: v.TagName,
		}
		for _, userID := range v.UserList {
			UserName, err := im_mysql_model.GetUserNameByUserID(userID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserNameByUserID failed", err.Error())
				continue
			}
			tag.UserList = append(tag.UserList, &pbOffice.TagUser{
				UserID:   userID,
				UserName: UserName,
			})
		}
		resp.Tags = append(resp.Tags, tag)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp ", resp.String())
	return resp, nil
}

func (s *officeServer) CreateTag(_ context.Context, req *pbOffice.CreateTagReq) (resp *pbOffice.CreateTagResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "CreateTag req", req.String())
	userIDList := utils.RemoveUserIDRepByMap(req.UserIDList)
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "userIDList: ", userIDList)
	resp = &pbOffice.CreateTagResp{CommonResp: &pbOffice.CommonResp{}}
	if err := db.DB.CreateTag(req.UserID, req.TagName, userIDList); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp", resp.String())
	return resp, nil
}

func (s *officeServer) DeleteTag(_ context.Context, req *pbOffice.DeleteTagReq) (resp *pbOffice.DeleteTagResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.DeleteTagResp{CommonResp: &pbOffice.CommonResp{}}
	if err := db.DB.DeleteTag(req.UserID, req.TagID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DeleteTag failed", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) SetTag(_ context.Context, req *pbOffice.SetTagReq) (resp *pbOffice.SetTagResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.SetTagResp{CommonResp: &pbOffice.CommonResp{}}
	IncreaseUserIDList := utils.RemoveUserIDRepByMap(req.IncreaseUserIDList)
	reduceUserIDList := utils.RemoveUserIDRepByMap(req.ReduceUserIDList)
	if err := db.DB.SetTag(req.UserID, req.TagID, req.NewName, IncreaseUserIDList, reduceUserIDList); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetTag failed", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) SendMsg2Tag(_ context.Context, req *pbOffice.SendMsg2TagReq) (resp *pbOffice.SendMsg2TagResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.SendMsg2TagResp{CommonResp: &pbOffice.CommonResp{}}
	var tagUserIDList []string
	for _, tagID := range req.TagList {
		userIDList, err := db.DB.GetUserIDListByTagID(req.SendID, tagID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserIDListByTagID failed", err.Error())
			continue
		}
		tagUserIDList = append(tagUserIDList, userIDList...)
	}
	var groupUserIDList []string
	for _, groupID := range req.GroupList {
		userIDList, err := im_mysql_model.GetGroupMemberIDListByGroupID(groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMemberIDListByGroupID failed", err.Error())
			continue
		}
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), userIDList)
		groupUserIDList = append(groupUserIDList, userIDList...)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), groupUserIDList, req.GroupList)
	var userIDList []string
	userIDList = append(userIDList, tagUserIDList...)
	userIDList = append(userIDList, groupUserIDList...)
	userIDList = append(userIDList, req.UserList...)
	userIDList = utils.RemoveUserIDRepByMap(userIDList)
	for i, userID := range userIDList {
		if userID == req.SendID || userID == "" {
			userIDList = append(userIDList[:i], userIDList[i+1:]...)
		}
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "total userIDList result: ", userIDList)
	for _, userID := range userIDList {
		msg.TagSendMessage(req.OperationID, req.SendID, userID, req.Content, req.SenderPlatformID)
	}

	if err := db.DB.SaveTagSendLog(req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SaveTagSendLog failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) GetTagSendLogs(_ context.Context, req *pbOffice.GetTagSendLogsReq) (resp *pbOffice.GetTagSendLogsResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.GetTagSendLogsResp{
		CommonResp: &pbOffice.CommonResp{},
		Pagination: &pbCommon.ResponsePagination{
			CurrentPage: req.Pagination.PageNumber,
			ShowNumber:  req.Pagination.ShowNumber,
		},
		TagSendLogs: []*pbOffice.TagSendLog{},
	}
	tagSendLogs, err := db.DB.GetTagSendLogs(req.UserID, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTagSendLogs", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	if err := utils.CopyStructFields(&resp.TagSendLogs, tagSendLogs); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) GetUserTagByID(_ context.Context, req *pbOffice.GetUserTagByIDReq) (resp *pbOffice.GetUserTagByIDResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.GetUserTagByIDResp{CommonResp: &pbOffice.CommonResp{}}
	tag, err := db.DB.GetTagByID(req.UserID, req.TagID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTagByID failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	for _, userID := range tag.UserList {
		userName, err := im_mysql_model.GetUserNameByUserID(userID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserNameByUserID failed", err.Error())
			continue
		}
		resp.Tag.UserList = append(resp.Tag.UserList, &pbOffice.TagUser{
			UserID:   userID,
			UserName: userName,
		})
	}
	resp.Tag.TagID = tag.TagID
	resp.Tag.TagName = tag.TagName
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

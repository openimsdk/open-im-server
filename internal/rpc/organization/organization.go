package organization

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/auth"
	rpc "Open_IM/pkg/proto/organization"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"
)

type organizationServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewServer(port int) *organizationServer {
	log.NewPrivateLog(constant.LogFileName)
	return &organizationServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImOrganizationName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *organizationServer) Run() {
	log.NewInfo("", "organization rpc start ")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.NewError("", "Listen failed ", err.Error(), address)
		return
	}
	log.NewInfo("", "listen network success, ", address, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	rpc.RegisterOrganizationServer(srv, s)
	rpcRegisterIP := ""
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("", "RegisterEtcd failed ", err.Error())
		return
	}
	log.NewInfo("", "organization rpc RegisterEtcd success", rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("", "organization rpc success")
}

func (s *organizationServer) CreateDepartment(ctx context.Context, req *rpc.CreateDepartmentReq) (*rpc.CreateDepartmentResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.CreateDepartmentResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}

	department := db.Department{}
	utils.CopyStructFields(&department, req.DepartmentInfo)
	if department.DepartmentID == "" {
		department.DepartmentID = utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	}
	log.Debug(req.OperationID, "dst ", department, "src ", req.DepartmentInfo)
	if err := imdb.CreateDepartment(&department); err != nil {
		errMsg := req.OperationID + " " + "CreateDepartment failed " + err.Error()
		log.Error(req.OperationID, errMsg)
		return &rpc.CreateDepartmentResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	err, createdDepartment := imdb.GetDepartment(department.DepartmentID)
	if err != nil {
		errMsg := req.OperationID + " " + "GetDepartment failed " + err.Error() + department.DepartmentID
		log.Error(req.OperationID, errMsg)
		return &rpc.CreateDepartmentResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "GetDepartment ", department.DepartmentID, *createdDepartment)
	resp := &rpc.CreateDepartmentResp{DepartmentInfo: &open_im_sdk.Department{}}
	utils.CopyStructFields(resp.DepartmentInfo, createdDepartment)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", *resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) UpdateDepartment(ctx context.Context, req *rpc.UpdateDepartmentReq) (*rpc.UpdateDepartmentResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.UpdateDepartmentResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}

	department := db.Department{}
	utils.CopyStructFields(&department, req.DepartmentInfo)

	log.Debug(req.OperationID, "dst ", department, "src ", req.DepartmentInfo)
	if err := imdb.UpdateDepartment(&department, nil); err != nil {
		errMsg := req.OperationID + " " + "UpdateDepartment failed " + err.Error()
		log.Error(req.OperationID, errMsg)
		return &rpc.UpdateDepartmentResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}

	resp := &rpc.UpdateDepartmentResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", *resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) GetSubDepartment(ctx context.Context, req *rpc.GetSubDepartmentReq) (*rpc.GetSubDepartmentResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	err, departmentList := imdb.GetSubDepartmentList(req.DepartmentID)
	if err != nil {
		errMsg := req.OperationID + " " + "GetDepartment failed " + err.Error()
		log.Error(req.OperationID, errMsg)
		return &rpc.GetSubDepartmentResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "GetSubDepartmentList ", req.DepartmentID, departmentList)
	resp := &rpc.GetSubDepartmentResp{}
	for _, v := range departmentList {
		v1 := open_im_sdk.Department{}
		utils.CopyStructFields(&v1, v)
		log.Debug(req.OperationID, "src ", v, "dst ", v1)
		err, v1.MemberNum = imdb.GetDepartmentMemberNum(v1.DepartmentID)
		if err != nil {
			log.Error(req.OperationID, "GetDepartmentMemberNum failed ", err.Error(), v1.DepartmentID)
			continue
		}
		err, v1.SubDepartmentNum = imdb.GetSubDepartmentNum(v1.DepartmentID)
		if err != nil {
			log.Error(req.OperationID, "GetSubDepartmentNum failed ", err.Error(), v1.DepartmentID)
			continue
		}
		resp.DepartmentList = append(resp.DepartmentList, &v1)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", *resp)
	return resp, nil
}

func (s *organizationServer) DeleteDepartment(ctx context.Context, req *rpc.DeleteDepartmentReq) (*rpc.DeleteDepartmentResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.DeleteDepartmentResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	err := imdb.DeleteDepartment(req.DepartmentID)
	if err != nil {
		errMsg := req.OperationID + " " + "DeleteDepartment failed " + err.Error()
		log.Error(req.OperationID, errMsg, req.DepartmentID)
		return &rpc.DeleteDepartmentResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "DeleteDepartment ", req.DepartmentID)
	resp := &rpc.DeleteDepartmentResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) CreateOrganizationUser(ctx context.Context, req *rpc.CreateOrganizationUserReq) (*rpc.CreateOrganizationUserResp, error) {
	authReq := &pbAuth.UserRegisterReq{UserInfo: &open_im_sdk.UserInfo{}}
	utils.CopyStructFields(authReq.UserInfo, req.OrganizationUser)
	authReq.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName)
	client := pbAuth.NewAuthClient(etcdConn)

	reply, err := client.UserRegister(context.Background(), authReq)
	if err != nil {
		errMsg := "UserRegister failed " + err.Error()
		log.NewError(req.OperationID, errMsg)
		return &rpc.CreateOrganizationUserResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	if reply.CommonResp.ErrCode != 0 {
		errMsg := "UserRegister failed " + reply.CommonResp.ErrMsg
		log.NewError(req.OperationID, errMsg)
		return &rpc.CreateOrganizationUserResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.CreateOrganizationUserResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	organizationUser := db.OrganizationUser{}
	utils.CopyStructFields(&organizationUser, req.OrganizationUser)
	organizationUser.Birth = utils.UnixSecondToTime(int64(req.OrganizationUser.Birth))
	log.Debug(req.OperationID, "src ", *req.OrganizationUser, "dst ", organizationUser)
	err = imdb.CreateOrganizationUser(&organizationUser)
	if err != nil {
		errMsg := req.OperationID + " " + "CreateOrganizationUser failed " + err.Error()
		log.Error(req.OperationID, errMsg, organizationUser)
		return &rpc.CreateOrganizationUserResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CreateOrganizationUser ", organizationUser)
	resp := &rpc.CreateOrganizationUserResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", *resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) UpdateOrganizationUser(ctx context.Context, req *rpc.UpdateOrganizationUserReq) (*rpc.UpdateOrganizationUserResp, error) {
	authReq := &pbAuth.UserRegisterReq{UserInfo: &open_im_sdk.UserInfo{}}
	utils.CopyStructFields(authReq.UserInfo, req.OrganizationUser)
	authReq.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName)
	client := pbAuth.NewAuthClient(etcdConn)

	reply, err := client.UserRegister(context.Background(), authReq)
	if err != nil {
		errMsg := "UserRegister failed " + err.Error()
		log.NewError(req.OperationID, errMsg)
		return &rpc.UpdateOrganizationUserResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	if reply.CommonResp.ErrCode != 0 {
		errMsg := "UserRegister failed " + reply.CommonResp.ErrMsg
		log.NewError(req.OperationID, errMsg)
		return &rpc.UpdateOrganizationUserResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.UpdateOrganizationUserResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	organizationUser := db.OrganizationUser{}
	utils.CopyStructFields(&organizationUser, req.OrganizationUser)
	if req.OrganizationUser.Birth != 0 {
		organizationUser.Birth = utils.UnixSecondToTime(int64(req.OrganizationUser.Birth))
		log.Debug(req.OperationID, "time: ", organizationUser.Birth, req.OrganizationUser.Birth)
	}

	log.Debug(req.OperationID, "src ", *req.OrganizationUser, "dst ", organizationUser)
	err = imdb.UpdateOrganizationUser(&organizationUser, nil)
	if err != nil {
		errMsg := req.OperationID + " " + "CreateOrganizationUser failed " + err.Error()
		log.Error(req.OperationID, errMsg, organizationUser)
		return &rpc.UpdateOrganizationUserResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "UpdateOrganizationUser ", organizationUser)
	resp := &rpc.UpdateOrganizationUserResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) CreateDepartmentMember(ctx context.Context, req *rpc.CreateDepartmentMemberReq) (*rpc.CreateDepartmentMemberResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.CreateDepartmentMemberResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}

	err, _ := imdb.GetOrganizationUser(req.DepartmentMember.UserID)
	if err != nil {
		errMsg := req.OperationID + " " + req.DepartmentMember.UserID + " is not exist"
		log.Error(req.OperationID, errMsg)
		return &rpc.CreateDepartmentMemberResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}

	departmentMember := db.DepartmentMember{}
	utils.CopyStructFields(&departmentMember, req.DepartmentMember)
	log.Debug(req.OperationID, "src ", *req.DepartmentMember, "dst ", departmentMember)
	err = imdb.CreateDepartmentMember(&departmentMember)
	if err != nil {
		errMsg := req.OperationID + " " + "CreateDepartmentMember failed " + err.Error()
		log.Error(req.OperationID, errMsg, departmentMember)
		return &rpc.CreateDepartmentMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "UpdateOrganizationUser ", departmentMember)

	resp := &rpc.CreateDepartmentMemberResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", *resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) GetUserInDepartmentByUserID(userID string, operationID string) (*open_im_sdk.UserInDepartment, error) {
	err, organizationUser := imdb.GetOrganizationUser(userID)
	if err != nil {
		return nil, utils.Wrap(err, "GetOrganizationUser failed")
	}
	err, departmentMemberList := imdb.GetUserInDepartment(userID)
	if err != nil {
		return nil, utils.Wrap(err, "GetUserInDepartment failed")
	}
	log.Debug(operationID, "GetUserInDepartment  ", departmentMemberList)
	resp := &open_im_sdk.UserInDepartment{OrganizationUser: &open_im_sdk.OrganizationUser{}}
	utils.CopyStructFields(resp.OrganizationUser, organizationUser)
	for _, v := range departmentMemberList {
		v1 := open_im_sdk.DepartmentMember{}
		utils.CopyStructFields(&v1, v)
		log.Debug(operationID, "DepartmentMember src  ", v, "dst ", v1)
		resp.DepartmentMemberList = append(resp.DepartmentMemberList, &v1)
	}
	return resp, nil
}

func (s *organizationServer) GetUserInDepartment(ctx context.Context, req *rpc.GetUserInDepartmentReq) (*rpc.GetUserInDepartmentResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	r, err := s.GetUserInDepartmentByUserID(req.UserID, req.OperationID)
	if err != nil {
		errMsg := req.OperationID + " " + "GetUserInDepartmentByUserID failed " + err.Error()
		log.Error(req.OperationID, errMsg, req.UserID)
		return &rpc.GetUserInDepartmentResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "GetUserInDepartmentByUserID success ", req.UserID, r)
	resp := rpc.GetUserInDepartmentResp{UserInDepartment: &open_im_sdk.UserInDepartment{OrganizationUser: &open_im_sdk.OrganizationUser{}}}
	resp.UserInDepartment.DepartmentMemberList = r.DepartmentMemberList
	resp.UserInDepartment.OrganizationUser = r.OrganizationUser
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", resp)
	return &resp, nil
}

func (s *organizationServer) UpdateUserInDepartment(ctx context.Context, req *rpc.UpdateUserInDepartmentReq) (*rpc.UpdateUserInDepartmentResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.UpdateUserInDepartmentResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	departmentMember := &db.DepartmentMember{}
	utils.CopyStructFields(departmentMember, req.DepartmentMember)
	log.Debug(req.OperationID, "dst ", departmentMember, "src ", req.DepartmentMember)
	err := imdb.UpdateUserInDepartment(departmentMember, nil)
	if err != nil {
		errMsg := req.OperationID + " " + "UpdateUserInDepartment failed " + err.Error()
		log.Error(req.OperationID, errMsg, *departmentMember)
		return &rpc.UpdateUserInDepartmentResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	resp := &rpc.UpdateUserInDepartmentResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", *resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) DeleteUserInDepartment(ctx context.Context, req *rpc.DeleteUserInDepartmentReq) (*rpc.DeleteUserInDepartmentResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.DeleteUserInDepartmentResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}

	err := imdb.DeleteUserInDepartment(req.DepartmentID, req.UserID)
	if err != nil {
		errMsg := req.OperationID + " " + "DeleteUserInDepartment failed " + err.Error()
		log.Error(req.OperationID, errMsg, req.DepartmentID, req.UserID)
		return &rpc.DeleteUserInDepartmentResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "DeleteUserInDepartment success ", req.DepartmentID, req.UserID)
	resp := &rpc.DeleteUserInDepartmentResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", *resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) DeleteOrganizationUser(ctx context.Context, req *rpc.DeleteOrganizationUserReq) (*rpc.DeleteOrganizationUserResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.DeleteOrganizationUserResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	err := imdb.DeleteOrganizationUser(req.UserID)
	if err != nil {
		errMsg := req.OperationID + " " + "DeleteOrganizationUser failed " + err.Error()
		log.Error(req.OperationID, errMsg, req.UserID)
		return &rpc.DeleteOrganizationUserResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "DeleteOrganizationUser success ", req.UserID)
	resp := &rpc.DeleteOrganizationUserResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", *resp)
	chat.OrganizationNotificationToAll(req.OpUserID, req.OperationID)
	return resp, nil
}

func (s *organizationServer) GetDepartmentMember(ctx context.Context, req *rpc.GetDepartmentMemberReq) (*rpc.GetDepartmentMemberResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	err, departmentMemberList := imdb.GetDepartmentMemberList(req.DepartmentID)
	if err != nil {
		errMsg := req.OperationID + " " + req.OpUserID + " is not app manager"
		log.Error(req.OperationID, errMsg)
		return &rpc.GetDepartmentMemberResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}

	log.Debug(req.OperationID, "GetDepartmentMemberList ", departmentMemberList)
	resp := rpc.GetDepartmentMemberResp{}
	for _, v := range departmentMemberList {
		err, organizationUser := imdb.GetOrganizationUser(v.UserID)
		if err != nil {
			log.Error(req.OperationID, "GetOrganizationUser failed ", err.Error())
			continue
		}
		respOrganizationUser := &open_im_sdk.OrganizationUser{}
		respDepartmentMember := &open_im_sdk.DepartmentMember{}

		utils.CopyStructFields(respOrganizationUser, organizationUser)
		utils.CopyStructFields(respDepartmentMember, &v)
		userDepartmentMember := open_im_sdk.UserDepartmentMember{OrganizationUser: respOrganizationUser, DepartmentMember: respDepartmentMember}

		log.Debug(req.OperationID, "GetUserInDepartmentByUserID success ", userDepartmentMember)
		resp.UserDepartmentMemberList = append(resp.UserDepartmentMemberList, &userDepartmentMember)
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", resp)
	return &resp, nil
}

package organization

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/organization"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"time"

	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type organizationServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewGroupServer(port int) *organizationServer {
	log.NewPrivateLog(constant.LogFileName)
	return &organizationServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImGroupName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *organizationServer) Run() {
	log.NewInfo("", "organization rpc start ")
	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("", "Listen failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("", "listen network success, ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	rpc.RegisterOrganizationServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("", "RegisterEtcd failed ", err.Error())
		return
	}
	log.NewInfo("", "organization rpc RegisterEtcd success", ip, s.rpcPort, s.rpcRegisterName, 10)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("", "organization rpc success")
}

func (s *organizationServer) CreateDepartment(ctx context.Context, req *rpc.CreateDepartmentReq) (*rpc.CreateDepartmentResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())
	if !token_verify.IsMangerUserID(req.OpUserID) {
		errMsg := "is not app manger " + req.OpUserID + " " + req.OperationID
		log.Error(req.OperationID, errMsg)
		return &rpc.CreateDepartmentResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	department := db.Department{}
	utils.CopyStructFields(&department, req.DepartmentInfo)
	if department.DepartmentID == "" {
		department.DepartmentID = utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	}
	if imdb.CreateDepartment(&department) != nil {

	}

	err, createdDepartment := imdb.GetDepartment(department.DepartmentID)
	if err != nil {

	}
	resp := &rpc.CreateDepartmentResp{DepartmentInfo: &open_im_sdk.Department{}}
	utils.CopyStructFields(resp.DepartmentInfo, createdDepartment)

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", resp)
	return resp, nil

}

func (s *organizationServer) UpdateDepartment(ctx context.Context, req *rpc.UpdateDepartmentReq) (*rpc.UpdateDepartmentResp, error) {
	return nil, nil
}

func (s *organizationServer) GetDepartment(ctx context.Context, req *rpc.GetDepartmentReq) (*rpc.GetDepartmentResp, error) {
	return nil, nil
}

func (s *organizationServer) DeleteDepartment(ctx context.Context, req *rpc.DeleteDepartmentReq) (*rpc.DeleteDepartmentResp, error) {
	return nil, nil
}

func (s *organizationServer) CreateOrganizationUser(ctx context.Context, req *rpc.CreateOrganizationUserReq) (*rpc.CreateOrganizationUserResp, error) {
	return nil, nil
}

func (s *organizationServer) UpdateOrganizationUser(ctx context.Context, req *rpc.UpdateOrganizationUserReq) (*rpc.UpdateOrganizationUserResp, error) {
	return nil, nil
}

func (s *organizationServer) CreateDepartmentMember(ctx context.Context, req *rpc.CreateDepartmentMemberReq) (*rpc.CreateDepartmentMemberResp, error) {
	return nil, nil
}

func (s *organizationServer) GetUserInDepartment(ctx context.Context, req *rpc.GetUserInDepartmentReq) (*rpc.GetUserInDepartmentResp, error) {
	return nil, nil
}

func (s *organizationServer) UpdateUserInDepartment(ctx context.Context, req *rpc.UpdateUserInDepartmentReq) (*rpc.UpdateUserInDepartmentResp, error) {
	return nil, nil
}

func (s *organizationServer) DeleteOrganizationUser(ctx context.Context, req *rpc.DeleteOrganizationUserReq) (*rpc.DeleteOrganizationUserResp, error) {
	return nil, nil
}

func (s *organizationServer) GetDepartmentMember(ctx context.Context, req *rpc.GetDepartmentMemberReq) (*rpc.GetDepartmentMemberResp, error) {
	return nil, nil
}

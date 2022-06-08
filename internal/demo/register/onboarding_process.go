package register

import (
	"Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	groupRpc "Open_IM/pkg/proto/group"
	organizationRpc "Open_IM/pkg/proto/organization"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"strings"
	"time"
)

func onboardingProcess(operationID, userID, userName string) {
	if err := createOrganizationUser(operationID, userID, userName); err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "createOrganizationUser failed", err.Error())
	}
	departmentID := config.Config.Demo.TestDepartMentID

	if err := joinTestDepartment(operationID, userID, departmentID); err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "joinTestDepartment failed", err.Error())
	}

	groupIDList, err := GetDepartmentGroupIDList(operationID, departmentID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), groupIDList)
	joinGroups(operationID, userID, userName, groupIDList)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "fineshed")
}

func createOrganizationUser(operationID, userID, userName string) error {
	defer func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), userID)
	}()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "start createOrganizationUser")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := organizationRpc.NewOrganizationClient(etcdConn)
	req := &organizationRpc.CreateOrganizationUserReq{
		OrganizationUser: &commonPb.OrganizationUser{
			UserID:      userID,
			Nickname:    userName,
			EnglishName: randomEnglishName(),
			Gender:      constant.Male,
			CreateTime:  uint32(time.Now().Unix()),
		},
		OperationID: operationID,
		OpUserID:    config.Config.Manager.AppManagerUid[0],
		IsRegister:  false,
	}
	if strings.Contains("@", userID) {
		req.OrganizationUser.Email = userID
	} else {
		req.OrganizationUser.Telephone = userID
	}
	resp, err := client.CreateOrganizationUser(context.Background(), req)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		return err
	}
	if resp.ErrCode != 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), resp)
		return errors.New(resp.ErrMsg)
	}
	return nil
}

func joinTestDepartment(operationID, userID, departmentID string) error {
	defer func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), userID)
	}()
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := organizationRpc.NewOrganizationClient(etcdConn)
	req := &organizationRpc.CreateDepartmentMemberReq{
		DepartmentMember: &commonPb.DepartmentMember{
			UserID:       userID,
			DepartmentID: departmentID,
			Position:     randomPosition(),
		},
		OperationID: operationID,
		OpUserID:    config.Config.Manager.AppManagerUid[0],
	}
	resp, err := client.CreateDepartmentMember(context.Background(), req)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		return err
	}
	if resp.ErrCode != 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), resp)
		return errors.New(resp.ErrMsg)
	}
	return nil
}

func GetDepartmentGroupIDList(operationID, departmentID string) ([]string, error) {
	defer func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), departmentID)
	}()
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := organizationRpc.NewOrganizationClient(etcdConn)
	req := organizationRpc.GetDepartmentParentIDListReq{
		DepartmentID: departmentID,
		OperationID:  operationID,
	}
	resp, err := client.GetDepartmentParentIDList(context.Background(), &req)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), req.String())
		return nil, err
	}
	if resp.ErrCode != 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), resp)
		return nil, errors.New(resp.ErrMsg)
	}

	resp.ParentIDList = append(resp.ParentIDList, departmentID)
	getDepartmentRelatedGroupIDListReq := organizationRpc.GetDepartmentRelatedGroupIDListReq{OperationID: operationID, DepartmentIDList: resp.ParentIDList}
	getDepartmentParentIDListResp, err := client.GetDepartmentRelatedGroupIDList(context.Background(), &getDepartmentRelatedGroupIDListReq)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), getDepartmentRelatedGroupIDListReq.String())
		return nil, err
	}
	if getDepartmentParentIDListResp.ErrCode != 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), getDepartmentParentIDListResp)
		return nil, errors.New(getDepartmentParentIDListResp.ErrMsg)
	}
	return getDepartmentParentIDListResp.GroupIDList, nil
}

func joinGroups(operationID, userID, userName string, groupIDList []string) {
	defer func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), userID, groupIDList)
	}()
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := groupRpc.NewGroupClient(etcdConn)
	for _, groupID := range groupIDList {
		req := &groupRpc.InviteUserToGroupReq{
			OperationID:       operationID,
			GroupID:           groupID,
			Reason:            "register auto join",
			InvitedUserIDList: []string{userID},
			OpUserID:          userID,
		}
		resp, err := client.InviteUserToGroup(context.Background(), req)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), req.String())
			continue
		}
		if resp.ErrCode != 0 {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), resp)
			continue
		}
		onboardingProcessNotification(operationID, userID, groupID)
	}
}

// welcome user join department notification
func onboardingProcessNotification(operationID, userID, groupID string) {
	defer func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), userID, groupID)
	}()
	var tips commonPb.TipsComm
	tips.DefaultTips = config.Config.Notification.JoinDepartmentNotification.DefaultTips.Tips
	tips.JsonDetail = ""
	content, err := proto.Marshal(&tips)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), "proto marshal failed")
		return
	}
	notification := &msg.NotificationMsg{
		SendID:      userID,
		RecvID:      groupID,
		Content:     content,
		MsgFrom:     constant.UserMsgType,
		ContentType: constant.Text,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}
	// notification user join group
	msg.Notification(notification)
}

func randomEnglishName() string {
	l := []string{"abandon", "entail", "nebula", "shrink", "accumulate", "etch", "nostalgia", "slide",
		"feudal", "adverse", "exploit", "occupy", "solve", "amazing", "fantasy", "orchid", "spiky", "approve", "flap"}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(l) - 1)
	fmt.Println(index)
	return l[index]
}

func randomPosition() string {
	l := []string{"Golang工程师", "前端工程师", "后端工程师", "产品经理", "测试开发工程师", "运维开发工程师"}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(l) - 1)
	fmt.Println(index)
	return l[index]
}

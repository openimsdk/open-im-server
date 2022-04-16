package organization

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/organization"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CreateDepartment(c *gin.Context) {
	params := api.CreateDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.CreateDepartmentReq{}
	utils.CopyStructFields(req, &params)
	utils.CopyStructFields(req.DepartmentInfo, &params.Department)
	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + " " + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.CreateDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc CreateDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.CreateDepartmentResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, Department: RpcResp.DepartmentInfo}
	apiResp.Data = jsonData.JsonDataOne(RpcResp.DepartmentInfo)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func UpdateDepartment(c *gin.Context) {
	params := api.UpdateDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.UpdateDepartmentReq{}
	utils.CopyStructFields(req, &params)
	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.UpdateDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc UpdateDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.UpdateDepartmentResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func GetSubDepartment(c *gin.Context) {
	params := api.GetDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetSubDepartmentReq{}
	utils.CopyStructFields(req, &params)
	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.GetSubDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc GetDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.GetDepartmentResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, DepartmentList: RpcResp.DepartmentList}
	apiResp.Data = jsonData.JsonDataList(RpcResp.DepartmentList)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func DeleteDepartment(c *gin.Context) {
	params := api.DeleteDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.DeleteDepartmentReq{}
	utils.CopyStructFields(req, &params)
	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.DeleteDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc DeleteDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.DeleteDepartmentResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func CreateOrganizationUser(c *gin.Context) {
	params := api.CreateOrganizationUserReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.CreateOrganizationUserReq{}
	utils.CopyStructFields(req, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.CreateOrganizationUser(context.Background(), req)
	if err != nil {
		errMsg := "rpc CreateOrganizationUser failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.CreateOrganizationUserResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func UpdateOrganizationUser(c *gin.Context) {
	params := api.UpdateOrganizationUserReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.UpdateOrganizationUserReq{}
	utils.CopyStructFields(req, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.UpdateOrganizationUser(context.Background(), req)
	if err != nil {
		errMsg := "rpc UpdateOrganizationUser failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.UpdateOrganizationUserResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func CreateDepartmentMember(c *gin.Context) {
	params := api.CreateDepartmentMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.CreateDepartmentMemberReq{}
	utils.CopyStructFields(req, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.CreateDepartmentMember(context.Background(), req)
	if err != nil {
		errMsg := "rpc CreateDepartmentMember failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.CreateDepartmentMemberResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func GetUserInDepartment(c *gin.Context) {
	params := api.GetUserInDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.GetUserInDepartmentReq{}
	utils.CopyStructFields(req, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.GetUserInDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc GetUserInDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.GetUserInDepartmentResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, UserInDepartment: RpcResp.UserInDepartment}
	apiResp.Data = jsonData.JsonDataOne(RpcResp.UserInDepartment)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func UpdateUserInDepartment(c *gin.Context) {
	params := api.UpdateUserInDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.UpdateUserInDepartmentReq{}
	utils.CopyStructFields(req, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.UpdateUserInDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc UpdateUserInDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.UpdateUserInDepartmentResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func DeleteOrganizationUser(c *gin.Context) {
	params := api.DeleteOrganizationUserReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.DeleteOrganizationUserReq{}
	utils.CopyStructFields(req, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.DeleteOrganizationUser(context.Background(), req)
	if err != nil {
		errMsg := "rpc DeleteOrganizationUser failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.DeleteOrganizationUserResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func GetDepartmentMember(c *gin.Context) {
	params := api.GetDepartmentMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.GetDepartmentMemberReq{}
	utils.CopyStructFields(req, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.GetDepartmentMember(context.Background(), req)
	if err != nil {
		errMsg := "rpc GetDepartmentMember failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.GetDepartmentMemberResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, UserInDepartmentList: RpcResp.UserInDepartmentList}
	apiResp.Data = jsonData.JsonDataList(RpcResp.UserInDepartmentList)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func DeleteUserInDepartment(c *gin.Context) {
	params := api.DeleteUserInDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.DeleteUserInDepartmentReq{}
	utils.CopyStructFields(req, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName)
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.DeleteUserInDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc DeleteUserInDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.GetDepartmentMemberResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

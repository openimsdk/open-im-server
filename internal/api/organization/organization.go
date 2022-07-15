package organization

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/organization"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// @Summary 创建部门
// @Description 创建部门
// @Tags 组织架构相关
// @ID CreateDepartment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CreateDepartmentReq true "请求"
// @Produce json
// @Success 0 {object} api.CreateDepartmentResp{data=open_im_sdk.Department}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/create_department [post]
func CreateDepartment(c *gin.Context) {
	params := api.CreateDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.CreateDepartmentReq{DepartmentInfo: &open_im_sdk.Department{}}
	utils.CopyStructFields(req, &params)
	utils.CopyStructFields(req.DepartmentInfo, &params)
	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + " " + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

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

// @Summary 更新部门信息
// @Description 更新部门信息
// @Tags 组织架构相关
// @ID UpdateDepartment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.UpdateDepartmentReq true "请求"
// @Produce json
// @Success 0 {object} api.UpdateDepartmentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/update_department [post]
func UpdateDepartment(c *gin.Context) {
	params := api.UpdateDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.UpdateDepartmentReq{DepartmentInfo: &open_im_sdk.Department{}}
	utils.CopyStructFields(req, &params)
	utils.CopyStructFields(req.DepartmentInfo, &params)
	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

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

// @Summary 获取子部门列表
// @Description 获取子部门列表
// @Tags 组织架构相关
// @ID GetSubDepartment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetSubDepartmentReq true "请求"
// @Produce json
// @Success 0 {object} api.GetSubDepartmentResp{data=[]open_im_sdk.Department}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/get_sub_department [post]
func GetSubDepartment(c *gin.Context) {
	params := api.GetSubDepartmentReq{}
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

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.GetSubDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc GetDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.GetSubDepartmentResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, DepartmentList: RpcResp.DepartmentList}
	apiResp.Data = jsonData.JsonDataList(RpcResp.DepartmentList)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

func GetAllDepartment(c *gin.Context) {

}

// @Summary 删除部门
// @Description 删除部门
// @Tags 组织架构相关
// @ID DeleteDepartment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.DeleteDepartmentReq true "请求"
// @Produce json
// @Success 0 {object} api.DeleteDepartmentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/delete_department [post]
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
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 组织架构导入用户
// @Description 组织架构导入用户
// @Tags 组织架构相关
// @ID CreateOrganizationUser
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CreateOrganizationUserReq true "请求"
// @Produce json
// @Success 0 {object} api.CreateOrganizationUserResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/create_organization_user [post]
func CreateOrganizationUser(c *gin.Context) {
	params := api.CreateOrganizationUserReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.CreateOrganizationUserReq{OrganizationUser: &open_im_sdk.OrganizationUser{}}
	utils.CopyStructFields(req, &params)
	utils.CopyStructFields(req.OrganizationUser, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 更新组织架构中的用户
// @Description 更新组织架构中的用户
// @Tags 组织架构相关
// @ID UpdateOrganizationUser
// @Accept json
// @Param token header string true "im token"
// @Param req body api.UpdateOrganizationUserReq true "请求"
// @Produce json
// @Success 0 {object} api.UpdateOrganizationUserResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/update_organization_user [post]
func UpdateOrganizationUser(c *gin.Context) {
	params := api.UpdateOrganizationUserReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.UpdateOrganizationUserReq{OrganizationUser: &open_im_sdk.OrganizationUser{}}
	utils.CopyStructFields(req, &params)
	utils.CopyStructFields(req.OrganizationUser, &params)
	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 创建部门用户
// @Description 创建部门用户
// @Tags 组织架构相关
// @ID CreateDepartmentMember
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CreateDepartmentMemberReq true "请求"
// @Produce json
// @Success 0 {object} api.CreateDepartmentMemberResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/create_department_member [post]
func CreateDepartmentMember(c *gin.Context) {
	params := api.CreateDepartmentMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.CreateDepartmentMemberReq{DepartmentMember: &open_im_sdk.DepartmentMember{}}
	utils.CopyStructFields(req, &params)
	utils.CopyStructFields(req.DepartmentMember, &params)

	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 获取部门中的所有用户
// @Description 获取部门中的所有用户
// @Tags 组织架构相关
// @ID GetUserInDepartment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUserInDepartmentReq true "请求"
// @Produce json
// @Success 0 {object} api.GetUserInDepartmentResp{data=open_im_sdk.UserInDepartment}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/get_user_in_department [post]
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
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 更新部门中某个用户
// @Description 更新部门中某个用户
// @Tags 组织架构相关
// @ID UpdateUserInDepartment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.UpdateUserInDepartmentReq true "请求"
// @Produce json
// @Success 0 {object} api.UpdateUserInDepartmentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/update_user_in_department [post]
func UpdateUserInDepartment(c *gin.Context) {
	params := api.UpdateUserInDepartmentReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.UpdateUserInDepartmentReq{DepartmentMember: &open_im_sdk.DepartmentMember{}}
	utils.CopyStructFields(req.DepartmentMember, &params)
	utils.CopyStructFields(req, &params)
	err, opUserID := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	req.OpUserID = opUserID
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 删除组织架构中某个用户
// @Description 删除组织架构中某个用户
// @Tags 组织架构相关
// @ID DeleteOrganizationUser
// @Accept json
// @Param token header string true "im token"
// @Param req body api.DeleteOrganizationUserReq true "请求"
// @Produce json
// @Success 0 {object} api.DeleteOrganizationUserResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/delete_organization_user [post]
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
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 获取部门中所有成员
// @Description 获取部门中所有成员
// @Tags 组织架构相关
// @ID GetDepartmentMember
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetDepartmentMemberReq true "请求"
// @Produce json
// @Success 0 {object} api.GetDepartmentMemberResp{data=[]open_im_sdk.UserDepartmentMember}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/get_department_member [post]
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
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.GetDepartmentMember(context.Background(), req)
	if err != nil {
		errMsg := "rpc GetDepartmentMember failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.GetDepartmentMemberResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, UserInDepartmentList: RpcResp.UserDepartmentMemberList}
	apiResp.Data = jsonData.JsonDataList(RpcResp.UserDepartmentMemberList)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

// @Summary 删除部门中某个用户
// @Description 删除部门中某个用户
// @Tags 组织架构相关
// @ID DeleteUserInDepartment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.DeleteUserInDepartmentReq true "请求"
// @Produce json
// @Success 0 {object} api.DeleteUserInDepartmentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /organization/delete_user_in_department [post]
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
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api args ", req.String(), "params", params)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOrganizationName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewOrganizationClient(etcdConn)
	RpcResp, err := client.DeleteUserInDepartment(context.Background(), req)
	if err != nil {
		errMsg := "rpc DeleteUserInDepartment failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	apiResp := api.DeleteUserInDepartmentResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "api return ", apiResp)
	c.JSON(http.StatusOK, apiResp)
}

/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/15 10:28).
 */
package manage

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbUser "Open_IM/pkg/proto/user"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// paramsDeleteUsers struct
type paramsDeleteUsers struct {
	OperationID   string   `json:"operationID" binding:"required"`
	DeleteUidList []string `json:"deleteUidList" binding:"required"`
}

// paramsGetAllUsersUid struct
type paramsGetAllUsersUid struct {
	OperationID string `json:"operationID" binding:"required"`
}

// deleteUserResult struct
type deleteUserResult struct {
	ErrCode       int      `json:"errCode" example:"0"`
	ErrMsg        string   `json:"errMsg"  example:"error"`
	FailedUidList []string `json:"failedUidList"  example:[]`
}

// uidListResult struct
type uidListResult struct {
	ErrCode int      `json:"errCode" example:"0"`
	ErrMsg  string   `json:"errMsg"  example:"error"`
	UidList []string `json:"uidList"  example:[]`
}

// @Summary
// @Schemes
// @Description delete user
// @Tags manage
// @Accept json
// @Produce json
// @Param body body manage.paramsDeleteUsers true "user to be deleted"
// @Param token header string true "token"
// @Success 200 {object} manage.deleteUserResult
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /manager/delete_user [post]
func DeleteUser(c *gin.Context) {
	params := paramsDeleteUsers{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.InfoByKv("DeleteUser req come here", params.OperationID, "DeleteUidList", params.DeleteUidList)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	//defer etcdConn.Close()

	req := &pbUser.DeleteUsersReq{
		OperationID:   params.OperationID,
		DeleteUidList: params.DeleteUidList,
		Token:         c.Request.Header.Get("token"),
	}
	RpcResp, err := client.DeleteUsers(context.Background(), req)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call delete users rpc server failed"})
		return
	}
	log.InfoByKv("call delete user rpc server is success", params.OperationID, "resp args", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.CommonResp.ErrorCode, "errMsg": RpcResp.CommonResp.ErrorMsg, "failedUidList": RpcResp.FailedUidList}
	c.JSON(http.StatusOK, resp)
}

// @Summary
// @Schemes
// @Description get all user ids
// @Tags manage
// @Accept json
// @Produce json
// @Param body body manage.paramsGetAllUsersUid true "all user ids"
// @Param token header string true "token"
// @Success 200 {object} manage.uidListResult
// @Failure 500 {object} manage.uidListResult
// @Router /manager/get_all_users_uid [post]
func GetAllUsersUid(c *gin.Context) {
	params := paramsGetAllUsersUid{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.InfoByKv("GetAllUsersUid req come here", params.OperationID)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	//defer etcdConn.Close()

	req := &pbUser.GetAllUsersUidReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	RpcResp, err := client.GetAllUsersUid(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error(), "uidList": []string{}})
		return
	}
	log.InfoByKv("call GetAllUsersUid rpc server is success", params.OperationID, "resp args", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.CommonResp.ErrorCode, "errMsg": RpcResp.CommonResp.ErrorMsg, "uidList": RpcResp.UidList}
	c.JSON(http.StatusOK, resp)

}

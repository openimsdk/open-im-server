package user

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbUser "Open_IM/pkg/proto/user"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// result struct
type result struct {
	ErrCode int         `json:"errCode" example:"0"`
	ErrMsg  string      `json:"errMsg"  example:"error"`
	Data    interface{} `json:"data,omitempty"`
}

// userInfo struct
type userInfo struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Gender int32  `json:"gender"`
	Mobile string `json:"mobile"`
	Birth  string `json:"birth"`
	Email  string `json:"email"`
	Ex     string `json:"ex"`
}

// queryUserInfoParam struct
type queryUserInfoParam struct {
	OperationID string   `json:"operationID" binding:"required"`
	UIDList     []string `json:"uidList"`
}

// @Summary
// @Description get user info by uid list
// @Tags user
// @Accept json
// @Produce json
// @Param body body user.queryUserInfoParam true "get userInfo by uidList"
// @Param token header string true "token"
// @Success 200 {object} user.result{data=[]user.userInfo}
// @Failure 400 {object} user.result
// @Router /user/get_user_info [post]
func GetUserInfo(c *gin.Context) {
	log.InfoByKv("api get userinfo init...", "")

	etcdConn := getcdv3.GetUserConn()
	client := pbUser.NewUserClient(etcdConn)
	//defer etcdConn.Close()

	params := queryUserInfoParam{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbUser.GetUserInfoReq{
		UserIDList:  params.UIDList,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.InfoByKv("api get user info is server", c.PostForm("OperationID"), c.Request.Header.Get("token"))
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call get user info rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"errCode": 500,
			"errMsg":  "call  rpc server failed",
		})
		return
	}
	log.InfoByKv("call get user info rpc server success", params.OperationID)
	if RpcResp.ErrorCode == 0 {
		userInfoList := make([]userInfo, 0)
		for _, user := range RpcResp.Data {
			var ui userInfo
			ui.UID = user.Uid
			ui.Name = user.Name
			ui.Icon = user.Icon
			ui.Gender = user.Gender
			ui.Mobile = user.Mobile
			ui.Birth = user.Birth
			ui.Email = user.Email
			ui.Ex = user.Ex
			userInfoList = append(userInfoList, ui)
		}
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg, "data": userInfoList}
		c.JSON(http.StatusOK, resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
		c.JSON(http.StatusOK, resp)
	}
	log.InfoByKv("api get user info return success", params.OperationID, "args=%s", RpcResp.String())
}

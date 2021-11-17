package friend

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// paramsRemoveBlackList struct
type paramsRemoveBlackList struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
}

// @Summary
// @Schemes
// @Description remove black list
// @Tags friend
// @Accept json
// @Produce json
// @Param body body friend.paramsSearchFriend true "remove black list params"
// @Param token header string true "token"
// @Success 200 {object} user.result
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /friend/remove_blacklist [post]
func RemoveBlacklist(c *gin.Context) {
	log.Info("", "", "api remove_blacklist init ....")

	etcdConn := getcdv3.GetFriendConn()
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsRemoveBlackList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.RemoveBlacklistReq{
		Uid:         params.UID,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api remove blacklist is server:userID=%s", req.Uid)
	RpcResp, err := client.RemoveBlacklist(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call remove blacklist rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call remove blacklist rpc server failed"})
		return
	}
	log.InfoByArgs("call remove blacklist rpc server success,args=%s", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
	c.JSON(http.StatusOK, resp)
	log.InfoByArgs("api remove blacklist success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}

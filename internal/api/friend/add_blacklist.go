package friend

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary
// @Schemes
// @Description add a user into black list
// @Tags friend
// @Accept json
// @Produce json
// @Param body body friend.paramsSearchFriend true "add black list params"
// @Param token header string true "token"
// @Success 200 {object} user.result
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /friend/add_blacklist [post]
func AddBlacklist(c *gin.Context) {
	log.Info("", "", "api add blacklist init ....")

	etcdConn := getcdv3.GetFriendConn()
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsSearchFriend{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.AddBlacklistReq{
		Uid:         params.UID,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
		OwnerUid:    params.OwnerUid,
	}
	log.Info(req.Token, req.OperationID, "api add blacklist is server:userID=%s", req.Uid)
	RpcResp, err := client.AddBlacklist(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call add blacklist rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add blacklist rpc server failed"})
		return
	}
	log.InfoByArgs("call add blacklist rpc server success,args=%s", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
	c.JSON(http.StatusOK, resp)
	log.InfoByArgs("api add blacklist success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}

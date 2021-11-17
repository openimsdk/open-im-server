package friend

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// paramsGetFriendLIst struct
type paramsGetFriendLIst struct {
	OperationID string `json:"operationID" binding:"required"`
}

// friendInfo struct
type friendInfo struct {
	UID           string `json:"uid"`
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	Gender        int32  `json:"gender"`
	Mobile        string `json:"mobile"`
	Birth         string `json:"birth"`
	Email         string `json:"email"`
	Ex            string `json:"ex"`
	Comment       string `json:"comment"`
	IsInBlackList int32  `json:"isInBlackList"`
}

// @Summary
// @Schemes
// @Description get friend apply list
// @Tags friend
// @Accept json
// @Produce json
// @Param body body friend.paramsGetFriendLIst true "get friend apply list"
// @Param token header string true "token"
// @Success 200 {object} user.result{data=[]friend.friendInfo}
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /friend/get_friend_list [post]
func GetFriendList(c *gin.Context) {
	log.Info("", "", fmt.Sprintf("api get_friendlist init ...."))

	etcdConn := getcdv3.GetFriendConn()
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsGetFriendLIst{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetFriendListReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api get friend list is server")
	RpcResp, err := client.GetFriendList(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call get friend list rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get friend list rpc server failed"})
		return
	}
	log.InfoByArgs("call get friend list rpc server success,args=%s", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		friendsInfo := make([]friendInfo, 0)
		for _, friend := range RpcResp.Data {
			var fi friendInfo
			fi.UID = friend.Uid
			fi.Name = friend.Name
			fi.Icon = friend.Icon
			fi.Gender = friend.Gender
			fi.Mobile = friend.Mobile
			fi.Birth = friend.Birth
			fi.Email = friend.Email
			fi.Ex = friend.Ex
			fi.Comment = friend.Comment
			fi.IsInBlackList = friend.IsInBlackList
			friendsInfo = append(friendsInfo, fi)
		}
		resp := gin.H{
			"errCode": RpcResp.ErrorCode,
			"errMsg":  RpcResp.ErrorMsg,
			"data":    friendsInfo,
		}
		c.JSON(http.StatusOK, resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
		c.JSON(http.StatusOK, resp)
	}
	log.InfoByArgs("api get friend list success return,get args=%s,return=%s", req.String(), RpcResp.String())
}

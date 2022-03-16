package apiChat

//
//import (
//	apiStruct "Open_IM/pkg/base_info"
//	"Open_IM/pkg/common/config"
//	"Open_IM/pkg/common/log"
//	"Open_IM/pkg/common/token_verify"
//	"Open_IM/pkg/grpc-etcdv3/getcdv3"
//	pbChat "Open_IM/pkg/proto/chat"
//	"Open_IM/pkg/utils"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strings"
//)
//
//func DelMsg(c *gin.Context) {
//	var (
//		req  apiStruct.DelMsgReq
//		resp apiStruct.DelMsgResp
//		reqPb pbChat.
//	)
//	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
//	if !ok {
//		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
//		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
//		return
//	}
//	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
//	msgClient := pbChat.NewChatClient(grpcConn)
//	//respPb, err := msgClient.DelMsgList()
//}

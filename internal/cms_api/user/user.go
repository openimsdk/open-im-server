package user

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetUser(c *gin.Context) {
	var (
		req cms_api_struct.RequestPagination
		resp cms_api_struct.GetUsersResponse
		reqPb rpc.GetUserInfoReq
		respPb *rpc.GetUserInfoResp
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	utils.CopyStructFields(req, &req)

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := rpc.NewUserClient(etcdConn)
	respPb, err := client.GetUserInfo(context.Background(), &reqPb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	//for _, v := range RpcResp.UserInfoList {
	//	publicUserInfoList = append(publicUserInfoList,
	//		&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, AppMangerLevel: v.AppMangerLevel})
	//}

	//resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
	//resp.Data = jsonData.JsonDataList(resp.UserInfoList)
	//log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func ResignUser(c *gin.Context) {

}



func AlterUser(c *gin.Context) {

}

func AddUser(c *gin.Context) {

}

func BlockUser(c *gin.Context) {

}

func UnblockUser(c *gin.Context) {

}

func GetBlockUsers(c *gin.Context) {

}

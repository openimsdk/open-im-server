package friend

import (
	common "Open_IM/internal/api_to_rpc"
	api "Open_IM/pkg/api_struct"
	"Open_IM/pkg/common/config"
	rpc "Open_IM/pkg/proto/friend"
	"github.com/gin-gonic/gin"
)

func AddBlack(c *gin.Context) {
	common.ApiToRpc(c, &api.AddBlacklistReq{}, &api.AddBlacklistResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func ImportFriend(c *gin.Context) {
	common.ApiToRpc(c, &api.ImportFriendReq{}, &api.ImportFriendResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func AddFriend(c *gin.Context) {
	common.ApiToRpc(c, &api.AddFriendReq{}, &api.AddFriendResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func AddFriendResponse(c *gin.Context) {
	common.ApiToRpc(c, &api.AddFriendResponseReq{}, &api.AddFriendResponseResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func DeleteFriend(c *gin.Context) {
	common.ApiToRpc(c, &api.DeleteFriendReq{}, &api.DeleteFriendResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func GetBlacklist(c *gin.Context) {
	common.ApiToRpc(c, &api.GetBlackListReq{}, &api.GetBlackListResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func SetFriendRemark(c *gin.Context) {
	common.ApiToRpc(c, &api.SetFriendRemarkReq{}, &api.SetFriendRemarkResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func RemoveBlacklist(c *gin.Context) {
	common.ApiToRpc(c, &api.RemoveBlacklistReq{}, &api.RemoveBlacklistResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func IsFriend(c *gin.Context) {
	common.ApiToRpc(c, &api.IsFriendReq{}, &api.IsFriendResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func GetFriendList(c *gin.Context) {
	common.ApiToRpc(c, &api.GetFriendListReq{}, &api.GetFriendListResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func GetFriendApplyList(c *gin.Context) {
	common.ApiToRpc(c, &api.GetFriendApplyListReq{}, &api.GetFriendApplyListResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

func GetSelfApplyList(c *gin.Context) {
	common.ApiToRpc(c, &api.GetSelfApplyListReq{}, &api.GetSelfApplyListResp{}, config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, "")
}

package friend

import (
	common "Open_IM/internal/api_to_rpc"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/token_verify"
	rpc "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 不一致
func AddBlacklist(c *gin.Context) {
	common.ApiToRpc(c, &api.AddBlacklistReq{}, &api.AddBlacklistResp{},
		config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, utils.GetSelfFuncName(), token_verify.ParseUserIDFromToken)
}

func ImportFriend1(c *gin.Context) {
	common.ApiToRpc(c, &api.ImportFriendReq{}, &api.ImportFriendResp{},
		config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, utils.GetSelfFuncName(), token_verify.ParseUserIDFromToken)
}

func AddFriend1(c *gin.Context) {
	common.ApiToRpc(c, &api.AddFriendReq{}, &api.AddFriendResp{},
		config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, utils.GetSelfFuncName(), token_verify.ParseUserIDFromToken)
}

func AddFriendResponse1(c *gin.Context) {
	common.ApiToRpc(c, &api.AddFriendResponseReq{}, &api.AddFriendResponseResp{},
		config.Config.RpcRegisterName.OpenImFriendName, rpc.NewFriendClient, utils.GetSelfFuncName(), token_verify.ParseUserIDFromToken)
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

type UserGlobalBlackApi struct {
	blacklistDB    controller.UserGlobalBlackDatabase
	userDB         database.User
	imAdminUserIDs []string
	authClient     *rpcli.AuthClient
}

func NewUserGlobalBlackApi(blacklistDB controller.UserGlobalBlackDatabase, userDB database.User, imAdminUserIDs []string, authClient *rpcli.AuthClient) UserGlobalBlackApi {
	return UserGlobalBlackApi{blacklistDB: blacklistDB, userDB: userDB, imAdminUserIDs: imAdminUserIDs, authClient: authClient}
}

type addGlobalBlacklistReq struct {
	Nicknames []string `json:"nicknames" binding:"required,min=1"`
	Reason    string   `json:"reason"`
}

type removeGlobalBlacklistReq struct {
	Nicknames []string `json:"nicknames" binding:"required,min=1"`
}

type getGlobalBlacklistReq struct {
	Pagination *sdkws.RequestPagination `json:"pagination" binding:"required"`
}

type globalBlackItem struct {
	UserID     string `json:"userID"`
	Nickname   string `json:"nickname"`
	OperatorID string `json:"operatorID"`
	Reason     string `json:"reason"`
	CreateTime int64  `json:"createTime"`
}

type getGlobalBlacklistResp struct {
	Total  int64             `json:"total"`
	Blacks []globalBlackItem `json:"blacks"`
}

// AddGlobalBlacklist 管理员将用户加入全局黑名单，并立即踢下线（所有平台 token 标记 KickedToken）
func (b *UserGlobalBlackApi) AddGlobalBlacklist(c *gin.Context) {
	var req addGlobalBlacklistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if err := authverify.CheckAdmin(c, b.imAdminUserIDs); err != nil {
		apiresp.GinError(c, err)
		return
	}
	operatorID := mcontext.GetOpUserID(c)
	blacks := make([]*model.UserGlobalBlack, 0, len(req.Nicknames))
	for _, nickname := range req.Nicknames {
		users, err := b.userDB.TakeByNickname(c, nickname)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		if len(users) == 0 {
			apiresp.GinError(c, errs.ErrRecordNotFound.WrapMsg("nickname not found", "nickname", nickname))
			return
		}
		if len(users) > 1 {
			apiresp.GinError(c, errs.ErrArgs.WrapMsg("nickname matched multiple users", "nickname", nickname))
			return
		}
		blacks = append(blacks, &model.UserGlobalBlack{
			UserID:     users[0].UserID,
			Nickname:   users[0].Nickname,
			OperatorID: operatorID,
			Reason:     req.Reason,
		})
	}
	if err := b.blacklistDB.AddBlack(c, blacks); err != nil {
		apiresp.GinError(c, err)
		return
	}
	// 黑名单写入成功后，对每个被封禁用户的所有非管理员平台执行 force_logout：
	// 1. 断开 WS 长连接（msggateway.KickUserOffline）
	// 2. 将 Redis 中该平台的所有 token 标记为 KickedToken
	for _, black := range blacks {
		for platformID := range constant.PlatformID2Name {
			if int32(platformID) == constant.AdminPlatformID {
				continue
			}
			if err := b.authClient.ForceLogout(c, black.UserID, int32(platformID)); err != nil {
				// 踢下线失败不阻断主流程，记录警告即可
				log.ZWarn(c, "AddGlobalBlacklist: ForceLogout failed", err,
					"userID", black.UserID, "platformID", platformID)
			}
		}
	}
	apiresp.GinSuccess(c, nil)
}

// RemoveGlobalBlacklist 管理员从全局黑名单移除用户
func (b *UserGlobalBlackApi) RemoveGlobalBlacklist(c *gin.Context) {
	var req removeGlobalBlacklistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if err := authverify.CheckAdmin(c, b.imAdminUserIDs); err != nil {
		apiresp.GinError(c, err)
		return
	}
	userIDs := make([]string, 0, len(req.Nicknames))
	for _, nickname := range req.Nicknames {
		users, err := b.userDB.TakeByNickname(c, nickname)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		if len(users) == 0 {
			apiresp.GinError(c, errs.ErrRecordNotFound.WrapMsg("nickname not found", "nickname", nickname))
			return
		}
		if len(users) > 1 {
			apiresp.GinError(c, errs.ErrArgs.WrapMsg("nickname matched multiple users", "nickname", nickname))
			return
		}
		userIDs = append(userIDs, users[0].UserID)
	}
	if err := b.blacklistDB.RemoveBlack(c, userIDs); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

// GetGlobalBlacklist 管理员分页查询全局黑名单
func (b *UserGlobalBlackApi) GetGlobalBlacklist(c *gin.Context) {
	var req getGlobalBlacklistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if err := authverify.CheckAdmin(c, b.imAdminUserIDs); err != nil {
		apiresp.GinError(c, err)
		return
	}
	total, blacks, err := b.blacklistDB.GetBlackList(c, req.Pagination)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	items := make([]globalBlackItem, 0, len(blacks))
	for _, blk := range blacks {
		items = append(items, globalBlackItem{
			UserID:     blk.UserID,
			Nickname:   blk.Nickname,
			OperatorID: blk.OperatorID,
			Reason:     blk.Reason,
			CreateTime: blk.CreateTime.UnixMilli(),
		})
	}
	apiresp.GinSuccess(c, getGlobalBlacklistResp{Total: total, Blacks: items})
}

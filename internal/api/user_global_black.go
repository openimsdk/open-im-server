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
	UserIDs []string `json:"userIDs" binding:"required,min=1"`
	Reason  string   `json:"reason"`
	// Status 限制类型：1=冻结（可登录，不能收发消息）；2=黑名单（不可登录，自动踢下线）
	Status int32 `json:"status" binding:"required,oneof=1 2"`
}

type removeGlobalBlacklistReq struct {
	UserIDs []string `json:"userIDs" binding:"required,min=1"`
	// Status 目标状态：0=恢复正常（同步从 blacklistDB 删除记录）；1=冻结；2=黑名单
	Status int32 `json:"status" binding:"oneof=0 1 2"`
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
	// Status 限制类型：1=冻结，2=黑名单
	Status int32 `json:"status"`
}

type getGlobalBlacklistResp struct {
	Total  int64             `json:"total"`
	Blacks []globalBlackItem `json:"blacks"`
}

// AddGlobalBlacklist 管理员设置用户限制状态。
// Status=1（冻结）：可登录，但不能收发消息；Status=2（黑名单）：不可登录，自动踢下线，不能收发消息。
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
	foundUsers, err := b.userDB.Find(c, req.UserIDs)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	userMap := make(map[string]*model.User, len(foundUsers))
	for _, u := range foundUsers {
		userMap[u.UserID] = u
	}
	blacks := make([]*model.UserGlobalBlack, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		u, ok := userMap[userID]
		if !ok {
			apiresp.GinError(c, errs.ErrRecordNotFound.WrapMsg("userID not found", "userID", userID))
			return
		}
		blacks = append(blacks, &model.UserGlobalBlack{
			UserID:     u.UserID,
			Nickname:   u.Nickname,
			OperatorID: operatorID,
			Reason:     req.Reason,
			Status:     req.Status,
		})
	}
	if err := b.blacklistDB.AddBlack(c, blacks); err != nil {
		apiresp.GinError(c, err)
		return
	}
	// 同步更新 user 集合中的状态字段
	for _, userID := range req.UserIDs {
		if err := b.userDB.UpdateByMap(c, userID, map[string]any{"status": req.Status}); err != nil {
			log.ZWarn(c, "AddGlobalBlacklist: UpdateByMap status failed", err,
				"userID", userID, "status", req.Status)
		}
	}
	// 仅黑名单（Status=2）需要踢下线：断开 WS 长连接并将 token 标记为 KickedToken
	if req.Status == model.UserStatusBlacklist {
		for _, black := range blacks {
			for platformID := range constant.PlatformID2Name {
				if int32(platformID) == constant.AdminPlatformID {
					continue
				}
				if err := b.authClient.ForceLogout(c, black.UserID, int32(platformID)); err != nil {
					log.ZWarn(c, "AddGlobalBlacklist: ForceLogout failed", err,
						"userID", black.UserID, "platformID", platformID)
				}
			}
		}
	}
	apiresp.GinSuccess(c, nil)
}

// RemoveGlobalBlacklist 管理员更新用户账号状态。
// 执行顺序：
//  1. 将 user 集合中的 status 字段更新为请求值
//  2. 仅当 status == 0（恢复正常）时，才从 blacklistDB 删除该用户的限制记录
//
// 说明：blacklistDB 是 auth/msg 层的拦截依据；状态先落 user 集合，
// 只有确认目标状态为"正常"时才清除黑名单记录，避免状态写入成功但记录未删导致仍被拦截。
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
	for _, userID := range req.UserIDs {
		if err := b.userDB.UpdateByMap(c, userID, map[string]any{"status": req.Status}); err != nil {
			log.ZError(c, "RemoveGlobalBlacklist: UpdateByMap status failed", err, "userID", userID, "status", req.Status)
			apiresp.GinError(c, err)
			return
		}
	}
	// 只有目标状态为 0（正常）时才删除 blacklistDB 中的限制记录
	if req.Status == model.UserStatusNormal {
		if err := b.blacklistDB.RemoveBlack(c, req.UserIDs); err != nil {
			apiresp.GinError(c, err)
			return
		}
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
			Status:     blk.Status,
		})
	}
	apiresp.GinSuccess(c, getGlobalBlacklistResp{Total: total, Blacks: items})
}

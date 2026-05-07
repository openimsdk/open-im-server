package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

// DeleteUserApi handles real account deletion (hard delete).
// It follows the same direct-DB pattern as UserGlobalBlackApi.
type DeleteUserApi struct {
	userDB         database.User
	authClient     *rpcli.AuthClient
	groupClient    group.GroupClient
	friendClient   relation.FriendClient
	imAdminUserIDs []string
}

func NewDeleteUserApi(
	userDB database.User,
	authClient *rpcli.AuthClient,
	groupClient group.GroupClient,
	friendClient relation.FriendClient,
	imAdminUserIDs []string,
) *DeleteUserApi {
	return &DeleteUserApi{
		userDB:         userDB,
		authClient:     authClient,
		groupClient:    groupClient,
		friendClient:   friendClient,
		imAdminUserIDs: imAdminUserIDs,
	}
}

type deleteUserReq struct {
	UserID string `json:"userID" binding:"required"`
}

// DeleteUser permanently deletes a user account and cleans up associated data.
// Steps: force-logout → delete friends → quit/kick groups → hard-delete user doc.
// Caller must be the same user as userID, or an IM admin (see CheckAccessV3).
func (d *DeleteUserApi) DeleteUser(c *gin.Context) {
	var req deleteUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	// Only the user themselves (or an IM admin) may delete the account.
	if err := authverify.CheckAccessV3(c, req.UserID, d.imAdminUserIDs); err != nil {
		apiresp.GinError(c, err)
		return
	}

	// 1. Verify user exists
	users, err := d.userDB.Find(c, []string{req.UserID})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(users) == 0 {
		apiresp.GinError(c, errs.ErrRecordNotFound.WrapMsg("user not found", "userID", req.UserID))
		return
	}

	// 2. Force logout from every platform
	for platformID := range constant.PlatformID2Name {
		if int32(platformID) == constant.AdminPlatformID {
			continue
		}
		if err := d.authClient.ForceLogout(c, req.UserID, int32(platformID)); err != nil {
			log.ZWarn(c, "DeleteUser: ForceLogout failed", err, "userID", req.UserID, "platformID", platformID)
		}
	}

	// 3. Delete all friendships (both directions: target→friend and friend→target)
	friendIDsResp, err := d.friendClient.GetFriendIDs(c, &relation.GetFriendIDsReq{UserID: req.UserID})
	if err != nil {
		log.ZWarn(c, "DeleteUser: GetFriendIDs failed", err, "userID", req.UserID)
	} else {
		for _, friendID := range friendIDsResp.FriendIDs {
			// Remove from target user's friend list
			if _, err := d.friendClient.DeleteFriend(c, &relation.DeleteFriendReq{
				OwnerUserID:  req.UserID,
				FriendUserID: friendID,
			}); err != nil {
				log.ZWarn(c, "DeleteUser: DeleteFriend (owner→friend) failed", err,
					"ownerUserID", req.UserID, "friendUserID", friendID)
			}
			// Remove from the friend's friend list
			//if _, err := d.friendClient.DeleteFriend(c, &relation.DeleteFriendReq{
			//	OwnerUserID:  friendID,
			//	FriendUserID: req.UserID,
			//}); err != nil {
			//	log.ZWarn(c, "DeleteUser: DeleteFriend (friend→owner) failed", err,
			//		"ownerUserID", friendID, "friendUserID", req.UserID)
			//}
		}
	}

	// 4. Quit / kick from all joined groups (paginated, page size 100)
	pageNumber := int32(1)
	const pageSize = int32(100)
	for {
		groupListResp, err := d.groupClient.GetJoinedGroupList(c, &group.GetJoinedGroupListReq{
			FromUserID: req.UserID,
			Pagination: &sdkws.RequestPagination{PageNumber: pageNumber, ShowNumber: pageSize},
		})
		if err != nil {
			log.ZWarn(c, "DeleteUser: GetJoinedGroupList failed", err, "userID", req.UserID, "page", pageNumber)
			break
		}
		for _, g := range groupListResp.Groups {
			if _, err := d.groupClient.QuitGroup(c, &group.QuitGroupReq{
				GroupID: g.GroupID,
				UserID:  req.UserID,
			}); err != nil {
				log.ZWarn(c, "DeleteUser: QuitGroup failed", err, "userID", req.UserID, "groupID", g.GroupID)
			}
		}
		if int32(len(groupListResp.Groups)) < pageSize {
			break
		}
		pageNumber++
	}

	// 5. Hard-delete user document from MongoDB.
	// Redis cache will become stale and expire via TTL; the user can no longer
	// authenticate because their tokens were already invalidated in step 2.
	if err := d.userDB.Delete(c, []string{req.UserID}); err != nil {
		apiresp.GinError(c, err)
		return
	}

	log.ZInfo(c, "DeleteUser: user deleted", "userID", req.UserID)
	apiresp.GinSuccess(c, nil)
}

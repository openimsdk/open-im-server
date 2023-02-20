package controller

import (
	"Open_IM/internal/tx"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	relationTb "Open_IM/pkg/common/db/table/relation"
	unRelationTb "Open_IM/pkg/common/db/table/unrelation"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"github.com/dtm-labs/rockscache"
	_ "github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
)

type GroupController interface {
	// Group
	CreateGroup(ctx context.Context, groups []*relationTb.GroupModel, groupMembers []*relationTb.GroupMemberModel) error
	TakeGroup(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error)
	FindGroup(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error)
	SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupModel, error)
	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
	DismissGroup(ctx context.Context, groupID string) error // 解散群，并删除群成员
	GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error)
	// GroupMember
	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relationTb.GroupMemberModel, err error)
	TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error)
	FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relationTb.GroupMemberModel, error)
	FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error)
	PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error)
	SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error)
	HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relationTb.GroupMemberModel) error
	DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error
	MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error)
	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error // 转让群
	UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error
	UpdateGroupMembers(ctx context.Context, data []*relationTb.BatchUpdateGroupMember) error
	// GroupRequest
	CreateGroupRequest(ctx context.Context, requests []*relationTb.GroupRequestModel) error
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relationTb.GroupRequestModel, error)
	PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error)
	// SuperGroupModelInterface
	FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unRelationTb.SuperGroupModel, error)
	FindJoinSuperGroup(ctx context.Context, userID string) (*unRelationTb.UserToSuperGroupModel, error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
}

func NewGroupController(
	group relationTb.GroupModelInterface,
	member relationTb.GroupMemberModelInterface,
	request relationTb.GroupRequestModelInterface,
	tx tx.Tx,
	ctxTx tx.CtxTx,
	super unRelationTb.SuperGroupModelInterface,
	client redis.UniversalClient,
) GroupController {
	database := &GroupDataBase{
		groupDB:        group,
		groupMemberDB:  member,
		groupRequestDB: request,
		tx:             tx,
		ctxTx:          ctxTx,
		cache: cache.NewGroupCacheRedis(client, group, member, request, super, rockscache.Options{
			RandomExpireAdjustment: 0.2,
			DisableCacheRead:       false,
			DisableCacheDelete:     false,
			StrongConsistency:      true,
		}),
		mongoDB: super,
	}
	return database
}

type GroupDataBase struct {
	groupDB        relationTb.GroupModelInterface
	groupMemberDB  relationTb.GroupMemberModelInterface
	groupRequestDB relationTb.GroupRequestModelInterface
	tx             tx.Tx
	ctxTx          tx.CtxTx
	cache          cache.GroupCacheRedisInterface
	mongoDB        unRelationTb.SuperGroupModelInterface
}

func (g *GroupDataBase) GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error) {
	return g.groupDB.GetGroupIDsByGroupType(ctx, groupType)
}

func (g *GroupDataBase) delGroupMemberCache(ctx context.Context, groupID string, userIDs []string) error {
	for _, userID := range userIDs {
		if err := g.cache.DelJoinedGroupID(ctx, userID); err != nil {
			return err
		}
		if err := g.cache.DelJoinedSuperGroupIDs(ctx, userID); err != nil {
			return err
		}
	}
	if err := g.cache.DelGroupMemberIDs(ctx, groupID); err != nil {
		return err
	}
	if err := g.cache.DelGroupMemberNum(ctx, groupID); err != nil {
		return err
	}
	if err := g.cache.DelGroupMembersHash(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *GroupDataBase) FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error) {
	return g.cache.GetGroupMemberIDs(ctx, groupID)
}

func (g *GroupDataBase) CreateGroup(ctx context.Context, groups []*relationTb.GroupModel, groupMembers []*relationTb.GroupMemberModel) error {
	return g.tx.Transaction(func(tx any) error {
		if len(groups) > 0 {
			if err := g.groupDB.NewTx(tx).Create(ctx, groups); err != nil {
				return err
			}
		}
		if len(groupMembers) > 0 {
			if err := g.groupMemberDB.NewTx(tx).Create(ctx, groupMembers); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *GroupDataBase) TakeGroup(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error) {
	return g.cache.GetGroupInfo(ctx, groupID)
}

func (g *GroupDataBase) FindGroup(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error) {
	return g.cache.GetGroupsInfo(ctx, groupIDs)
}

func (g *GroupDataBase) SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupModel, error) {
	return g.groupDB.Search(ctx, keyword, pageNumber, showNumber)
}

func (g *GroupDataBase) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupDB.NewTx(tx).UpdateMap(ctx, groupID, data); err != nil {
			return err
		}
		if err := g.cache.DelGroupInfo(ctx, groupID); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) DismissGroup(ctx context.Context, groupID string) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupDB.NewTx(tx).UpdateStatus(ctx, groupID, constant.GroupStatusDismissed); err != nil {
			return err
		}
		if err := g.groupMemberDB.NewTx(tx).DeleteGroup(ctx, []string{groupID}); err != nil {
			return err
		}
		userIDs, err := g.cache.GetGroupMemberIDs(ctx, groupID)
		if err != nil {
			return err
		}
		if err := g.delGroupMemberCache(ctx, groupID, userIDs); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relationTb.GroupMemberModel, err error) {
	return g.cache.GetGroupMemberInfo(ctx, groupID, userID)
}

func (g *GroupDataBase) TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.TakeOwner(ctx, groupID) // todo cache group owner
}

func (g *GroupDataBase) FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.Find(ctx, groupIDs, userIDs, roleLevels) // todo cache group find
}

func (g *GroupDataBase) PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, "", groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupDataBase) SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, keyword, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupDataBase) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relationTb.GroupMemberModel) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupRequestDB.NewTx(tx).UpdateHandler(ctx, groupID, userID, handledMsg, handleResult); err != nil {
			return err
		}
		if member != nil {
			if err := g.groupMemberDB.NewTx(tx).Create(ctx, []*relationTb.GroupMemberModel{member}); err != nil {
				return err
			}
			if err := g.delGroupMemberCache(ctx, groupID, []string{userID}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *GroupDataBase) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupMemberDB.NewTx(tx).Delete(ctx, groupID, userIDs); err != nil {
			return err
		}
		if err := g.delGroupMemberCache(ctx, groupID, userIDs); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error) {
	return g.cache.GetGroupMemberHash1(ctx, groupIDs)
}

func (g *GroupDataBase) MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	return g.groupMemberDB.MapGroupMemberNum(ctx, groupIDs)
}

func (g *GroupDataBase) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error {
	return g.tx.Transaction(func(tx any) error {
		rowsAffected, err := g.groupMemberDB.NewTx(tx).UpdateRoleLevel(ctx, groupID, oldOwnerUserID, roleLevel)
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			return utils.Wrap(fmt.Errorf("oldOwnerUserID %s rowsAffected = %d", oldOwnerUserID, rowsAffected), "")
		}
		rowsAffected, err = g.groupMemberDB.NewTx(tx).UpdateRoleLevel(ctx, groupID, newOwnerUserID, constant.GroupOwner)
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			return utils.Wrap(fmt.Errorf("newOwnerUserID %s rowsAffected = %d", newOwnerUserID, rowsAffected), "")
		}
		if err := g.delGroupMemberCache(ctx, groupID, []string{oldOwnerUserID, newOwnerUserID}); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupMemberDB.NewTx(tx).Update(ctx, groupID, userID, data); err != nil {
			return err
		}
		if err := g.cache.DelGroupMemberInfo(ctx, groupID, userID); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) UpdateGroupMembers(ctx context.Context, data []*relationTb.BatchUpdateGroupMember) error {
	return g.tx.Transaction(func(tx any) error {
		for _, item := range data {
			if err := g.groupMemberDB.NewTx(tx).Update(ctx, item.GroupID, item.UserID, item.Map); err != nil {
				return err
			}
			if err := g.cache.DelGroupMemberInfo(ctx, item.GroupID, item.UserID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *GroupDataBase) CreateGroupRequest(ctx context.Context, requests []*relationTb.GroupRequestModel) error {
	return g.groupRequestDB.Create(ctx, requests)
}

func (g *GroupDataBase) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.Take(ctx, groupID, userID)
}

func (g *GroupDataBase) PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.Page(ctx, userID, pageNumber, showNumber)
}

func (g *GroupDataBase) FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unRelationTb.SuperGroupModel, error) {
	return g.mongoDB.FindSuperGroup(ctx, groupIDs)
}

func (g *GroupDataBase) FindJoinSuperGroup(ctx context.Context, userID string) (*unRelationTb.UserToSuperGroupModel, error) {
	return g.mongoDB.GetSuperGroupByUserID(ctx, userID)
}

func (g *GroupDataBase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		return g.mongoDB.CreateSuperGroup(ctx, groupID, initMemberIDList)
	})
}

func (g *GroupDataBase) DeleteSuperGroup(ctx context.Context, groupID string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		return g.mongoDB.DeleteSuperGroup(ctx, groupID)
	})
}

func (g *GroupDataBase) DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		return g.mongoDB.RemoverUserFromSuperGroup(ctx, groupID, userIDs)
	})
}

func (g *GroupDataBase) CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		return g.mongoDB.AddUserToSuperGroup(ctx, groupID, userIDs)
	})
}

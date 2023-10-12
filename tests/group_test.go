package controller_test

import (
	"context"
	"testing"

	"github.com/OpenIMSDK/open-im-server/v3/pkg/common/db/controller"
	"github.com/OpenIMSDK/open-im-server/v3/pkg/common/db/relation"
	"github.com/OpenIMSDK/open-im-server/v3/pkg/common/db/table/relationtb"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func TestDeleteGroupMemberHash(t *testing.T) {
	db := gorm.DB{}
	rdb := redis.UniversalClient{}
	database := mongo.Database{}
	hashCode := func(ctx context.Context, groupID string) (uint64, error) { return 0, nil }

	groupDatabase := controller.InitGroupDatabase(&db, rdb, &database, hashCode)

	group := &relationtb.GroupModel{GroupID: "testGroupID"}
	member := &relationtb.GroupMemberModel{GroupID: "testGroupID", UserID: "testUserID"}

	groupDatabase.CreateGroup(context.Background(), []*relationtb.GroupModel{group}, []*relationtb.GroupMemberModel{member})

	err := groupDatabase.DeleteGroupMemberHash(context.Background(), "testGroupID", "testUserID")
	if err != nil {
		t.Errorf("Failed to delete group member hash: %v", err)
	}

	members, err := groupDatabase.FindGroupMember(context.Background(), []string{"testGroupID"}, []string{"testUserID"}, nil)
	if err != nil {
		t.Errorf("Failed to find group member: %v", err)
	}

	if len(members) != 0 {
		t.Errorf("Group member hash was not deleted")
	}
}

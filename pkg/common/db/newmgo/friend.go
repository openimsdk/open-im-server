package newmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo/mgotool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FriendMgo implements FriendModelInterface using MongoDB as the storage backend.
type FriendMgo struct {
	coll *mongo.Collection
}

// NewFriendMongo creates a new instance of FriendMgo with the provided MongoDB database.
func NewFriendMongo(db *mongo.Database) (relation.FriendModelInterface, error) {
	return &FriendMgo{
		coll: db.Collection("friend"),
	}, nil
}

// Create inserts multiple friend records.
func (f *FriendMgo) Create(ctx context.Context, friends []*relation.FriendModel) error {
	return mgotool.InsertMany(ctx, f.coll, friends)
}

// Delete removes specified friends of the owner user.
func (f *FriendMgo) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) error {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": bson.M{"$in": friendUserIDs},
	}
	return mgotool.DeleteOne(ctx, f.coll, filter)
}

// UpdateByMap updates specific fields of a friend document using a map.
func (f *FriendMgo) UpdateByMap(ctx context.Context, ownerUserID string, friendUserID string, args map[string]interface{}) error {
	if len(args) == 0 {
		return nil
	}
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": friendUserID,
	}
	return mgotool.UpdateOne(ctx, f.coll, filter, bson.M{"$set": args}, true)
}

// Update modifies multiple friend documents.
// func (f *FriendMgo) Update(ctx context.Context, friends []*relation.FriendModel) error {
// 	filter := bson.M{
// 		"owner_user_id":  ownerUserID,
// 		"friend_user_id": friendUserID,
// 	}
// 	return mgotool.UpdateMany(ctx, f.coll, filter, friends)
// }

// UpdateRemark updates the remark for a specific friend.
func (f *FriendMgo) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) error {
	return f.UpdateByMap(ctx, ownerUserID, friendUserID, map[string]any{"remark": remark})
}

// Take retrieves a single friend document. Returns an error if not found.
func (f *FriendMgo) Take(ctx context.Context, ownerUserID, friendUserID string) (*relation.FriendModel, error) {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": friendUserID,
	}
	return mgotool.FindOne[*relation.FriendModel](ctx, f.coll, filter)
}

// FindUserState finds the friendship status between two users.
func (f *FriendMgo) FindUserState(ctx context.Context, userID1, userID2 string) ([]*relation.FriendModel, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"owner_user_id": userID1, "friend_user_id": userID2},
			{"owner_user_id": userID2, "friend_user_id": userID1},
		},
	}
	return mgotool.Find[*relation.FriendModel](ctx, f.coll, filter)
}

// FindFriends retrieves a list of friends for a given owner. Missing friends do not cause an error.
func (f *FriendMgo) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) ([]*relation.FriendModel, error) {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": bson.M{"$in": friendUserIDs},
	}
	return mgotool.Find[*relation.FriendModel](ctx, f.coll, filter)
}

// FindReversalFriends finds users who have added the specified user as a friend.
func (f *FriendMgo) FindReversalFriends(ctx context.Context, friendUserID string, ownerUserIDs []string) ([]*relation.FriendModel, error) {
	filter := bson.M{
		"owner_user_id":  bson.M{"$in": ownerUserIDs},
		"friend_user_id": friendUserID,
	}
	return mgotool.Find[*relation.FriendModel](ctx, f.coll, filter)
}

// FindOwnerFriends retrieves a paginated list of friends for a given owner.
func (f *FriendMgo) FindOwnerFriends(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (int64, []*relation.FriendModel, error) {
	filter := bson.M{"owner_user_id": ownerUserID}
	return mgotool.FindPage[*relation.FriendModel](ctx, f.coll, filter, pagination)
}

// FindInWhoseFriends finds users who have added the specified user as a friend, with pagination.
func (f *FriendMgo) FindInWhoseFriends(ctx context.Context, friendUserID string, pagination pagination.Pagination) (int64, []*relation.FriendModel, error) {
	filter := bson.M{"friend_user_id": friendUserID}
	return mgotool.FindPage[*relation.FriendModel](ctx, f.coll, filter, pagination)
}

// FindFriendUserIDs retrieves a list of friend user IDs for a given owner.
func (f *FriendMgo) FindFriendUserIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	filter := bson.M{"owner_user_id": ownerUserID}
	return mgotool.Find[string](ctx, f.coll, filter, options.Find().SetProjection(bson.M{"_id": 0, "friend_user_id": 1}))
}

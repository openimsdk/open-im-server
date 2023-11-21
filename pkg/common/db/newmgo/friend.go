package newmgo

import (
	"context"

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
		coll: db.Collection(relation.FriendModelCollectionName),
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
	_, err := f.coll.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
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
	update := bson.M{"$set": args}
	err := mgotool.UpdateOne(ctx, f.coll, filter, update, true)
	if err != nil {
		return err
	}
	return nil
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
	args := map[string]interface{}{"remark": remark}
	return f.UpdateByMap(ctx, ownerUserID, friendUserID, args)
}

// Take retrieves a single friend document. Returns an error if not found.
func (f *FriendMgo) Take(ctx context.Context, ownerUserID, friendUserID string) (*relation.FriendModel, error) {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": friendUserID,
	}
	friend, err := mgotool.FindOne[*relation.FriendModel](ctx, f.coll, filter)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

// FindUserState finds the friendship status between two users.
func (f *FriendMgo) FindUserState(ctx context.Context, userID1, userID2 string) ([]*relation.FriendModel, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"owner_user_id": userID1, "friend_user_id": userID2},
			{"owner_user_id": userID2, "friend_user_id": userID1},
		},
	}
	friends, err := mgotool.Find[*relation.FriendModel](ctx, f.coll, filter)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// FindFriends retrieves a list of friends for a given owner. Missing friends do not cause an error.
func (f *FriendMgo) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) ([]*relation.FriendModel, error) {
	filter := bson.M{
		"owner_user_id":  ownerUserID,
		"friend_user_id": bson.M{"$in": friendUserIDs},
	}
	friends, err := mgotool.Find[*relation.FriendModel](ctx, f.coll, filter)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// FindReversalFriends finds users who have added the specified user as a friend.
func (f *FriendMgo) FindReversalFriends(ctx context.Context, friendUserID string, ownerUserIDs []string) ([]*relation.FriendModel, error) {
	filter := bson.M{
		"owner_user_id":  bson.M{"$in": ownerUserIDs},
		"friend_user_id": friendUserID,
	}
	friends, err := mgotool.Find[*relation.FriendModel](ctx, f.coll, filter)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

// FindOwnerFriends retrieves a paginated list of friends for a given owner.
func (f *FriendMgo) FindOwnerFriends(ctx context.Context, ownerUserID string, pagination pagination.Pagination, showNumber int32) ([]*relation.FriendModel, int64, error) {
	filter := bson.M{"owner_user_id": ownerUserID}
	count, friends, err := mgotool.FindPage[*relation.FriendModel](ctx, f.coll, filter, pagination)
	if err != nil {
		return nil, 0, err
	}
	return friends, count, nil
}

// FindInWhoseFriends finds users who have added the specified user as a friend, with pagination.
func (f *FriendMgo) FindInWhoseFriends(ctx context.Context, friendUserID string, pagination.Pagination, showNumber int32) ([]*relation.FriendModel, int64, error) {
	filter := bson.M{"friend_user_id": friendUserID}
	count, friends, err := mgotool.FindPage[*relation.FriendModel](ctx, f.coll, filter, pagination)
	if err != nil {
		return nil, 0, err
	}
	return friends, count, nil
}

// FindFriendUserIDs retrieves a list of friend user IDs for a given owner.
func (f *FriendMgo) FindFriendUserIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	filter := bson.M{"owner_user_id": ownerUserID}
	friends := []*relation.FriendModel{}
	friends, err := mgotool.Find[*relation.FriendModel](ctx, f.coll, filter)
	if err != nil {
		return nil, err
	}

	friendUserIDs := make([]string, len(friends))
	for i, friend := range friends {
		friendUserIDs[i] = friend.FriendUserID
	}
	return friendUserIDs, nil
}

// NewTx creates a new transaction.
func (f *FriendMgo) NewTx(tx any) relation.FriendModelInterface {
	panic("not implemented")
	return nil
}

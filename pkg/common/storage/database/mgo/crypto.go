package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ---- CryptoDevice ----

type CryptoDeviceMgo struct {
	coll *mongo.Collection
}

func NewCryptoDeviceMongo(db *mongo.Database) (database.CryptoDevice, error) {
	coll := db.Collection("crypto_device")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "device_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &CryptoDeviceMgo{coll: coll}, nil
}

func (m *CryptoDeviceMgo) Create(ctx context.Context, device *model.CryptoDevice) error {
	_, err := m.coll.InsertOne(ctx, device)
	return err
}

func (m *CryptoDeviceMgo) FindByUserID(ctx context.Context, userID string) ([]*model.CryptoDevice, error) {
	cursor, err := m.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var devices []*model.CryptoDevice
	if err := cursor.All(ctx, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func (m *CryptoDeviceMgo) FindByUserIDAndDeviceID(ctx context.Context, userID, deviceID string) (*model.CryptoDevice, error) {
	var device model.CryptoDevice
	err := m.coll.FindOne(ctx, bson.M{"user_id": userID, "device_id": deviceID}).Decode(&device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("crypto device not found", "userID", userID, "deviceID", deviceID)
		}
		return nil, err
	}
	return &device, nil
}

func (m *CryptoDeviceMgo) UpdateStatus(ctx context.Context, userID, deviceID, status string) error {
	result, err := m.coll.UpdateOne(ctx,
		bson.M{"user_id": userID, "device_id": deviceID},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errs.ErrRecordNotFound.WrapMsg("crypto device not found", "userID", userID, "deviceID", deviceID)
	}
	return nil
}

func (m *CryptoDeviceMgo) UpdateLastSeen(ctx context.Context, userID, deviceID string) error {
	result, err := m.coll.UpdateOne(ctx,
		bson.M{"user_id": userID, "device_id": deviceID},
		bson.M{"$set": bson.M{"last_seen_at": time.Now()}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errs.ErrRecordNotFound.WrapMsg("crypto device not found", "userID", userID, "deviceID", deviceID)
	}
	return nil
}

// ---- GroupKeyVersion ----

type GroupKeyVersionMgo struct {
	coll *mongo.Collection
}

func NewGroupKeyVersionMongo(db *mongo.Database) (database.GroupKeyVersion, error) {
	coll := db.Collection("group_key_version")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "group_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &GroupKeyVersionMgo{coll: coll}, nil
}

func (m *GroupKeyVersionMgo) Find(ctx context.Context, groupID string) (*model.GroupKeyVersion, error) {
	var v model.GroupKeyVersion
	err := m.coll.FindOne(ctx, bson.M{"group_id": groupID}).Decode(&v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("group key version not found", "groupID", groupID)
		}
		return nil, err
	}
	return &v, nil
}

func (m *GroupKeyVersionMgo) IncrVersion(ctx context.Context, groupID string) (int64, error) {
	var result model.GroupKeyVersion
	err := m.coll.FindOneAndUpdate(ctx,
		bson.M{"group_id": groupID},
		bson.M{"$inc": bson.M{"group_key_version": int64(1)}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.GroupKeyVersion, nil
}

// ---- GroupKeyEvent ----

type GroupKeyEventMgo struct {
	coll *mongo.Collection
}

const maxGroupKeyEventsPerQuery = 500

func NewGroupKeyEventMongo(db *mongo.Database) (database.GroupKeyEvent, error) {
	coll := db.Collection("group_key_event")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "group_id", Value: 1}, {Key: "group_key_version", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "event_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, err
	}
	return &GroupKeyEventMgo{coll: coll}, nil
}

func (m *GroupKeyEventMgo) Create(ctx context.Context, event *model.GroupKeyEvent) error {
	_, err := m.coll.InsertOne(ctx, event)
	return err
}

func (m *GroupKeyEventMgo) FindSinceVersion(ctx context.Context, groupID string, sinceVersion int64) ([]*model.GroupKeyEvent, error) {
	cursor, err := m.coll.Find(ctx, bson.M{
		"group_id":          groupID,
		"group_key_version": bson.M{"$gt": sinceVersion},
	}, options.Find().
		SetSort(bson.D{{Key: "group_key_version", Value: 1}}).
		SetLimit(maxGroupKeyEventsPerQuery))
	if err != nil {
		return nil, err
	}
	var events []*model.GroupKeyEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

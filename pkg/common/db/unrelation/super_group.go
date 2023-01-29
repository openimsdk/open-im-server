package unrelation

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
)

const (
	cSuperGroup       = "super_group"
	cUserToSuperGroup = "user_to_super_group"
)

type SuperGroupMgoDB struct {
	MgoClient                  *mongo.Client
	MgoDB                      *mongo.Database
	superGroupCollection       *mongo.Collection
	userToSuperGroupCollection *mongo.Collection
}

type SuperGroup struct {
	GroupID      string   `bson:"group_id" json:"groupID"`
	MemberIDList []string `bson:"member_id_list" json:"memberIDList"`
}

type UserToSuperGroup struct {
	UserID      string   `bson:"user_id" json:"userID"`
	GroupIDList []string `bson:"group_id_list" json:"groupIDList"`
}

func NewSuperGroupMgoDB(mgoClient *mongo.Client) *SuperGroupMgoDB {
	mgoDB := mgoClient.Database(config.Config.Mongo.DBDatabase)
	return &SuperGroupMgoDB{MgoDB: mgoDB, MgoClient: mgoClient, superGroupCollection: mgoDB.Collection(cSuperGroup), userToSuperGroupCollection: mgoDB.Collection(cUserToSuperGroup)}
}

func (db *SuperGroupMgoDB) CreateSuperGroup(sCtx mongo.SessionContext, groupID string, initMemberIDList []string, memberNumCount int) error {
	superGroup := SuperGroup{
		GroupID:      groupID,
		MemberIDList: initMemberIDList,
	}
	_, err := db.superGroupCollection.InsertOne(sCtx, superGroup)
	if err != nil {
		return err
	}
	upsert := true
	opts := &options.UpdateOptions{
		Upsert: &upsert,
	}
	for _, userID := range initMemberIDList {
		_, err = db.userToSuperGroupCollection.UpdateOne(sCtx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
		if err != nil {
			return err
		}
	}
	return nil

}

func (db *SuperGroupMgoDB) GetSuperGroup(ctx context.Context, groupID string) (*SuperGroup, error) {
	superGroup := SuperGroup{}
	err := db.superGroupCollection.FindOne(ctx, bson.M{"group_id": groupID}).Decode(&superGroup)
	return &superGroup, err
}

func (db *SuperGroupMgoDB) AddUserToSuperGroup(ctx context.Context, groupID string, userIDList []string) error {
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	return db.MgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
		_, err := db.superGroupCollection.UpdateOne(sCtx, bson.M{"group_id": groupID}, bson.M{"$addToSet": bson.M{"member_id_list": bson.M{"$each": userIDList}}})
		if err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		upsert := true
		opts := &options.UpdateOptions{
			Upsert: &upsert,
		}
		for _, userID := range userIDList {
			_, err = db.userToSuperGroupCollection.UpdateOne(sCtx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
			if err != nil {
				_ = sCtx.AbortTransaction(ctx)
				return utils.Wrap(err, "transaction failed")
			}
		}
		return sCtx.CommitTransaction(ctx)
	})
}

func (db *SuperGroupMgoDB) RemoverUserFromSuperGroup(ctx context.Context, groupID string, userIDList []string) error {
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	return db.MgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
		_, err := db.superGroupCollection.UpdateOne(sCtx, bson.M{"group_id": groupID}, bson.M{"$pull": bson.M{"member_id_list": bson.M{"$in": userIDList}}})
		if err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		err = db.RemoveGroupFromUser(sCtx, groupID, userIDList)
		if err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		return sCtx.CommitTransaction(ctx)
	})
}

func (db *SuperGroupMgoDB) GetSuperGroupByUserID(ctx context.Context, userID string) (*UserToSuperGroup, error) {
	var user UserToSuperGroup
	_ = db.userToSuperGroupCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	return &user, nil
}

func (db *SuperGroupMgoDB) DeleteSuperGroup(ctx context.Context, groupID string) error {
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	return db.MgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
		superGroup := &SuperGroup{}
		_, err := db.superGroupCollection.DeleteOne(sCtx, bson.M{"group_id": groupID})
		if err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		if err = db.RemoveGroupFromUser(sCtx, groupID, superGroup.MemberIDList); err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		return sCtx.CommitTransaction(ctx)
	})
}

func (db *SuperGroupMgoDB) RemoveGroupFromUser(sCtx context.Context, groupID string, userIDList []string) error {
	_, err := db.userToSuperGroupCollection.UpdateOne(sCtx, bson.M{"user_id": bson.M{"$in": userIDList}}, bson.M{"$pull": bson.M{"group_id_list": groupID}})
	return err
}

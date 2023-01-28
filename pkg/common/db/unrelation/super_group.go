package unrelation

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"time"
)

const (
	cSuperGroup       = "super_group"
	cUserToSuperGroup = "user_to_super_group"
)

type SuperGroupMgo struct {
	mgoDB                      *mongo.Database
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

func NewSuperGroupMgoDB(mgoDB *mongo.Database) *SuperGroupMgo {
	return &SuperGroupMgo{mgoDB: mgoDB, superGroupCollection: mgoDB.Collection(cSuperGroup), userToSuperGroupCollection: mgoDB.Collection(cUserToSuperGroup)}
}

func (db *SuperGroupMgo) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string, memberNumCount int) error {
	//ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	//c := db.mgoDB.Database(config.Config.Mongo.DBDatabase).Collection(cSuperGroup)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	return db.mgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
		err := sCtx.StartTransaction()
		if err != nil {
			return err
		}
		superGroup := SuperGroup{
			GroupID:      groupID,
			MemberIDList: initMemberIDList,
		}
		_, err = db.superGroupCollection.InsertOne(sCtx, superGroup)
		if err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		upsert := true
		opts := &options.UpdateOptions{
			Upsert: &upsert,
		}
		for _, userID := range initMemberIDList {
			_, err = db.userToSuperGroupCollection.UpdateOne(sCtx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
			if err != nil {
				_ = sCtx.AbortTransaction(ctx)
				return err
			}
		}
		return sCtx.CommitTransaction(context.Background())
	})
}

func (db *SuperGroupMgo) GetSuperGroup(ctx context.Context, groupID string) (*SuperGroup, error) {
	superGroup := SuperGroup{}
	err := db.superGroupCollection.FindOne(ctx, bson.M{"group_id": groupID}).Decode(&superGroup)
	return &superGroup, err
}

func (db *SuperGroupMgo) AddUserToSuperGroup(ctx context.Context, groupID string, userIDList []string) error {
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	return db.mgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
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
		return sCtx.CommitTransaction(context.Background())
	})
}

func (d *SuperGroupMgo) RemoverUserFromSuperGroup(groupID string, userIDList []string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSuperGroup)
	session, err := d.mongoClient.StartSession()
	if err != nil {
		return utils.Wrap(err, "start session failed")
	}
	defer session.EndSession(ctx)
	sCtx := mongo.NewSessionContext(ctx, session)
	_, err = c.UpdateOne(ctx, bson.M{"group_id": groupID}, bson.M{"$pull": bson.M{"member_id_list": bson.M{"$in": userIDList}}})
	if err != nil {
		_ = session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	err = d.RemoveGroupFromUser(ctx, sCtx, groupID, userIDList)
	if err != nil {
		_ = session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	_ = session.CommitTransaction(ctx)
	return err
}

func (d *SuperGroupMgo) GetSuperGroupByUserID(userID string) (UserToSuperGroup, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cUserToSuperGroup)
	var user UserToSuperGroup
	_ = c.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	return user, nil
}

func (d *SuperGroupMgo) DeleteSuperGroup(groupID string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSuperGroup)
	session, err := d.mongoClient.StartSession()
	if err != nil {
		return utils.Wrap(err, "start session failed")
	}
	defer session.EndSession(ctx)
	sCtx := mongo.NewSessionContext(ctx, session)
	superGroup := &SuperGroup{}
	result := c.FindOneAndDelete(sCtx, bson.M{"group_id": groupID})
	err = result.Decode(superGroup)
	if err != nil {
		session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	if err = d.RemoveGroupFromUser(ctx, sCtx, groupID, superGroup.MemberIDList); err != nil {
		session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	session.CommitTransaction(ctx)
	return nil
}

func (d *SuperGroupMgo) RemoveGroupFromUser(ctx, sCtx context.Context, groupID string, userIDList []string) error {
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cUserToSuperGroup)
	_, err := c.UpdateOne(sCtx, bson.M{"user_id": bson.M{"$in": userIDList}}, bson.M{"$pull": bson.M{"group_id_list": groupID}})
	if err != nil {
		return utils.Wrap(err, "UpdateOne transaction failed")
	}
	return err
}

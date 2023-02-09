package unrelation

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db/table/unrelation"
	"Open_IM/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
)

var _ unrelation.SuperGroupModelInterface = (*SuperGroupMongoDriver)(nil)

func NewSuperGroupMongoDriver(mgoClient *mongo.Client) *SuperGroupMongoDriver {
	mgoDB := mgoClient.Database(config.Config.Mongo.DBDatabase)
	return &SuperGroupMongoDriver{MgoDB: mgoDB, MgoClient: mgoClient, superGroupCollection: mgoDB.Collection(unrelation.CSuperGroup), userToSuperGroupCollection: mgoDB.Collection(unrelation.CUserToSuperGroup)}
}

type SuperGroupMongoDriver struct {
	MgoClient                  *mongo.Client
	MgoDB                      *mongo.Database
	superGroupCollection       *mongo.Collection
	userToSuperGroupCollection *mongo.Collection
}

//	func (s *SuperGroupMongoDriver) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDs []string, tx ...interface{}) error {
//			superGroup := unrelation.SuperGroupModel{
//				GroupID:   groupID,
//				MemberIDs: initMemberIDs,
//			}
//			coll := getTxCtx(s.superGroupCollection, tx)
//			_, err := coll.InsertOne(ctx, superGroup)
//			if err != nil {
//				return err
//			}
//			opts := &options.UpdateOptions{
//				Upsert: utils.ToPtr(true),
//			}
//			for _, userID := range initMemberIDs {
//				_, err = coll.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
//				if err != nil {
//					return err
//				}
//			}
//			return nil
//	}
//
//	func (s *SuperGroupMongoDriver) FindSuperGroup(ctx context.Context, groupIDs []string, tx ...interface{}) (groups []*unrelation.SuperGroupModel, err error) {
//		cursor, err := s.superGroupCollection.Find(ctx, bson.M{"group_id": bson.M{
//			"$in": groupIDs,
//		}})
//		if err != nil {
//			return nil, utils.Wrap(err, "")
//		}
//		defer cursor.Close(ctx)
//		if err := cursor.All(ctx, &groups); err != nil {
//			return nil, utils.Wrap(err, "")
//		}
//		return groups, nil
//	}
//
//	func (s *SuperGroupMongoDriver) AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string, tx ...interface{}) error {
//			opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
//			return s.MgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
//				_, err := s.superGroupCollection.UpdateOne(sCtx, bson.M{"group_id": groupID}, bson.M{"$addToSet": bson.M{"member_id_list": bson.M{"$each": userIDs}}})
//				if err != nil {
//					_ = sCtx.AbortTransaction(ctx)
//					return err
//				}
//				upsert := true
//				opts := &options.UpdateOptions{
//					Upsert: &upsert,
//				}
//				for _, userID := range userIDs {
//					_, err = s.userToSuperGroupCollection.UpdateOne(sCtx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
//					if err != nil {
//						_ = sCtx.AbortTransaction(ctx)
//						return utils.Wrap(err, "transaction failed")
//					}
//				}
//				return sCtx.CommitTransaction(ctx)
//			})
//	}
//
//	func (s *SuperGroupMongoDriver) RemoverUserFromSuperGroup(ctx context.Context, groupID string, userIDs []string, tx ...interface{}) error {
//			opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
//			return s.MgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
//				_, err := s.superGroupCollection.UpdateOne(sCtx, bson.M{"group_id": groupID}, bson.M{"$pull": bson.M{"member_id_list": bson.M{"$in": userIDs}}})
//				if err != nil {
//					_ = sCtx.AbortTransaction(ctx)
//					return err
//				}
//				err = s.RemoveGroupFromUser(sCtx, groupID, userIDs)
//				if err != nil {
//					_ = sCtx.AbortTransaction(ctx)
//					return err
//				}
//				return sCtx.CommitTransaction(ctx)
//			})
//	}
//
//	func (s *SuperGroupMongoDriver) GetSuperGroupByUserID(ctx context.Context, userID string, tx ...interface{}) (*unrelation.UserToSuperGroupModel, error) {
//		//TODO implement me
//		panic("implement me")
//	}
//
//	func (s *SuperGroupMongoDriver) DeleteSuperGroup(ctx context.Context, groupID string, tx ...interface{}) error {
//		//TODO implement me
//		panic("implement me")
//	}

func (s *SuperGroupMongoDriver) Transaction(ctx context.Context, fn func(s unrelation.SuperGroupModelInterface, tx any) error) error {
	sess, err := s.MgoClient.StartSession()
	if err != nil {
		return err
	}
	txCtx := mongo.NewSessionContext(ctx, sess)
	defer sess.EndSession(txCtx)
	if err := fn(s, txCtx); err != nil {
		_ = sess.AbortTransaction(txCtx)
		return err
	}
	return utils.Wrap(sess.CommitTransaction(txCtx), "")
}

func (s *SuperGroupMongoDriver) getTxCtx(ctx context.Context, tx []any) context.Context {
	if len(tx) > 0 {
		if ctx, ok := tx[0].(mongo.SessionContext); ok {
			return ctx
		}
	}
	return ctx
}

//func (s *SuperGroupMongoDriver) Transaction(ctx context.Context, fn func(ctx mongo.SessionContext) error) error {
//	sess, err := s.MgoClient.StartSession()
//	if err != nil {
//		return err
//	}
//	sCtx := mongo.NewSessionContext(ctx, sess)
//
//	defer sess.EndSession(sCtx)
//	if err := fn(sCtx); err != nil {
//		_ = sess.AbortTransaction(sCtx)
//		return err
//	}
//	return utils.Wrap(sess.CommitTransaction(sCtx), "")
//}

func (s *SuperGroupMongoDriver) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDs []string, tx ...any) error {
	ctx = s.getTxCtx(ctx, tx)
	_, err := s.superGroupCollection.InsertOne(ctx, &unrelation.SuperGroupModel{
		GroupID:   groupID,
		MemberIDs: initMemberIDs,
	})
	if err != nil {
		return err
	}
	for _, userID := range initMemberIDs {
		_, err = s.userToSuperGroupCollection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, &options.UpdateOptions{
			Upsert: utils.ToPtr(true),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SuperGroupMongoDriver) TakeSuperGroup(ctx context.Context, groupID string, tx ...any) (group *unrelation.SuperGroupModel, err error) {
	ctx = s.getTxCtx(ctx, tx)
	if err := s.superGroupCollection.FindOne(ctx, bson.M{"group_id": groupID}).Decode(&group); err != nil {
		return nil, utils.Wrap(err, "")
	}
	return group, nil
}

func (s *SuperGroupMongoDriver) FindSuperGroup(ctx context.Context, groupIDs []string, tx ...any) (groups []*unrelation.SuperGroupModel, err error) {
	ctx = s.getTxCtx(ctx, tx)
	cursor, err := s.superGroupCollection.Find(ctx, bson.M{"group_id": bson.M{
		"$in": groupIDs,
	}})
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, utils.Wrap(err, "")
	}
	return groups, nil
}

func (s *SuperGroupMongoDriver) AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string, tx ...any) error {
	ctx = s.getTxCtx(ctx, tx)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	return s.MgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
		_, err := s.superGroupCollection.UpdateOne(sCtx, bson.M{"group_id": groupID}, bson.M{"$addToSet": bson.M{"member_id_list": bson.M{"$each": userIDs}}})
		if err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		upsert := true
		opts := &options.UpdateOptions{
			Upsert: &upsert,
		}
		for _, userID := range userIDs {
			_, err = s.userToSuperGroupCollection.UpdateOne(sCtx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
			if err != nil {
				_ = sCtx.AbortTransaction(ctx)
				return utils.Wrap(err, "transaction failed")
			}
		}
		return sCtx.CommitTransaction(ctx)
	})
}

func (s *SuperGroupMongoDriver) RemoverUserFromSuperGroup(ctx context.Context, groupID string, userIDs []string, tx ...any) error {
	ctx = s.getTxCtx(ctx, tx)
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	return s.MgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
		_, err := s.superGroupCollection.UpdateOne(sCtx, bson.M{"group_id": groupID}, bson.M{"$pull": bson.M{"member_id_list": bson.M{"$in": userIDs}}})
		if err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		err = s.RemoveGroupFromUser(sCtx, groupID, userIDs)
		if err != nil {
			_ = sCtx.AbortTransaction(ctx)
			return err
		}
		return sCtx.CommitTransaction(ctx)
	})
}

func (s *SuperGroupMongoDriver) GetSuperGroupByUserID(ctx context.Context, userID string, tx ...any) (*unrelation.UserToSuperGroupModel, error) {
	ctx = s.getTxCtx(ctx, tx)
	var user unrelation.UserToSuperGroupModel
	err := s.userToSuperGroupCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	return &user, utils.Wrap(err, "")
}

func (s *SuperGroupMongoDriver) DeleteSuperGroup(ctx context.Context, groupID string, tx ...any) error {
	ctx = s.getTxCtx(ctx, tx)
	group, err := s.TakeSuperGroup(ctx, groupID, tx...)
	if err != nil {
		return err
	}
	if _, err := s.superGroupCollection.DeleteOne(ctx, bson.M{"group_id": groupID}); err != nil {
		return utils.Wrap(err, "")
	}
	return s.RemoveGroupFromUser(ctx, groupID, group.MemberIDs)
}

//func (s *SuperGroupMongoDriver) DeleteSuperGroup(ctx context.Context, groupID string, tx ...any) error {
//	ctx = s.getTxCtx(ctx, tx)
//	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
//	return s.MgoDB.Client().UseSessionWithOptions(ctx, opts, func(sCtx mongo.SessionContext) error {
//		superGroup := &unrelation.SuperGroupModel{}
//		_, err := s.superGroupCollection.DeleteOne(sCtx, bson.M{"group_id": groupID})
//		if err != nil {
//			_ = sCtx.AbortTransaction(ctx)
//			return err
//		}
//		if err = s.RemoveGroupFromUser(sCtx, groupID, superGroup.MemberIDs); err != nil {
//			_ = sCtx.AbortTransaction(ctx)
//			return err
//		}
//		return sCtx.CommitTransaction(ctx)
//	})
//}

func (s *SuperGroupMongoDriver) RemoveGroupFromUser(ctx context.Context, groupID string, userIDs []string, tx ...any) error {
	ctx = s.getTxCtx(ctx, tx)
	_, err := s.userToSuperGroupCollection.UpdateOne(ctx, bson.M{"user_id": bson.M{"$in": userIDs}}, bson.M{"$pull": bson.M{"group_id_list": groupID}})
	return utils.Wrap(err, "")
}

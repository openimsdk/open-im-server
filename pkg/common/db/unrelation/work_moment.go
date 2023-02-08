package unrelation

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/table/unrelation"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type WorkMomentMongoDriver struct {
	mgoDB                *mongo.Database
	WorkMomentCollection *mongo.Collection
}

func NewWorkMomentMongoDriver(mgoDB *mongo.Database) *WorkMomentMongoDriver {
	return &WorkMomentMongoDriver{mgoDB: mgoDB, WorkMomentCollection: mgoDB.Collection(unrelation.CWorkMoment)}
}

func (db *WorkMomentMongoDriver) CreateOneWorkMoment(ctx context.Context, workMoment *unrelation.WorkMoment) error {
	workMomentID := generateWorkMomentID(workMoment.UserID)
	workMoment.WorkMomentID = workMomentID
	workMoment.CreateTime = int32(time.Now().Unix())
	_, err := db.WorkMomentCollection.InsertOne(ctx, workMoment)
	return err
}

func (db *WorkMomentMongoDriver) DeleteOneWorkMoment(ctx context.Context, workMomentID string) error {
	_, err := db.WorkMomentCollection.DeleteOne(ctx, bson.M{"work_moment_id": workMomentID})
	return err
}

func (db *WorkMomentMongoDriver) DeleteComment(ctx context.Context, workMomentID, contentID, opUserID string) error {
	_, err := db.WorkMomentCollection.UpdateOne(ctx, bson.D{{"work_moment_id", workMomentID},
		{"$or", bson.A{
			bson.D{{"user_id", opUserID}},
			bson.D{{"comments", bson.M{"$elemMatch": bson.M{"user_id": opUserID}}}},
		},
		}}, bson.M{"$pull": bson.M{"comments": bson.M{"content_id": contentID}}})
	return err
}

func (db *WorkMomentMongoDriver) GetWorkMomentByID(ctx context.Context, workMomentID string) (*unrelation.WorkMoment, error) {
	workMoment := &unrelation.WorkMoment{}
	err := db.WorkMomentCollection.FindOne(ctx, bson.M{"work_moment_id": workMomentID}).Decode(workMoment)
	return workMoment, err
}

func (db *WorkMomentMongoDriver) LikeOneWorkMoment(ctx context.Context, likeUserID, userName, workMomentID string) (*unrelation.WorkMoment, bool, error) {
	workMoment, err := db.GetWorkMomentByID(ctx, workMomentID)
	if err != nil {
		return nil, false, err
	}
	var isAlreadyLike bool
	for i, user := range workMoment.LikeUserList {
		if likeUserID == user.UserID {
			isAlreadyLike = true
			workMoment.LikeUserList = append(workMoment.LikeUserList[0:i], workMoment.LikeUserList[i+1:]...)
		}
	}
	if !isAlreadyLike {
		workMoment.LikeUserList = append(workMoment.LikeUserList, &unrelation.CommonUserModel{UserID: likeUserID, UserName: userName})
	}
	_, err = db.WorkMomentCollection.UpdateOne(ctx, bson.M{"work_moment_id": workMomentID}, bson.M{"$set": bson.M{"like_user_list": workMoment.LikeUserList}})
	return workMoment, !isAlreadyLike, err
}

func (db *WorkMomentMongoDriver) CommentOneWorkMoment(ctx context.Context, comment *unrelation.Comment, workMomentID string) (unrelation.WorkMoment, error) {
	comment.ContentID = generateWorkMomentCommentID(workMomentID)
	var workMoment unrelation.WorkMoment
	err := db.WorkMomentCollection.FindOneAndUpdate(ctx, bson.M{"work_moment_id": workMomentID}, bson.M{"$push": bson.M{"comments": comment}}).Decode(&workMoment)
	return workMoment, err
}

func (db *WorkMomentMongoDriver) GetUserSelfWorkMoments(ctx context.Context, userID string, showNumber, pageNumber int32) ([]unrelation.WorkMoment, error) {
	var workMomentList []unrelation.WorkMoment
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"create_time": -1})
	result, err := db.WorkMomentCollection.Find(ctx, bson.M{"user_id": userID}, findOpts)
	if err != nil {
		return workMomentList, nil
	}
	err = result.All(ctx, &workMomentList)
	return workMomentList, err
}

func (db *WorkMomentMongoDriver) GetUserWorkMoments(ctx context.Context, opUserID, userID string, showNumber, pageNumber int32, friendIDList []string) ([]unrelation.WorkMoment, error) {
	var workMomentList []unrelation.WorkMoment
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"create_time": -1})
	result, err := db.WorkMomentCollection.Find(ctx, bson.D{ // 等价条件: select * from
		{"user_id", userID},
		{"$or", bson.A{
			bson.D{{"permission", constant.WorkMomentPermissionCantSee}, {"permission_user_id_list", bson.D{{"$nin", bson.A{opUserID}}}}},
			bson.D{{"permission", constant.WorkMomentPermissionCanSee}, {"permission_user_id_list", bson.D{{"$in", bson.A{opUserID}}}}},
			bson.D{{"permission", constant.WorkMomentPublic}},
		}},
	}, findOpts)
	if err != nil {
		return workMomentList, nil
	}
	err = result.All(ctx, &workMomentList)
	return workMomentList, err
}

func (db *WorkMomentMongoDriver) GetUserFriendWorkMoments(ctx context.Context, showNumber, pageNumber int32, userID string, friendIDList []string) ([]unrelation.WorkMoment, error) {
	var workMomentList []unrelation.WorkMoment
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"create_time": -1})
	var filter bson.D
	permissionFilter := bson.D{
		{"$or", bson.A{
			bson.D{{"permission", constant.WorkMomentPermissionCantSee}, {"permission_user_id_list", bson.D{{"$nin", bson.A{userID}}}}},
			bson.D{{"permission", constant.WorkMomentPermissionCanSee}, {"permission_user_id_list", bson.D{{"$in", bson.A{userID}}}}},
			bson.D{{"permission", constant.WorkMomentPublic}},
		}}}
	if config.Config.WorkMoment.OnlyFriendCanSee {
		filter = bson.D{
			{"$or", bson.A{
				bson.D{{"user_id", userID}}, //self
				bson.D{{"$and", bson.A{permissionFilter, bson.D{{"user_id", bson.D{{"$in", friendIDList}}}}}}},
			},
			},
		}
	} else {
		filter = bson.D{
			{"$or", bson.A{
				bson.D{{"user_id", userID}}, //self
				permissionFilter,
			},
			},
		}
	}
	result, err := db.WorkMomentCollection.Find(ctx, filter, findOpts)
	if err != nil {
		return workMomentList, err
	}
	err = result.All(ctx, &workMomentList)
	return workMomentList, err
}

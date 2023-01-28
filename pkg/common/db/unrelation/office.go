package unrelation

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const cTag = "tag"
const cSendLog = "send_log"
const cWorkMoment = "work_moment"

type OfficeMgoDB struct {
	mgoDB                *mongo.Database
	TagCollection        *mongo.Collection
	TagSendLogCollection *mongo.Collection
	WorkMomentCollection *mongo.Collection
}

type Tag struct {
	UserID   string   `bson:"user_id"`
	TagID    string   `bson:"tag_id"`
	TagName  string   `bson:"tag_name"`
	UserList []string `bson:"user_list"`
}

type commonUser struct {
	UserID   string `bson:"user_id"`
	UserName string `bson:"user_name"`
}

type TagSendLog struct {
	UserList         []commonUser `bson:"tag_list"`
	SendID           string       `bson:"send_id"`
	SenderPlatformID int32        `bson:"sender_platform_id"`
	Content          string       `bson:"content"`
	SendTime         int64        `bson:"send_time"`
}

type WorkMoment struct {
	WorkMomentID         string        `bson:"work_moment_id"`
	UserID               string        `bson:"user_id"`
	UserName             string        `bson:"user_name"`
	FaceURL              string        `bson:"face_url"`
	Content              string        `bson:"content"`
	LikeUserList         []*commonUser `bson:"like_user_list"`
	AtUserList           []*commonUser `bson:"at_user_list"`
	PermissionUserList   []*commonUser `bson:"permission_user_list"`
	Comments             []*commonUser `bson:"comments"`
	PermissionUserIDList []string      `bson:"permission_user_id_list"`
	Permission           int32         `bson:"permission"`
	CreateTime           int32         `bson:"create_time"`
}

type Comment struct {
	UserID        string `bson:"user_id" json:"user_id"`
	UserName      string `bson:"user_name" json:"user_name"`
	ReplyUserID   string `bson:"reply_user_id" json:"reply_user_id"`
	ReplyUserName string `bson:"reply_user_name" json:"reply_user_name"`
	ContentID     string `bson:"content_id" json:"content_id"`
	Content       string `bson:"content" json:"content"`
	CreateTime    int32  `bson:"create_time" json:"create_time"`
}

func NewOfficeMgoDB(mgoDB *mongo.Database) *OfficeMgoDB {
	return &OfficeMgoDB{mgoDB: mgoDB, TagCollection: mgoDB.Collection(cTag), TagSendLogCollection: mgoDB.Collection(cSendLog), WorkMomentCollection: mgoDB.Collection(cSendLog)}
}

func (db *OfficeMgoDB) GetUserTags(ctx context.Context, userID string) ([]Tag, error) {
	var tags []Tag
	cursor, err := db.TagCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return tags, err
	}
	if err = cursor.All(ctx, &tags); err != nil {
		return tags, err
	}
	return tags, nil
}

func (db *OfficeMgoDB) CreateTag(ctx context.Context, userID, tagName string, userList []string) error {
	tagID := generateTagID(tagName, userID)
	tag := Tag{
		UserID:   userID,
		TagID:    tagID,
		TagName:  tagName,
		UserList: userList,
	}
	_, err := db.TagCollection.InsertOne(ctx, tag)
	return err
}

func (db *OfficeMgoDB) GetTagByID(ctx context.Context, userID, tagID string) (Tag, error) {
	var tag Tag
	err := db.TagCollection.FindOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}).Decode(&tag)
	return tag, err
}

func (db *OfficeMgoDB) DeleteTag(ctx context.Context, userID, tagID string) error {
	_, err := db.TagCollection.DeleteOne(ctx, bson.M{"user_id": userID, "tag_id": tagID})
	return err
}

func (db *OfficeMgoDB) SetTag(ctx context.Context, userID, tagID, newName string, increaseUserIDList []string, reduceUserIDList []string) error {
	var tag Tag
	if err := db.TagCollection.FindOne(ctx, bson.M{"tag_id": tagID, "user_id": userID}).Decode(&tag); err != nil {
		return err
	}
	if newName != "" {
		_, err := db.TagCollection.UpdateOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}, bson.M{"$set": bson.M{"tag_name": newName}})
		if err != nil {
			return err
		}
	}
	tag.UserList = append(tag.UserList, increaseUserIDList...)
	tag.UserList = utils.RemoveRepeatedStringInList(tag.UserList)
	for _, v := range reduceUserIDList {
		for i2, v2 := range tag.UserList {
			if v == v2 {
				tag.UserList[i2] = ""
			}
		}
	}
	var newUserList []string
	for _, v := range tag.UserList {
		if v != "" {
			newUserList = append(newUserList, v)
		}
	}
	_, err := db.TagCollection.UpdateOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}, bson.M{"$set": bson.M{"user_list": newUserList}})
	if err != nil {
		return err
	}
	return nil
}

func (db *OfficeMgoDB) GetUserIDListByTagID(ctx context.Context, userID, tagID string) ([]string, error) {
	var tag Tag
	err := db.TagCollection.FindOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}).Decode(&tag)
	return tag.UserList, err
}

func (db *OfficeMgoDB) SaveTagSendLog(ctx context.Context, tagSendLog *TagSendLog) error {
	_, err := db.TagSendLogCollection.InsertOne(ctx, tagSendLog)
	return err
}

func (db *OfficeMgoDB) GetTagSendLogs(ctx context.Context, userID string, showNumber, pageNumber int32) ([]TagSendLog, error) {
	var tagSendLogs []TagSendLog
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"send_time": -1})
	cursor, err := db.TagSendLogCollection.Find(ctx, bson.M{"send_id": userID}, findOpts)
	if err != nil {
		return tagSendLogs, err
	}
	err = cursor.All(ctx, &tagSendLogs)
	return tagSendLogs, err
}

func (db *OfficeMgoDB) CreateOneWorkMoment(ctx context.Context, workMoment *WorkMoment) error {
	workMomentID := generateWorkMomentID(workMoment.UserID)
	workMoment.WorkMomentID = workMomentID
	workMoment.CreateTime = int32(time.Now().Unix())
	_, err := db.WorkMomentCollection.InsertOne(ctx, workMoment)
	return err
}

func (db *OfficeMgoDB) DeleteOneWorkMoment(ctx context.Context, workMomentID string) error {
	_, err := db.WorkMomentCollection.DeleteOne(ctx, bson.M{"work_moment_id": workMomentID})
	return err
}

func (db *OfficeMgoDB) DeleteComment(ctx context.Context, workMomentID, contentID, opUserID string) error {
	_, err := db.WorkMomentCollection.UpdateOne(ctx, bson.D{{"work_moment_id", workMomentID},
		{"$or", bson.A{
			bson.D{{"user_id", opUserID}},
			bson.D{{"comments", bson.M{"$elemMatch": bson.M{"user_id": opUserID}}}},
		},
		}}, bson.M{"$pull": bson.M{"comments": bson.M{"content_id": contentID}}})
	return err
}

func (db *OfficeMgoDB) GetWorkMomentByID(ctx context.Context, workMomentID string) (*WorkMoment, error) {
	workMoment := &WorkMoment{}
	err := db.WorkMomentCollection.FindOne(ctx, bson.M{"work_moment_id": workMomentID}).Decode(workMoment)
	return workMoment, err
}

func (db *OfficeMgoDB) LikeOneWorkMoment(ctx context.Context, likeUserID, userName, workMomentID string) (*WorkMoment, bool, error) {
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
		workMoment.LikeUserList = append(workMoment.LikeUserList, &commonUser{UserID: likeUserID, UserName: userName})
	}
	_, err = db.WorkMomentCollection.UpdateOne(ctx, bson.M{"work_moment_id": workMomentID}, bson.M{"$set": bson.M{"like_user_list": workMoment.LikeUserList}})
	return workMoment, !isAlreadyLike, err
}

func (db *OfficeMgoDB) SetUserWorkMomentsLevel(ctx context.Context, userID string, level int32) error {
	return nil
}

func (db *OfficeMgoDB) CommentOneWorkMoment(ctx context.Context, comment *Comment, workMomentID string) (WorkMoment, error) {
	comment.ContentID = generateWorkMomentCommentID(workMomentID)
	var workMoment WorkMoment
	err := db.WorkMomentCollection.FindOneAndUpdate(ctx, bson.M{"work_moment_id": workMomentID}, bson.M{"$push": bson.M{"comments": comment}}).Decode(&workMoment)
	return workMoment, err
}

func (db *OfficeMgoDB) GetUserSelfWorkMoments(ctx context.Context, userID string, showNumber, pageNumber int32) ([]WorkMoment, error) {
	var workMomentList []WorkMoment
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"create_time": -1})
	result, err := db.WorkMomentCollection.Find(ctx, bson.M{"user_id": userID}, findOpts)
	if err != nil {
		return workMomentList, nil
	}
	err = result.All(ctx, &workMomentList)
	return workMomentList, err
}

func (db *OfficeMgoDB) GetUserWorkMoments(ctx context.Context, opUserID, userID string, showNumber, pageNumber int32, friendIDList []string) ([]WorkMoment, error) {
	var workMomentList []WorkMoment
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

func (db *OfficeMgoDB) GetUserFriendWorkMoments(ctx context.Context, showNumber, pageNumber int32, userID string, friendIDList []string) ([]WorkMoment, error) {
	var workMomentList []WorkMoment
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

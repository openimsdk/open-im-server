package unrelation

import (
	"Open_IM/pkg/common/db/table/unrelation"
	"Open_IM/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TagMongoDriver struct {
	mgoDB                *mongo.Database
	TagCollection        *mongo.Collection
	TagSendLogCollection *mongo.Collection
}

func NewTagMongoDriver(mgoDB *mongo.Database) *TagMongoDriver {
	return &TagMongoDriver{mgoDB: mgoDB, TagCollection: mgoDB.Collection(unrelation.CTag), TagSendLogCollection: mgoDB.Collection(unrelation.CSendLog)}
}

func (db *TagMongoDriver) GetUserTags(ctx context.Context, userID string) ([]unrelation.TagModel, error) {
	var tags []unrelation.TagModel
	cursor, err := db.TagCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return tags, err
	}
	if err = cursor.All(ctx, &tags); err != nil {
		return tags, err
	}
	return tags, nil
}

func (db *TagMongoDriver) CreateTag(ctx context.Context, userID, tagName string, userList []string) error {
	tagID := generateTagID(tagName, userID)
	tag := unrelation.TagModel{
		UserID:   userID,
		TagID:    tagID,
		TagName:  tagName,
		UserList: userList,
	}
	_, err := db.TagCollection.InsertOne(ctx, tag)
	return err
}

func (db *TagMongoDriver) GetTagByID(ctx context.Context, userID, tagID string) (unrelation.TagModel, error) {
	var tag unrelation.TagModel
	err := db.TagCollection.FindOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}).Decode(&tag)
	return tag, err
}

func (db *TagMongoDriver) DeleteTag(ctx context.Context, userID, tagID string) error {
	_, err := db.TagCollection.DeleteOne(ctx, bson.M{"user_id": userID, "tag_id": tagID})
	return err
}

func (db *TagMongoDriver) SetTag(ctx context.Context, userID, tagID, newName string, increaseUserIDList []string, reduceUserIDList []string) error {
	var tag unrelation.TagModel
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

func (db *TagMongoDriver) GetUserIDListByTagID(ctx context.Context, userID, tagID string) ([]string, error) {
	var tag unrelation.TagModel
	err := db.TagCollection.FindOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}).Decode(&tag)
	return tag.UserList, err
}

func (db *TagMongoDriver) SaveTagSendLog(ctx context.Context, tagSendLog *unrelation.TagSendLogModel) error {
	_, err := db.TagSendLogCollection.InsertOne(ctx, tagSendLog)
	return err
}

func (db *TagMongoDriver) GetTagSendLogs(ctx context.Context, userID string, showNumber, pageNumber int32) ([]unrelation.TagSendLogModel, error) {
	var tagSendLogs []unrelation.TagSendLogModel
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"send_time": -1})
	cursor, err := db.TagSendLogCollection.Find(ctx, bson.M{"send_id": userID}, findOpts)
	if err != nil {
		return tagSendLogs, err
	}
	err = cursor.All(ctx, &tagSendLogs)
	return tagSendLogs, err
}

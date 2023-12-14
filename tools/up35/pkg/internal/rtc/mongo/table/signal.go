package table

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/pagination"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignalModel struct {
	SID           string    `bson:"sid"`
	InviterUserID string    `bson:"inviter_user_id"`
	CustomData    string    `bson:"custom_data"`
	GroupID       string    `bson:"group_id"`
	RoomID        string    `bson:"room_id"`
	Timeout       int32     `bson:"timeout"`
	MediaType     string    `bson:"media_type"`
	PlatformID    int32     `bson:"platform_id"`
	SessionType   int32     `bson:"session_type"`
	InitiateTime  time.Time `bson:"initiate_time"`
	EndTime       time.Time `bson:"end_time"`
	FileURL       string    `bson:"file_url"`

	Title         string `bson:"title"`
	Desc          string `bson:"desc"`
	Ex            string `bson:"ex"`
	IOSPushSound  string `bson:"ios_push_sound"`
	IOSBadgeCount bool   `bson:"ios_badge_count"`
	SignalInfo    string `bson:"signal_info"`
}

type SignalInterface interface {
	Find(ctx context.Context, sids []string) ([]*SignalModel, error)
	CreateSignal(ctx context.Context, signalModel *SignalModel) error
	Update(ctx context.Context, sid string, update map[string]any) error
	UpdateSignalFileURL(ctx context.Context, sID, fileURL string) error
	UpdateSignalEndTime(ctx context.Context, sID string, endTime time.Time) error
	Delete(ctx context.Context, sids []string) error
	PageSignal(ctx context.Context, sesstionType int32, sendID string, startTime, endTime time.Time, pagination pagination.Pagination) (int64, []*SignalModel, error)
}

type SignalInvitationModel struct {
	SID          string    `bson:"sid"`
	UserID       string    `bson:"user_id"`
	Status       int32     `bson:"status"`
	InitiateTime time.Time `bson:"initiate_time"`
	HandleTime   time.Time `bson:"handle_time"`
}

type SignalInvitationInterface interface {
	Find(ctx context.Context, sid string) ([]*SignalInvitationModel, error)
	CreateSignalInvitation(ctx context.Context, sid string, inviteeUserIDs []string) error
	HandleSignalInvitation(ctx context.Context, sID, InviteeUserID string, status int32) error
	PageSID(ctx context.Context, recvID string, startTime, endTime time.Time, pagination pagination.Pagination) (int64, []string, error)
	Delete(ctx context.Context, sids []string) error
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	err = errs.Unwrap(err)
	return err == mongo.ErrNoDocuments || err == redis.Nil
}

func IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	return mongo.IsDuplicateKeyError(errs.Unwrap(err))
}

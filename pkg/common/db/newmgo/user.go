package newmgo

//import (
//	"context"
//	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo/mgotool"
//	"time"
//)
//
//type UserModel struct {
//	UserID           string    `bson:"user_id"`
//	Nickname         string    `bson:"nickname"`
//	FaceURL          string    `bson:"face_url"`
//	Ex               string    `bson:"ex"`
//	AppMangerLevel   int32     `bson:"app_manger_level"`
//	GlobalRecvMsgOpt int32     `bson:"global_recv_msg_opt"`
//	CreateTime       time.Time `bson:"create_time"`
//}
//
//type UserModelInterface interface {
//	Create(ctx context.Context, users []*UserModel) (err error)
//	UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error)
//	Find(ctx context.Context, userIDs []string) (users []*UserModel, err error)
//	Take(ctx context.Context, userID string) (user *UserModel, err error)
//	Page(ctx context.Context, pagination mgotool.Pagination) (count int64, users []*UserModel, err error)
//	Exist(ctx context.Context, userID string) (exist bool, err error)
//	GetAllUserID(ctx context.Context, pagination mgotool.Pagination) (count int64, userIDs []string, err error)
//	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
//	// 获取用户总数
//	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
//	// 获取范围内用户增量
//	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
//}

package relation

import (
	"context"
	"time"
)

const (
	GroupRequestModelTableName = "group_requests"
)

type GroupRequestModel struct {
	UserID        string    `gorm:"column:user_id;primary_key;size:64"`
	GroupID       string    `gorm:"column:group_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:1024"`
	HandledMsg    string    `gorm:"column:handle_msg;size:1024"`
	ReqTime       time.Time `gorm:"column:req_time"`
	HandleUserID  string    `gorm:"column:handle_user_id;size:64"`
	HandledTime   time.Time `gorm:"column:handle_time"`
	JoinSource    int32     `gorm:"column:join_source"`
	InviterUserID string    `gorm:"column:inviter_user_id;size:64"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

func (GroupRequestModel) TableName() string {
	return GroupRequestModelTableName
}

type GroupRequestModelInterface interface {
	Create(ctx context.Context, groupRequests []*GroupRequestModel, tx ...any) (err error)
	//Delete(ctx context.Context, groupRequests []*GroupRequestModel, tx ...any) (err error)
	//UpdateMap(ctx context.Context, groupID string, userID string, args map[string]interface{}, tx ...any) (err error)
	UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, tx ...any) (err error)
	//Update(ctx context.Context, groupRequests []*GroupRequestModel, tx ...any) (err error)
	//Find(ctx context.Context, groupRequests []*GroupRequestModel, tx ...any) (resultGroupRequests []*GroupRequestModel, err error)
	Take(ctx context.Context, groupID string, userID string, tx ...any) (groupRequest *GroupRequestModel, err error)
	Page(ctx context.Context, userID string, pageNumber, showNumber int32, tx ...any) (total uint32, groups []*GroupRequestModel, err error)
}

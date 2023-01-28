package relation

import (
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

var GroupRequestDB *gorm.DB

type GroupRequest struct {
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

func (GroupRequest) TableName() string {
	return "friend_requests"
}

func (*GroupRequest) Create(ctx context.Context, groupRequests []*GroupRequest) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
	}()
	return utils.Wrap(GroupRequestDB.Create(&groupRequests).Error, utils.GetSelfFuncName())
}

func (*GroupRequest) Delete(ctx context.Context, groupRequests []*GroupRequest) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
	}()
	return utils.Wrap(GroupRequestDB.Delete(&groupRequests).Error, utils.GetSelfFuncName())
}

func (*GroupRequest) UpdateByMap(ctx context.Context, groupID string, userID string, args map[string]interface{}) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "args", args)
	}()
	return utils.Wrap(GroupRequestDB.Where("group_id = ? and user_id = ? ", groupID, userID).Updates(args).Error, utils.GetSelfFuncName())
}

func (*GroupRequest) Update(ctx context.Context, groupRequests []*GroupRequest) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
	}()
	return utils.Wrap(GroupRequestDB.Updates(&groupRequests).Error, utils.GetSelfFuncName())
}

func (*GroupRequest) Find(ctx context.Context, groupRequests []*GroupRequest) (resultGroupRequests []*GroupRequest, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests, "resultGroupRequests", resultGroupRequests)
	}()
	var where [][]interface{}
	for _, groupMember := range groupRequests {
		where = append(where, []interface{}{groupMember.GroupID, groupMember.UserID})
	}
	return resultGroupRequests, utils.Wrap(GroupRequestDB.Where("(group_id, user_id) in ?", where).Find(&resultGroupRequests).Error, utils.GetSelfFuncName())
}

func (*GroupRequest) Take(ctx context.Context, groupID string, userID string) (groupRequest *GroupRequest, err error) {
	groupRequest = &GroupRequest{}
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "groupRequest", *groupRequest)
	}()
	return groupRequest, utils.Wrap(GroupRequestDB.Where("group_id = ? and user_id = ? ", groupID, userID).Take(groupRequest).Error, utils.GetSelfFuncName())
}

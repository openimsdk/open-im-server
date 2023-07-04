package msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"time"
)

func (m *msgServer) GetActiveUser(ctx context.Context, req *msg.GetActiveUserReq) (*msg.GetActiveUserResp, error) {
	msgCount, userCount, users, dateCount, err := m.MsgDatabase.RangeUserSendCount(ctx, time.UnixMilli(req.Start), time.UnixMilli(req.End), req.Ase, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	var pbUsers []*msg.ActiveUser
	if len(users) > 0 {
		userIDs := utils.Slice(users, func(e *unrelation.UserCount) string { return e.UserID })
		userMap, err := m.User.GetUsersInfoMap(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		pbUsers = make([]*msg.ActiveUser, 0, len(users))
		for _, user := range users {
			pbUser := userMap[user.UserID]
			if pbUser == nil {
				pbUser = &sdkws.UserInfo{
					UserID:   user.UserID,
					Nickname: user.UserID,
				}
			}
			pbUsers = append(pbUsers, &msg.ActiveUser{
				User:  pbUser,
				Count: user.Count,
			})
		}
	}
	return &msg.GetActiveUserResp{
		MsgCount:  msgCount,
		UserCount: userCount,
		DateCount: dateCount,
		Users:     pbUsers,
	}, nil
}

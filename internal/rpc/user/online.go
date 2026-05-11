package user

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/protocol/constant"
	pbuser "github.com/openimsdk/protocol/user"
)

func (s *userServer) getUserOnlineStatus(ctx context.Context, userID string) (*pbuser.OnlineStatus, error) {
	platformIDs, err := s.online.GetOnline(ctx, userID)
	if err != nil {
		return nil, err
	}
	status := pbuser.OnlineStatus{
		UserID:      userID,
		PlatformIDs: platformIDs,
	}
	if len(platformIDs) > 0 {
		status.Status = constant.Online
	} else {
		status.Status = constant.Offline
	}
	return &status, nil
}

func (s *userServer) getUsersOnlineStatus(ctx context.Context, userIDs []string) ([]*pbuser.OnlineStatus, error) {
	res := make([]*pbuser.OnlineStatus, 0, len(userIDs))
	for _, userID := range userIDs {
		status, err := s.getUserOnlineStatus(ctx, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, status)
	}
	return res, nil
}

// SubscribeOrCancelUsersStatus Subscribe online or cancel online users.
func (s *userServer) SubscribeOrCancelUsersStatus(ctx context.Context, req *pbuser.SubscribeOrCancelUsersStatusReq) (*pbuser.SubscribeOrCancelUsersStatusResp, error) {
	return &pbuser.SubscribeOrCancelUsersStatusResp{}, nil
}

// GetUserStatus Get the online status of the user.
func (s *userServer) GetUserStatus(ctx context.Context, req *pbuser.GetUserStatusReq) (*pbuser.GetUserStatusResp, error) {
	res, err := s.getUsersOnlineStatus(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetUserStatusResp{StatusList: res}, nil
}

// SetUserStatus Synchronize user's online status.
func (s *userServer) SetUserStatus(ctx context.Context, req *pbuser.SetUserStatusReq) (*pbuser.SetUserStatusResp, error) {
	var (
		online  []int32
		offline []int32
	)
	switch req.Status {
	case constant.Online:
		online = []int32{req.PlatformID}
	case constant.Offline:
		offline = []int32{req.PlatformID}
	}
	if err := s.online.SetUserOnline(ctx, req.UserID, online, offline); err != nil {
		return nil, err
	}
	return &pbuser.SetUserStatusResp{}, nil
}

// GetSubscribeUsersStatus Get the online status of subscribers.
func (s *userServer) GetSubscribeUsersStatus(ctx context.Context, req *pbuser.GetSubscribeUsersStatusReq) (*pbuser.GetSubscribeUsersStatusResp, error) {
	return &pbuser.GetSubscribeUsersStatusResp{}, nil
}

func (s *userServer) SetUserOnlineStatus(ctx context.Context, req *pbuser.SetUserOnlineStatusReq) (*pbuser.SetUserOnlineStatusResp, error) {
	for _, status := range req.Status {
		if err := s.online.SetUserOnline(ctx, status.UserID, status.Online, status.Offline); err != nil {
			return nil, err
		}
		s.updateOfflineRecord(ctx, status.UserID, len(status.Offline) > 0, len(status.Online) > 0)
	}
	return &pbuser.SetUserOnlineStatusResp{}, nil
}

// updateOfflineRecord 根据用户当前在线状态维护 user_offline_record 集合：
//   - 若某平台刚断开且用户已全平台离线 → upsert 离线记录（仅首次写入，保留最早离线时刻）
//   - 若某平台刚上线且用户当前有在线平台 → 删除离线记录（停止计时）
func (s *userServer) updateOfflineRecord(ctx context.Context, userID string, hasOffline, hasOnline bool) {
	if !hasOffline && !hasOnline {
		return
	}
	platformIDs, err := s.online.GetOnline(ctx, userID)
	if err != nil {
		log.ZWarn(ctx, "updateOfflineRecord: GetOnline failed", err, "userID", userID)
		return
	}
	if len(platformIDs) == 0 {
		// 所有平台已离线，写入离线记录（含预计算的删除截止时间）
		offlineTime := time.Now()
		deadline := s.calcDeleteDeadline(ctx, userID, offlineTime)
		if err := s.userOfflineRecord.Upsert(ctx, userID, offlineTime, deadline); err != nil {
			log.ZWarn(ctx, "updateOfflineRecord: Upsert failed", err, "userID", userID)
		}
	} else if hasOnline {
		// 用户重新上线，删除离线记录，停止计时
		if err := s.userOfflineRecord.Delete(ctx, userID); err != nil {
			log.ZWarn(ctx, "updateOfflineRecord: Delete failed", err, "userID", userID)
		}
	}
}

// calcDeleteDeadline 查询用户的 delete_account_interval 并计算删除截止时间。
// 若查询失败或 interval 为 0，则使用系统默认值（18 个月）。
func (s *userServer) calcDeleteDeadline(ctx context.Context, userID string, from time.Time) time.Time {
	interval := int32(model.DefaultDeleteAccountIntervalSec)
	users, err := s.db.Find(ctx, []string{userID})
	if err == nil && len(users) > 0 && users[0].DeleteAccountInterval > 0 {
		interval = users[0].DeleteAccountInterval
	}
	return from.Add(time.Duration(interval) * time.Second)
}

func (s *userServer) GetAllOnlineUsers(ctx context.Context, req *pbuser.GetAllOnlineUsersReq) (*pbuser.GetAllOnlineUsersResp, error) {
	resMap, nextCursor, err := s.online.GetAllOnlineUsers(ctx, req.Cursor)
	if err != nil {
		return nil, err
	}
	resp := &pbuser.GetAllOnlineUsersResp{
		StatusList: make([]*pbuser.OnlineStatus, 0, len(resMap)),
		NextCursor: nextCursor,
	}
	for userID, plats := range resMap {
		resp.StatusList = append(resp.StatusList, &pbuser.OnlineStatus{
			UserID:      userID,
			Status:      int32(datautil.If(len(plats) > 0, constant.Online, constant.Offline)),
			PlatformIDs: plats,
		})
	}
	return resp, nil
}

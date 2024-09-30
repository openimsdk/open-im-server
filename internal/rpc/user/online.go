package user

import (
	"context"

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
	}
	return &pbuser.SetUserOnlineStatusResp{}, nil
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

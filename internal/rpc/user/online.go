package user

import (
	"context"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
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
	if req.Genre == constant.SubscriberUser {
		err := s.db.SubscribeUsersStatus(ctx, req.UserID, req.UserIDs)
		if err != nil {
			return nil, err
		}
		var status []*pbuser.OnlineStatus
		status, err = s.getUsersOnlineStatus(ctx, req.UserIDs)
		if err != nil {
			return nil, err
		}
		return &pbuser.SubscribeOrCancelUsersStatusResp{StatusList: status}, nil
	} else if req.Genre == constant.Unsubscribe {
		err := s.db.UnsubscribeUsersStatus(ctx, req.UserID, req.UserIDs)
		if err != nil {
			return nil, err
		}
	}
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
		online = []int32{req.PlatformID}
	}
	if err := s.online.SetUserOnline(ctx, req.UserID, online, offline); err != nil {
		return nil, err
	}
	list, err := s.db.GetSubscribedList(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, userID := range list {
		tips := &sdkws.UserStatusChangeTips{
			FromUserID: req.UserID,
			ToUserID:   userID,
			Status:     req.Status,
			PlatformID: req.PlatformID,
		}
		s.userNotificationSender.UserStatusChangeNotification(ctx, tips)
	}

	return &pbuser.SetUserStatusResp{}, nil
}

// GetSubscribeUsersStatus Get the online status of subscribers.
func (s *userServer) GetSubscribeUsersStatus(ctx context.Context, req *pbuser.GetSubscribeUsersStatusReq) (*pbuser.GetSubscribeUsersStatusResp, error) {
	userList, err := s.db.GetAllSubscribeList(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	onlineStatusList, err := s.getUsersOnlineStatus(ctx, userList)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetSubscribeUsersStatusResp{StatusList: onlineStatusList}, nil
}

func (s *userServer) SetUserOnlineStatus(ctx context.Context, req *pbuser.SetUserOnlineStatusReq) (*pbuser.SetUserOnlineStatusResp, error) {
	for _, status := range req.Status {
		if err := s.online.SetUserOnline(ctx, status.UserID, status.Online, status.Offline); err != nil {
			return nil, err
		}
	}
	return &pbuser.SetUserOnlineStatusResp{}, nil
}

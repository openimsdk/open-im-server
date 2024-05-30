package friend

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbfriend "github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/tools/errs"
)

func (s *friendServer) NotificationUserInfoUpdate(ctx context.Context, req *pbfriend.NotificationUserInfoUpdateReq) (*pbfriend.NotificationUserInfoUpdateResp, error) {
	if req.NewUserInfo == nil {
		var err error
		req.NewUserInfo, err = s.userRpcClient.GetUserInfo(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
	}
	if req.UserID != req.NewUserInfo.UserID {
		return nil, errs.ErrArgs.WrapMsg("req.UserID != req.NewUserInfo.UserID")
	}
	userIDs, err := s.friendDatabase.FindFriendUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if len(userIDs) > 0 {
		if err := s.friendDatabase.UpdateFriendUserInfo(ctx, req.UserID, userIDs, req.NewUserInfo.Nickname, req.NewUserInfo.FaceURL); err != nil {
			return nil, err
		}
		s.notificationSender.FriendsInfoUpdateNotification(ctx, req.UserID, userIDs)
	}
	return &pbfriend.NotificationUserInfoUpdateResp{}, nil
}

func (s *friendServer) SearchFriends(ctx context.Context, req *pbfriend.SearchFriendsReq) (*pbfriend.SearchFriendsResp, error) {
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	if req.Keyword == "" {
		total, friends, err := s.friendDatabase.PageOwnerFriends(ctx, req.UserID, req.Pagination)
		if err != nil {
			return nil, err
		}
		return &pbfriend.SearchFriendsResp{
			Total:   total,
			Friends: friendsDB2PB(friends),
		}, nil
	}
	total, friends, err := s.friendDatabase.SearchFriend(ctx, req.UserID, req.Keyword, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &pbfriend.SearchFriendsResp{
		Total:   total,
		Friends: friendsDB2PB(friends),
	}, nil
}

func (s *friendServer) sortFriendUserIDsHash(userIDs []string) uint64 {
	data, _ := json.Marshal(userIDs)
	sum := md5.Sum(data)
	return binary.BigEndian.Uint64(sum[:])
}

func (s *friendServer) GetIncrementalFriends(ctx context.Context, req *pbfriend.GetIncrementalFriendsReq) (*pbfriend.GetIncrementalFriendsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	var limit int
	if req.Version > 0 {
		limit = s.config.RpcConfig.FriendSyncCount
	}
	incrVer, err := s.friendDatabase.FindFriendIncrVersion(ctx, req.UserID, uint(req.Version), limit)
	if err != nil {
		return nil, err
	}
	var (
		deleteUserIDs []string
		changeUserIDs []string
	)
	if incrVer.Full() {
		changeUserIDs, err = s.friendDatabase.FindSortFriendUserIDs(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		deleteUserIDs, changeUserIDs = incrVer.DeleteAndChangeIDs()
	}
	var friends []*model.Friend
	if len(changeUserIDs) > 0 {
		friends, err = s.friendDatabase.FindFriendsWithError(ctx, req.UserID, changeUserIDs)
		if err != nil {
			return nil, err
		}
	}
	return &pbfriend.GetIncrementalFriendsResp{
		Version:       uint64(incrVer.Version),
		VersionID:     incrVer.ID.Hex(),
		Full:          incrVer.Full(),
		SyncCount:     uint32(s.config.RpcConfig.FriendSyncCount),
		DeleteUserIds: deleteUserIDs,
		Changes:       friendsDB2PB(friends),
	}, nil
}

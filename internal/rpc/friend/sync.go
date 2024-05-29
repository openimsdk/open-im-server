package friend

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	pbfriend "github.com/openimsdk/protocol/friend"
)

func (s *friendServer) SearchFriends(ctx context.Context, req *pbfriend.SearchFriendsReq) (*pbfriend.SearchFriendsResp, error) {
	//TODO implement me
	panic("implement me")
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
	var friends []*relation.FriendModel
	if len(changeUserIDs) > 0 {
		friends, err = s.friendDatabase.FindFriendsWithError(ctx, req.UserID, changeUserIDs)
		if err != nil {
			return nil, err
		}
	}
	return &pbfriend.GetIncrementalFriendsResp{
		Version:       uint64(incrVer.Version),
		VersionID:     incrVer.ID.String(),
		Full:          incrVer.Full(),
		SyncCount:     uint32(s.config.RpcConfig.FriendSyncCount),
		DeleteUserIds: deleteUserIDs,
		Changes:       friendsDB2PB(friends),
	}, nil
}

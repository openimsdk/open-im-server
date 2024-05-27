package friend

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/dataver"
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
	sortUserIDs, err := s.friendDatabase.FindSortFriendUserIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if len(sortUserIDs) == 0 {
		return &pbfriend.GetIncrementalFriendsResp{
			Version:   uint64(incrVer.Version),
			VersionID: dataver.VersionIDStr(incrVer.ID),
			Full:      true,
			SyncCount: uint32(s.config.RpcConfig.FriendSyncCount),
		}, nil
	}
	var changes []*relation.FriendModel
	res := dataver.NewSyncResult(incrVer, sortUserIDs, req.VersionID)
	if len(res.Changes) > 0 {
		changes, err = s.friendDatabase.FindFriendsWithError(ctx, req.UserID, res.Changes)
		if err != nil {
			return nil, err
		}
	}
	calcHash := s.sortFriendUserIDsHash(sortUserIDs)
	if calcHash == req.IdHash {
		sortUserIDs = nil
	}
	return &pbfriend.GetIncrementalFriendsResp{
		Version:        uint64(res.Version),
		VersionID:      res.VersionID,
		Full:           res.Full,
		SyncCount:      uint32(s.config.RpcConfig.FriendSyncCount),
		SortUserIdHash: calcHash,
		SortUserIds:    sortUserIDs,
		DeleteUserIds:  res.DeleteEID,
		Changes:        friendsDB2PB(changes),
	}, nil
}

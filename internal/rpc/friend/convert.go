package friend

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/tools/utils/datautil"
)

func friendDB2PB(db *model.Friend) *friend.FriendInfo {
	return &friend.FriendInfo{
		OwnerUserID:    db.OwnerUserID,
		FriendUserID:   db.FriendUserID,
		FriendNickname: db.FriendNickname,
		FriendFaceURL:  db.FriendFaceURL,
		Remark:         db.Remark,
		CreateTime:     db.CreateTime.UnixMilli(),
		AddSource:      db.AddSource,
		OperatorUserID: db.OperatorUserID,
		Ex:             db.Ex,
		IsPinned:       db.IsPinned,
	}
}

func friendsDB2PB(db []*model.Friend) []*friend.FriendInfo {
	return datautil.Slice(db, friendDB2PB)
}

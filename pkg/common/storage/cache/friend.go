package cache

import (
	"context"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// FriendCache is an interface for caching friend-related data.
type FriendCache interface {
	BatchDeleter
	CloneFriendCache() FriendCache
	GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error)
	// Called when friendID list changed
	DelFriendIDs(ownerUserID ...string) FriendCache
	// Get single friendInfo from the cache
	GetFriend(ctx context.Context, ownerUserID, friendUserID string) (friend *relationtb.Friend, err error)
	// Delete friend when friend info changed
	DelFriend(ownerUserID, friendUserID string) FriendCache
	// Delete friends when friends' info changed
	DelFriends(ownerUserID string, friendUserIDs []string) FriendCache

	DelOwner(friendUserID string, ownerUserIDs []string) FriendCache

	DelMaxFriendVersion(ownerUserIDs ...string) FriendCache

	//DelSortFriendUserIDs(ownerUserIDs ...string) FriendCache

	//FindSortFriendUserIDs(ctx context.Context, ownerUserID string) ([]string, error)

	//FindFriendIncrVersion(ctx context.Context, ownerUserID string, version uint, limit int) (*relationtb.VersionLog, error)

	FindMaxFriendVersion(ctx context.Context, ownerUserID string) (*relationtb.VersionLog, error)
}

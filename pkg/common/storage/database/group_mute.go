package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// GroupMute persists per-user group notification mute settings.
type GroupMute interface {
	Upsert(ctx context.Context, mute *model.GroupMute) error
	Delete(ctx context.Context, ownerUserID, groupID string) error
	// ListActiveMutedUserIDs returns which of candidateUserIDs currently have an active mute on this group.
	ListActiveMutedUserIDs(ctx context.Context, groupID string, candidateUserIDs []string) ([]string, error)
	// Get returns one document by owner + group; nil if not found.
	Get(ctx context.Context, ownerUserID, groupID string) (*model.GroupMute, error)
}

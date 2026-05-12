package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// GroupMuteDatabase per-user group notification mute.
type GroupMuteDatabase interface {
	Upsert(ctx context.Context, mute *model.GroupMute) error
	Delete(ctx context.Context, ownerUserID, groupID string) error
	ListActiveMutedUserIDs(ctx context.Context, groupID string, candidateUserIDs []string) ([]string, error)
	Get(ctx context.Context, ownerUserID, groupID string) (*model.GroupMute, error)
}

type groupMuteDatabase struct {
	db database.GroupMute
}

func NewGroupMuteDatabase(db database.GroupMute) GroupMuteDatabase {
	return &groupMuteDatabase{db: db}
}

func (g *groupMuteDatabase) Upsert(ctx context.Context, mute *model.GroupMute) error {
	return g.db.Upsert(ctx, mute)
}

func (g *groupMuteDatabase) Delete(ctx context.Context, ownerUserID, groupID string) error {
	return g.db.Delete(ctx, ownerUserID, groupID)
}

func (g *groupMuteDatabase) ListActiveMutedUserIDs(ctx context.Context, groupID string, candidateUserIDs []string) ([]string, error) {
	return g.db.ListActiveMutedUserIDs(ctx, groupID, candidateUserIDs)
}

func (g *groupMuteDatabase) Get(ctx context.Context, ownerUserID, groupID string) (*model.GroupMute, error) {
	return g.db.Get(ctx, ownerUserID, groupID)
}

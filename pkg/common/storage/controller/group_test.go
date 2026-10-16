package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createGroupTestTx func(context.Context, func(context.Context) error) error

func (fn createGroupTestTx) Transaction(ctx context.Context, callback func(context.Context) error) error {
	return fn(ctx, callback)
}

type createGroupTestGroupDB struct {
	database.Group
	create func(context.Context, []*model.Group) error
}

func (db createGroupTestGroupDB) Create(ctx context.Context, groups []*model.Group) error {
	return db.create(ctx, groups)
}

type createGroupTestMemberDB struct {
	database.GroupMember
	create func(context.Context, []*model.GroupMember) error
}

func (db createGroupTestMemberDB) Create(ctx context.Context, members []*model.GroupMember) error {
	return db.create(ctx, members)
}

type createGroupTestCache struct {
	cache.GroupCache
	clone      func() cache.GroupCache
	invalidate func(context.Context) error
}

func (c *createGroupTestCache) CloneGroupCache() cache.GroupCache {
	if c.clone != nil {
		return c.clone()
	}
	return c
}

func (c *createGroupTestCache) ChainExecDel(ctx context.Context) error {
	return c.invalidate(ctx)
}

func (c *createGroupTestCache) DelGroupsInfo(...string) cache.GroupCache               { return c }
func (c *createGroupTestCache) DelGroupMembersHash(string) cache.GroupCache            { return c }
func (c *createGroupTestCache) DelGroupsMemberNum(...string) cache.GroupCache          { return c }
func (c *createGroupTestCache) DelGroupMemberIDs(string) cache.GroupCache              { return c }
func (c *createGroupTestCache) DelGroupAllRoleLevel(string) cache.GroupCache           { return c }
func (c *createGroupTestCache) DelJoinedGroupID(...string) cache.GroupCache            { return c }
func (c *createGroupTestCache) DelGroupMembersInfo(string, ...string) cache.GroupCache { return c }
func (c *createGroupTestCache) DelMaxJoinGroupVersion(...string) cache.GroupCache      { return c }
func (c *createGroupTestCache) DelMaxGroupMemberVersion(...string) cache.GroupCache    { return c }

func TestCreateGroupInvalidatesCacheAfterCommit(t *testing.T) {
	for _, tt := range []struct {
		name    string
		groups  []*model.Group
		members []*model.GroupMember
	}{
		{name: "invite members", members: []*model.GroupMember{{GroupID: "group", UserID: "new-member"}}},
		{name: "create group with members", groups: []*model.Group{{GroupID: "group"}}, members: []*model.GroupMember{{GroupID: "group", UserID: "new-member"}}},
		{name: "create group without members", groups: []*model.Group{{GroupID: "group"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			type sessionKey struct{}
			ctx := context.Background()
			sessionCtx := context.WithValue(ctx, sessionKey{}, true)
			committedIDs := []string{"existing-member"}
			var pendingIDs, cachedIDs []string
			var groupWrites, memberWrites, invalidations int
			var invalidationCtx context.Context
			committed := false
			db := &groupDatabase{
				groupDB: createGroupTestGroupDB{create: func(ctx context.Context, groups []*model.Group) error {
					require.Equal(t, true, ctx.Value(sessionKey{}))
					assert.Equal(t, tt.groups, groups)
					groupWrites++
					return nil
				}},
				groupMemberDB: createGroupTestMemberDB{create: func(ctx context.Context, members []*model.GroupMember) error {
					require.Equal(t, true, ctx.Value(sessionKey{}))
					assert.Equal(t, tt.members, members)
					memberWrites++
					for _, member := range members {
						pendingIDs = append(pendingIDs, member.UserID)
					}
					return nil
				}},
				ctxTx: createGroupTestTx(func(_ context.Context, callback func(context.Context) error) error {
					if err := callback(sessionCtx); err != nil {
						return err
					}
					committedIDs = append(committedIDs, pendingIDs...)
					committed = true
					return nil
				}),
				cache: &createGroupTestCache{invalidate: func(ctx context.Context) error {
					invalidations++
					invalidationCtx = ctx
					// A concurrent cache miss sees only the committed MongoDB snapshot.
					// Invalidating inside the transaction would repopulate the old IDs.
					cachedIDs = append([]string(nil), committedIDs...)
					assert.True(t, committed, "cache invalidation must follow the transaction commit")
					return nil
				}},
			}
			require.NoError(t, db.CreateGroup(ctx, tt.groups, tt.members))
			assert.Equal(t, committedIDs, cachedIDs, "cache must not retain the pre-commit member IDs")
			assert.Equal(t, 1, invalidations)
			assert.Nil(t, invalidationCtx.Value(sessionKey{}), "do not reuse the completed MongoDB session context")
			assert.Equal(t, min(1, len(tt.groups)), groupWrites)
			assert.Equal(t, min(1, len(tt.members)), memberWrites)
		})
	}
}

func TestCreateGroupDoesNotInvalidateOnTransactionFailure(t *testing.T) {
	writeErr := errors.New("write failed")
	commitErr := errors.New("commit failed")
	for _, tt := range []struct {
		name      string
		groupErr  error
		memberErr error
		commitErr error
		wantErr   error
	}{
		{name: "group write", groupErr: writeErr, wantErr: writeErr},
		{name: "member write", memberErr: writeErr, wantErr: writeErr},
		{name: "commit", commitErr: commitErr, wantErr: commitErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			invalidations := 0
			db := &groupDatabase{
				groupDB:       createGroupTestGroupDB{create: func(context.Context, []*model.Group) error { return tt.groupErr }},
				groupMemberDB: createGroupTestMemberDB{create: func(context.Context, []*model.GroupMember) error { return tt.memberErr }},
				ctxTx: createGroupTestTx(func(ctx context.Context, callback func(context.Context) error) error {
					if err := callback(ctx); err != nil {
						return err
					}
					return tt.commitErr
				}),
				cache: &createGroupTestCache{invalidate: func(context.Context) error {
					invalidations++
					return nil
				}},
			}
			err := db.CreateGroup(context.Background(), []*model.Group{{GroupID: "group"}}, []*model.GroupMember{{GroupID: "group", UserID: "member"}})
			require.ErrorIs(t, err, tt.wantErr)
			assert.Zero(t, invalidations)
		})
	}
}

func TestCreateGroupInvalidatesOnceAfterTransactionRetry(t *testing.T) {
	clones, invalidations, writes := 0, 0, 0
	committed := false
	newBatch := func() cache.GroupCache {
		clones++
		return &createGroupTestCache{invalidate: func(context.Context) error {
			invalidations++
			assert.True(t, committed)
			return nil
		}}
	}
	db := &groupDatabase{
		groupMemberDB: createGroupTestMemberDB{create: func(context.Context, []*model.GroupMember) error {
			writes++
			return nil
		}},
		ctxTx: createGroupTestTx(func(ctx context.Context, callback func(context.Context) error) error {
			// WithTransaction may retry the callback before a successful commit.
			for range 2 {
				if err := callback(ctx); err != nil {
					return err
				}
			}
			committed = true
			return nil
		}),
		cache: &createGroupTestCache{clone: newBatch},
	}
	require.NoError(t, db.CreateGroup(context.Background(), nil, []*model.GroupMember{{GroupID: "group", UserID: "member"}}))
	assert.Equal(t, 2, writes)
	assert.Equal(t, 2, clones)
	assert.Equal(t, 1, invalidations)
}

func TestCreateGroupReturnsPostCommitCacheError(t *testing.T) {
	cacheErr := errors.New("cache invalidation failed")
	committed := false
	db := &groupDatabase{
		groupMemberDB: createGroupTestMemberDB{create: func(context.Context, []*model.GroupMember) error { return nil }},
		ctxTx: createGroupTestTx(func(ctx context.Context, callback func(context.Context) error) error {
			if err := callback(ctx); err != nil {
				return err
			}
			committed = true
			return nil
		}),
		cache: &createGroupTestCache{invalidate: func(context.Context) error { return cacheErr }},
	}
	err := db.CreateGroup(context.Background(), nil, []*model.GroupMember{{GroupID: "group", UserID: "member"}})
	require.ErrorIs(t, err, cacheErr)
	assert.True(t, committed)
}

func TestCreateGroupEmptyInput(t *testing.T) {
	// Nil dependencies make any attempted transaction or invalidation fail.
	require.NoError(t, (&groupDatabase{}).CreateGroup(context.Background(), nil, nil))
}

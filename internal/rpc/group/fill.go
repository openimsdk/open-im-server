package group

import (
	"context"

	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func (s *groupServer) FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relationTb.GroupMemberModel, error) {
	members, err := s.GroupDatabase.FindGroupMember(ctx, groupIDs, userIDs, roleLevels)
	if err != nil {
		return nil, err
	}
	emptyUserIDs := make(map[string]struct{})
	for _, member := range members {
		if member.Nickname == "" || member.FaceURL == "" {
			emptyUserIDs[member.UserID] = struct{}{}
		}
	}
	if len(emptyUserIDs) > 0 {
		users, err := s.User.GetPublicUserInfoMap(ctx, utils.Keys(emptyUserIDs), true)
		if err != nil {
			return nil, err
		}
		for i, member := range members {
			user, ok := users[member.UserID]
			if !ok {
				continue
			}
			if member.Nickname == "" {
				members[i].Nickname = user.Nickname
			}
			if member.FaceURL == "" {
				members[i].FaceURL = user.FaceURL
			}
		}
	}
	return members, nil
}

func (s *groupServer) TakeGroupMember(ctx context.Context, groupID string, userID string) (*relationTb.GroupMemberModel, error) {
	member, err := s.GroupDatabase.TakeGroupMember(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if member.Nickname == "" || member.FaceURL == "" {
		user, err := s.User.GetPublicUserInfo(ctx, userID)
		if err != nil {
			return nil, err
		}
		if member.Nickname == "" {
			member.Nickname = user.Nickname
		}
		if member.FaceURL == "" {
			member.FaceURL = user.FaceURL
		}
	}
	return member, nil
}

func (s *groupServer) TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error) {
	owner, err := s.GroupDatabase.TakeGroupOwner(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if owner.Nickname == "" || owner.FaceURL == "" {
		user, err := s.User.GetUserInfo(ctx, owner.UserID)
		if err != nil {
			return nil, err
		}
		if owner.Nickname == "" {
			owner.Nickname = user.Nickname
		}
		if owner.FaceURL == "" {
			owner.FaceURL = user.FaceURL
		}
	}
	return owner, nil
}

func (s *groupServer) PageGetGroupMember(ctx context.Context, groupID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	total, members, err := s.GroupDatabase.PageGetGroupMember(ctx, groupID, pageNumber, showNumber)
	if err != nil {
		return 0, nil, err
	}
	emptyUserIDs := make(map[string]struct{})
	for _, member := range members {
		if member.Nickname == "" || member.FaceURL == "" {
			emptyUserIDs[member.UserID] = struct{}{}
		}
	}
	if len(emptyUserIDs) > 0 {
		users, err := s.User.GetPublicUserInfoMap(ctx, utils.Keys(emptyUserIDs), true)
		if err != nil {
			return 0, nil, err
		}
		for i, member := range members {
			user, ok := users[member.UserID]
			if !ok {
				continue
			}
			if member.Nickname == "" {
				members[i].Nickname = user.Nickname
			}
			if member.FaceURL == "" {
				members[i].FaceURL = user.FaceURL
			}
		}
	}
	return total, members, nil
}

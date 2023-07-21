// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package friend

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func (x *GetPaginationFriendsReq) Check() error {
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	return nil
}

func (x *ApplyToAddFriendReq) Check() error {
	if x.ToUserID == "" {
		return errs.ErrArgs.Wrap("toUserID is empty")
	}
	if x.FromUserID == "" {
		return errs.ErrArgs.Wrap("fromUserID is empty")
	}
	return nil
}

func (x *ImportFriendReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	if x.FriendUserIDs == nil {
		return errs.ErrArgs.Wrap("friendUserIDS is empty")
	}
	return nil
}

func (x *GetPaginationFriendsApplyToReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	return nil
}

func (x *GetDesignatedFriendsReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	if x.FriendUserIDs == nil {
		return errs.ErrArgs.Wrap("friendUserIDS is empty")
	}
	return nil
}

func (x *AddBlackReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	if x.BlackUserID == "" {
		return errs.ErrArgs.Wrap("BlackUserID is empty")
	}
	return nil
}

func (x *RemoveBlackReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	if x.BlackUserID == "" {
		return errs.ErrArgs.Wrap("BlackUserID is empty")
	}
	return nil
}

func (x *GetPaginationBlacksReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	return nil
}

func (x *IsFriendReq) Check() error {
	if x.UserID1 == "" {
		return errs.ErrArgs.Wrap("userID1 is empty")
	}
	if x.UserID2 == "" {
		return errs.ErrArgs.Wrap("userID2 is empty")
	}
	return nil
}

func (x *IsBlackReq) Check() error {
	if x.UserID1 == "" {
		return errs.ErrArgs.Wrap("userID1 is empty")
	}
	if x.UserID2 == "" {
		return errs.ErrArgs.Wrap("userID2 is empty")
	}
	return nil
}

func (x *DeleteFriendReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("userID1 is empty")
	}
	if x.FriendUserID == "" {
		return errs.ErrArgs.Wrap("userID2 is empty")
	}
	return nil
}

func (x *RespondFriendApplyReq) Check() error {
	if x.ToUserID == "" {
		return errs.ErrArgs.Wrap("toUserID is empty")
	}
	if x.FromUserID == "" {
		return errs.ErrArgs.Wrap("fromUserID is empty")
	}
	return nil
}

func (x *SetFriendRemarkReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	if x.FriendUserID == "" {
		return errs.ErrArgs.Wrap("fromUserID is empty")
	}
	return nil
}

func (x *GetPaginationFriendsApplyFromReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.Pagination == nil {
		return errs.ErrArgs.Wrap("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errs.ErrArgs.Wrap("pageNumber is invalid")
	}
	return nil
}

func (x *GetFriendIDsReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	return nil
}

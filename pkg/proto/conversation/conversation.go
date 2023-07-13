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

package conversation

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func (x *ConversationReq) Check() error {
	if x.ConversationID == "" {
		return errs.ErrArgs.Wrap("conversation is empty")
	}
	return nil
}

func (x *Conversation) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("OwnerUserID is empty")
	}
	if x.ConversationID == "" {
		return errs.ErrArgs.Wrap("ConversationID is empty")
	}
	if x.ConversationType < 1 || x.ConversationType > 4 {
		return errs.ErrArgs.Wrap("ConversationType is invalid")
	}
	if x.RecvMsgOpt < 0 || x.RecvMsgOpt > 2 {
		return errs.ErrArgs.Wrap("RecvMsgOpt is invalid")
	}
	return nil
}

//func (x *ModifyConversationFieldReq) Check() error {
//	if x.UserIDList == nil {
//		return errs.ErrArgs.Wrap("userIDList is empty")
//	}
//	if x.Conversation == nil {
//		return errs.ErrArgs.Wrap("conversation is empty")
//	}
//	return nil
//}

func (x *SetConversationReq) Check() error {
	if x.Conversation == nil {
		return errs.ErrArgs.Wrap("Conversation is empty")
	}
	if x.Conversation.ConversationID == "" {
		return errs.ErrArgs.Wrap("conversationID is empty")
	}
	return nil
}

//func (x *SetRecvMsgOptReq) Check() error {
//	if x.OwnerUserID == "" {
//		return errs.ErrArgs.Wrap("ownerUserID is empty")
//	}
//	if x.ConversationID == "" {
//		return errs.ErrArgs.Wrap("conversationID is empty")
//	}
//	if x.RecvMsgOpt > 2 || x.RecvMsgOpt < 0 {
//		return errs.ErrArgs.Wrap("MsgReceiveOpt is invalid")
//	}
//	return nil
//}

func (x *GetConversationReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	if x.ConversationID == "" {
		return errs.ErrArgs.Wrap("conversationID is empty")
	}
	return nil
}

func (x *GetConversationsReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	if x.ConversationIDs == nil {
		return errs.ErrArgs.Wrap("conversationIDs is empty")
	}
	return nil
}

func (x *GetAllConversationsReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	return nil
}

//
//func (x *BatchSetConversationsReq) Check() error {
//	if x.Conversations == nil {
//		return errs.ErrArgs.Wrap("conversations is empty")
//	}
//	if x.OwnerUserID == "" {
//		return errs.ErrArgs.Wrap("conversation is empty")
//	}
//	return nil
//}

func (x *GetRecvMsgNotNotifyUserIDsReq) Check() error {
	if x.GroupID == "" {
		return errs.ErrArgs.Wrap("groupID is empty")
	}
	return nil
}

func (x *CreateGroupChatConversationsReq) Check() error {
	if x.GroupID == "" {
		return errs.ErrArgs.Wrap("groupID is empty")
	}
	return nil
}

func (x *SetConversationMaxSeqReq) Check() error {
	if x.ConversationID == "" {
		return errs.ErrArgs.Wrap("conversationID is empty")
	}
	if x.OwnerUserID == nil {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	if x.MaxSeq <= 0 {
		return errs.ErrArgs.Wrap("maxSeq is invalid")
	}
	return nil
}

func (x *SetConversationsReq) Check() error {
	if x.UserIDs == nil {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.Conversation == nil {
		return errs.ErrArgs.Wrap("conversation is empty")
	}
	return nil
}

func (x *GetUserConversationIDsHashReq) Check() error {
	if x.OwnerUserID == "" {
		return errs.ErrArgs.Wrap("ownerUserID is empty")
	}
	return nil
}

func (x *GetConversationsByConversationIDReq) Check() error {
	if x.ConversationIDs == nil {
		return errs.ErrArgs.Wrap("conversationIDs is empty")
	}
	return nil
}

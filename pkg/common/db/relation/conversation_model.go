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

package relation

import (
	"context"

	"gorm.io/gorm"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type ConversationGorm struct {
	*MetaDB
}

func NewConversationGorm(db *gorm.DB) relation.ConversationModelInterface {
	return &ConversationGorm{NewMetaDB(db, &relation.ConversationModel{})}
}

func (c *ConversationGorm) NewTx(tx any) relation.ConversationModelInterface {
	return &ConversationGorm{NewMetaDB(tx.(*gorm.DB), &relation.ConversationModel{})}
}

func (c *ConversationGorm) Create(ctx context.Context, conversations []*relation.ConversationModel) (err error) {
	return utils.Wrap(c.db(ctx).Create(&conversations).Error, "")
}

func (c *ConversationGorm) Delete(ctx context.Context, groupIDs []string) (err error) {
	return utils.Wrap(c.db(ctx).Where("group_id in (?)", groupIDs).Delete(&relation.ConversationModel{}).Error, "")
}

func (c *ConversationGorm) UpdateByMap(
	ctx context.Context,
	userIDList []string,
	conversationID string,
	args map[string]interface{},
) (rows int64, err error) {
	result := c.db(ctx).Where("owner_user_id IN (?) and  conversation_id=?", userIDList, conversationID).Updates(args)
	return result.RowsAffected, utils.Wrap(result.Error, "")
}

func (c *ConversationGorm) Update(ctx context.Context, conversation *relation.ConversationModel) (err error) {
	return utils.Wrap(
		c.db(ctx).
			Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
			Updates(conversation).
			Error,
		"",
	)
}

func (c *ConversationGorm) Find(
	ctx context.Context,
	ownerUserID string,
	conversationIDs []string,
) (conversations []*relation.ConversationModel, err error) {
	err = utils.Wrap(
		c.db(ctx).
			Where("owner_user_id=? and conversation_id IN (?)", ownerUserID, conversationIDs).
			Find(&conversations).
			Error,
		"",
	)
	return conversations, err
}

func (c *ConversationGorm) Take(
	ctx context.Context,
	userID, conversationID string,
) (conversation *relation.ConversationModel, err error) {
	cc := &relation.ConversationModel{}
	return cc, utils.Wrap(
		c.db(ctx).Where("conversation_id = ? And owner_user_id = ?", conversationID, userID).Take(cc).Error,
		"",
	)
}

func (c *ConversationGorm) FindUserID(
	ctx context.Context,
	userIDs []string,
	conversationIDs []string,
) (existUserID []string, err error) {
	return existUserID, utils.Wrap(
		c.db(ctx).
			Where(" owner_user_id IN (?) and conversation_id in (?)", userIDs, conversationIDs).
			Pluck("owner_user_id", &existUserID).
			Error,
		"",
	)
}

func (c *ConversationGorm) FindConversationID(
	ctx context.Context,
	userID string,
	conversationIDList []string,
) (existConversationID []string, err error) {
	return existConversationID, utils.Wrap(
		c.db(ctx).
			Where(" conversation_id IN (?) and owner_user_id=?", conversationIDList, userID).
			Pluck("conversation_id", &existConversationID).
			Error,
		"",
	)
}

func (c *ConversationGorm) FindUserIDAllConversationID(
	ctx context.Context,
	userID string,
) (conversationIDList []string, err error) {
	return conversationIDList, utils.Wrap(
		c.db(ctx).Where("owner_user_id=?", userID).Pluck("conversation_id", &conversationIDList).Error,
		"",
	)
}

func (c *ConversationGorm) FindUserIDAllConversations(
	ctx context.Context,
	userID string,
) (conversations []*relation.ConversationModel, err error) {
	return conversations, utils.Wrap(c.db(ctx).Where("owner_user_id=?", userID).Find(&conversations).Error, "")
}

func (c *ConversationGorm) FindRecvMsgNotNotifyUserIDs(
	ctx context.Context,
	groupID string,
) (userIDs []string, err error) {
	return userIDs, utils.Wrap(
		c.db(ctx).
			Where("group_id = ? and recv_msg_opt = ?", groupID, constant.ReceiveNotNotifyMessage).
			Pluck("user_id", &userIDs).
			Error,
		"",
	)
}

func (c *ConversationGorm) FindSuperGroupRecvMsgNotNotifyUserIDs(
	ctx context.Context,
	groupID string,
) (userIDs []string, err error) {
	return userIDs, utils.Wrap(
		c.db(ctx).
			Where("group_id = ? and recv_msg_opt = ? and conversation_type = ?", groupID, constant.ReceiveNotNotifyMessage, constant.SuperGroupChatType).
			Pluck("user_id", &userIDs).
			Error,
		"",
	)
}

func (c *ConversationGorm) GetUserRecvMsgOpt(
	ctx context.Context,
	ownerUserID, conversationID string,
) (opt int, err error) {
	var conversation relation.ConversationModel
	return int(
			conversation.RecvMsgOpt,
		), utils.Wrap(
			c.db(ctx).
				Where("conversation_id = ? And owner_user_id = ?", conversationID, ownerUserID).
				Select("recv_msg_opt").
				Find(&conversation).
				Error,
			"",
		)
}

func (c *ConversationGorm) GetAllConversationIDs(ctx context.Context) (conversationIDs []string, err error) {
	return conversationIDs, utils.Wrap(
		c.db(ctx).Distinct("conversation_id").Pluck("conversation_id", &conversationIDs).Error,
		"",
	)
}

func (c *ConversationGorm) GetUserAllHasReadSeqs(
	ctx context.Context,
	ownerUserID string,
) (hasReadSeqs map[string]int64, err error) {
	return nil, nil
}

func (c *ConversationGorm) GetConversationsByConversationID(
	ctx context.Context,
	conversationIDs []string,
) (conversations []*relation.ConversationModel, err error) {
	return conversations, utils.Wrap(
		c.db(ctx).Where("conversation_id IN (?)", conversationIDs).Find(&conversations).Error,
		"",
	)
}

func (c *ConversationGorm) GetConversationIDsNeedDestruct(
	ctx context.Context,
) (conversations []*relation.ConversationModel, err error) {
	return conversations, utils.Wrap(
		c.db(ctx).
			Where("is_msg_destruct = 1 && UNIX_TIMESTAMP(NOW()) > (msg_destruct_time + UNIX_TIMESTAMP(latest_msg_destruct_time)) && msg_destruct_time != 0").
			Find(&conversations).
			Error,
		"",
	)
}

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

package unrelation

import (
	"context"
	"fmt"

	"github.com/OpenIMSDK/tools/log"
	table "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/unrelation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MsgMongoDriver) ConvertMsgsDocLen(ctx context.Context, conversationIDs []string) {
	for _, conversationID := range conversationIDs {
		regex := primitive.Regex{Pattern: fmt.Sprintf("^%s:", conversationID)}
		cursor, err := m.MsgCollection.Find(ctx, bson.M{"doc_id": regex})
		if err != nil {
			log.ZError(ctx, "convertAll find msg doc failed", err, "conversationID", conversationID)
			continue
		}
		var msgDocs []table.MsgDocModel
		err = cursor.All(ctx, &msgDocs)
		if err != nil {
			log.ZError(ctx, "convertAll cursor all failed", err, "conversationID", conversationID)
			continue
		}
		if len(msgDocs) < 1 {
			continue
		}
		log.ZInfo(ctx, "msg doc convert", "conversationID", conversationID, "len(msgDocs)", len(msgDocs))
		if len(msgDocs[0].Msg) == int(m.model.GetSingleGocMsgNum5000()) {
			if _, err := m.MsgCollection.DeleteMany(ctx, bson.M{"doc_id": regex}); err != nil {
				log.ZError(ctx, "convertAll delete many failed", err, "conversationID", conversationID)
				continue
			}
			var newMsgDocs []any
			for _, msgDoc := range msgDocs {
				if int64(len(msgDoc.Msg)) == m.model.GetSingleGocMsgNum() {
					continue
				}
				var index int64
				for index < int64(len(msgDoc.Msg)) {
					msg := msgDoc.Msg[index]
					if msg != nil && msg.Msg != nil {
						msgDocModel := table.MsgDocModel{DocID: m.model.GetDocID(conversationID, msg.Msg.Seq)}
						end := index + m.model.GetSingleGocMsgNum()
						if int(end) >= len(msgDoc.Msg) {
							msgDocModel.Msg = msgDoc.Msg[index:]
						} else {
							msgDocModel.Msg = msgDoc.Msg[index:end]
						}
						newMsgDocs = append(newMsgDocs, msgDocModel)
						index = end
					} else {
						break
					}
				}
			}
			_, err = m.MsgCollection.InsertMany(ctx, newMsgDocs)
			if err != nil {
				log.ZError(ctx, "convertAll insert many failed", err, "conversationID", conversationID, "len(newMsgDocs)", len(newMsgDocs))
			} else {
				log.ZInfo(ctx, "msg doc convert", "conversationID", conversationID, "len(newMsgDocs)", len(newMsgDocs))
			}
		}
	}
}

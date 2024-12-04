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

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Friend represents the data structure for a friend relationship in MongoDB.
type Friend struct {
	ID             primitive.ObjectID `bson:"_id"`
	OwnerUserID    string             `bson:"owner_user_id"`
	FriendUserID   string             `bson:"friend_user_id"`
	Remark         string             `bson:"remark"`
	CreateTime     time.Time          `bson:"create_time"`
	AddSource      int32              `bson:"add_source"`
	OperatorUserID string             `bson:"operator_user_id"`
	Ex             string             `bson:"ex"`
	IsPinned       bool               `bson:"is_pinned"`
}

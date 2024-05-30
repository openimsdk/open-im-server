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
)

type Object struct {
	Name        string    `bson:"name"`
	UserID      string    `bson:"user_id"`
	Hash        string    `bson:"hash"`
	Engine      string    `bson:"engine"`
	Key         string    `bson:"key"`
	Size        int64     `bson:"size"`
	ContentType string    `bson:"content_type"`
	Group       string    `bson:"group"`
	CreateTime  time.Time `bson:"create_time"`
}

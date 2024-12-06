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

type Log struct {
	LogID        string    `bson:"log_id"`
	Platform     string    `bson:"platform"`
	UserID       string    `bson:"user_id"`
	CreateTime   time.Time `bson:"create_time"`
	Url          string    `bson:"url"`
	FileName     string    `bson:"file_name"`
	SystemType   string    `bson:"system_type"`
	AppFramework string    `bson:"app_framework"`
	Version      string    `bson:"version"`
	Ex           string    `bson:"ex"`
}

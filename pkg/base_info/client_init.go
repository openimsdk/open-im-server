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

package base_info

type SetClientInitConfigReq struct {
	OperationID     string  `json:"operationID"  binding:"required"`
	DiscoverPageURL *string `json:"discoverPageURL"`
}

type SetClientInitConfigResp struct {
	CommResp
}

type GetClientInitConfigReq struct {
	OperationID string `json:"operationID"  binding:"required"`
}

type GetClientInitConfigResp struct {
	CommResp
	Data struct {
		DiscoverPageURL string `json:"discoverPageURL"`
	} `json:"data"`
}

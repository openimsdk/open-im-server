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

package call_back_struct

import commonPb "Open_IM/pkg/proto/sdk_ws"

type CallbackBeforePushReq struct {
	UserStatusBatchCallbackReq
	*commonPb.OfflinePushInfo
	ClientMsgID  string   `json:"clientMsgID"`
	SendID       string   `json:"sendID"`
	GroupID      string   `json:"groupID"`
	ContentType  int32    `json:"contentType"`
	SessionType  int32    `json:"sessionType"`
	AtUserIDList []string `json:"atUserIDList"`
	Content      string   `json:"content"`
}

type CallbackBeforePushResp struct {
	*CommonCallbackResp
	UserIDList      []string                  `json:"userIDList"`
	OfflinePushInfo *commonPb.OfflinePushInfo `json:"offlinePushInfo"`
}

type CallbackBeforeSuperGroupOnlinePushReq struct {
	//*commonPb.OfflinePushInfo
	UserStatusBaseCallback
	ClientMsgID  string   `json:"clientMsgID"`
	SendID       string   `json:"sendID"`
	GroupID      string   `json:"groupID"`
	ContentType  int32    `json:"contentType"`
	SessionType  int32    `json:"sessionType"`
	AtUserIDList []string `json:"atUserIDList"`
	Content      string   `json:"content"`
	Seq          uint32   `json:"seq"`
}

type CallbackBeforeSuperGroupOnlinePushResp struct {
	*CommonCallbackResp
	UserIDList      []string                  `json:"userIDList"`
	OfflinePushInfo *commonPb.OfflinePushInfo `json:"offlinePushInfo"`
}

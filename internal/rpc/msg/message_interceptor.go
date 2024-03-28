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

package msg

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
)

type MessageInterceptorFunc func(ctx context.Context, globalConfig *config.GlobalConfig, req *msg.SendMsgReq) (*sdkws.MsgData, error)

func MessageHasReadEnabled(_ context.Context, globalConfig *config.GlobalConfig, req *msg.SendMsgReq) (*sdkws.MsgData, error) {
	switch {
	case req.MsgData.ContentType == constant.HasReadReceipt && req.MsgData.SessionType == constant.SingleChatType:
		if !globalConfig.SingleMessageHasReadReceiptEnable {
			return nil, servererrs.ErrMessageHasReadDisable.Wrap()
		}
		return req.MsgData, nil
	case req.MsgData.ContentType == constant.HasReadReceipt && req.MsgData.SessionType == constant.SuperGroupChatType:
		if !globalConfig.GroupMessageHasReadReceiptEnable {
			return nil, servererrs.ErrMessageHasReadDisable.Wrap()
		}
		return req.MsgData, nil
	}
	return req.MsgData, nil
}

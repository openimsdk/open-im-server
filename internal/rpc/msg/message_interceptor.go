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

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type MessageInterceptorFunc func(ctx context.Context, req *msg.SendMsgReq) (*sdkws.MsgData, error)

func MessageHasReadEnabled(_ context.Context, req *msg.SendMsgReq) (*sdkws.MsgData, error) {
	switch {
	case req.MsgData.ContentType == constant.HasReadReceipt && req.MsgData.SessionType == constant.SingleChatType:
		if config.Config.SingleMessageHasReadReceiptEnable {
			return req.MsgData, nil
		} else {
			return nil, errs.ErrMessageHasReadDisable.Wrap()
		}
	case req.MsgData.ContentType == constant.HasReadReceipt && req.MsgData.SessionType == constant.SuperGroupChatType:
		if config.Config.GroupMessageHasReadReceiptEnable {
			return req.MsgData, nil
		} else {
			return nil, errs.ErrMessageHasReadDisable.Wrap()
		}
	}
	return req.MsgData, nil
}

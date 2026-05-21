// Copyright © 2024 OpenIM. All rights reserved.
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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/rtc"
	"github.com/openimsdk/tools/a2r"
)

type RtcApi struct {
	Client rtc.RtcServiceClient
}

func NewRtcApi(client rtc.RtcServiceClient) RtcApi {
	return RtcApi{Client: client}
}

func (o *RtcApi) SignalMessageAssemble(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.SignalMessageAssemble, o.Client)
}

func (o *RtcApi) SignalGetRoomByGroupID(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.SignalGetRoomByGroupID, o.Client)
}

func (o *RtcApi) SignalGetTokenByRoomID(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.SignalGetTokenByRoomID, o.Client)
}

func (o *RtcApi) SignalGetRooms(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.SignalGetRooms, o.Client)
}

func (o *RtcApi) GetSignalInvitationInfo(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.GetSignalInvitationInfo, o.Client)
}

func (o *RtcApi) GetSignalInvitationInfoStartApp(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.GetSignalInvitationInfoStartApp, o.Client)
}

func (o *RtcApi) SignalSendCustomSignal(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.SignalSendCustomSignal, o.Client)
}

func (o *RtcApi) GetSignalInvitationRecords(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.GetSignalInvitationRecords, o.Client)
}

func (o *RtcApi) DeleteSignalRecords(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.DeleteSignalRecords, o.Client)
}

func (o *RtcApi) GetCallRecords(c *gin.Context) {
	a2r.Call(c, rtc.RtcServiceClient.GetCallRecords, o.Client)
}

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

package getui

import (
	"fmt"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (r *Resp) parseError() (err error) {
	switch r.Code {
	case tokenExpireCode:
		err = ErrTokenExpire
	case 0:
		err = nil
	default:
		err = fmt.Errorf("code %d, msg %s", r.Code, r.Msg)
	}
	return err
}

type RespI interface {
	parseError() error
}

type AuthReq struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	AppKey    string `json:"appkey"`
}

type AuthResp struct {
	ExpireTime string `json:"expire_time"`
	Token      string `json:"token"`
}

type TaskResp struct {
	TaskID string `json:"taskID"`
}

type Settings struct {
	TTL *int64 `json:"ttl"`
}

type Audience struct {
	Alias []string `json:"alias"`
}

type PushMessage struct {
	Notification *Notification `json:"notification,omitempty"`
	Transmission *string       `json:"transmission,omitempty"`
}

type PushChannel struct {
	Ios     *Ios     `json:"ios"`
	Android *Android `json:"android"`
}

type PushReq struct {
	RequestID   *string      `json:"request_id"`
	Settings    *Settings    `json:"settings"`
	Audience    *Audience    `json:"audience"`
	PushMessage *PushMessage `json:"push_message"`
	PushChannel *PushChannel `json:"push_channel"`
	IsAsync     *bool        `json:"is_async"`
	TaskID      *string      `json:"taskid"`
}

type Ios struct {
	NotificationType *string `json:"type"`
	AutoBadge        *string `json:"auto_badge"`
	Aps              struct {
		Sound string `json:"sound"`
		Alert Alert  `json:"alert"`
	} `json:"aps"`
}

type Alert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Android struct {
	Ups struct {
		Notification Notification `json:"notification"`
		Options      Options      `json:"options"`
	} `json:"ups"`
}

type Notification struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	ChannelID   string `json:"channelID"`
	ChannelName string `json:"ChannelName"`
	ClickType   string `json:"click_type"`
}

type Options struct {
	HW struct {
		DefaultSound bool   `json:"/message/android/notification/default_sound"`
		ChannelID    string `json:"/message/android/notification/channel_id"`
		Sound        string `json:"/message/android/notification/sound"`
		Importance   string `json:"/message/android/notification/importance"`
	} `json:"HW"`
	XM struct {
		ChannelID string `json:"/extra.channel_id"`
	} `json:"XM"`
	VV struct {
		Classification int `json:"/classification"`
	} `json:"VV"`
}

type Payload struct {
	IsSignal bool `json:"isSignal"`
}

func newPushReq(title, content string) PushReq {
	pushReq := PushReq{PushMessage: &PushMessage{Notification: &Notification{
		Title:       title,
		Body:        content,
		ClickType:   "startapp",
		ChannelID:   config.Config.Push.GeTui.ChannelID,
		ChannelName: config.Config.Push.GeTui.ChannelName,
	}}}
	return pushReq
}

func newBatchPushReq(userIDs []string, taskID string) PushReq {
	IsAsync := true
	return PushReq{Audience: &Audience{Alias: userIDs}, IsAsync: &IsAsync, TaskID: &taskID}
}

func (pushReq *PushReq) setPushChannel(title string, body string) {
	pushReq.PushChannel = &PushChannel{}
	// autoBadge := "+1"
	pushReq.PushChannel.Ios = &Ios{}
	notify := "notify"
	pushReq.PushChannel.Ios.NotificationType = &notify
	pushReq.PushChannel.Ios.Aps.Sound = "default"
	pushReq.PushChannel.Ios.Aps.Alert = Alert{
		Title: title,
		Body:  body,
	}
	pushReq.PushChannel.Android = &Android{}
	pushReq.PushChannel.Android.Ups.Notification = Notification{
		Title:     title,
		Body:      body,
		ClickType: "startapp",
	}
	pushReq.PushChannel.Android.Ups.Options = Options{
		HW: struct {
			DefaultSound bool   `json:"/message/android/notification/default_sound"`
			ChannelID    string `json:"/message/android/notification/channel_id"`
			Sound        string `json:"/message/android/notification/sound"`
			Importance   string `json:"/message/android/notification/importance"`
		}{ChannelID: "RingRing4", Sound: "/raw/ring001", Importance: "NORMAL"},
		XM: struct {
			ChannelID string `json:"/extra.channel_id"`
		}{ChannelID: "high_system"},
		VV: struct {
			Classification int "json:\"/classification\""
		}{
			Classification: 1,
		},
	}
}

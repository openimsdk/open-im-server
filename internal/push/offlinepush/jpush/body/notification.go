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

package body

import (
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/options"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type Notification struct {
	Alert   string  `json:"alert,omitempty"`
	Android Android `json:"android,omitempty"`
	IOS     Ios     `json:"ios,omitempty"`
}

type Android struct {
	Alert  string `json:"alert,omitempty"`
	Title  string `json:"title,omitempty"`
	Intent struct {
		URL string `json:"url,omitempty"`
	} `json:"intent,omitempty"`
	Extras map[string]string `json:"extras,omitempty"`
}
type Ios struct {
	Alert          IosAlert          `json:"alert,omitempty"`
	Sound          string            `json:"sound,omitempty"`
	Badge          string            `json:"badge,omitempty"`
	Extras         map[string]string `json:"extras,omitempty"`
	MutableContent bool              `json:"mutable-content"`
}

type IosAlert struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

func (n *Notification) SetAlert(alert string, title string, opts *options.Opts) {
	n.Alert = alert
	n.Android.Alert = alert
	n.Android.Title = title
	n.IOS.Alert.Body = alert
	n.IOS.Alert.Title = title
	n.IOS.Sound = opts.IOSPushSound
	if opts.IOSBadgeCount {
		n.IOS.Badge = "+1"
	}
}

func (n *Notification) SetExtras(extras map[string]string) {
	n.IOS.Extras = extras
	n.Android.Extras = extras
}

func (n *Notification) SetAndroidIntent(pushConf *config.Push) {
	n.Android.Intent.URL = pushConf.JPush.PushIntent
}

func (n *Notification) IOSEnableMutableContent() {
	n.IOS.MutableContent = true
}

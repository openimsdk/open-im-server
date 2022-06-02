package requestBody

import (
	"Open_IM/pkg/common/config"
)

type Notification struct {
	Alert   string  `json:"alert,omitempty"`
	Android Android `json:"android,omitempty"`
	IOS     Ios     `json:"ios,omitempty"`
}

type Android struct {
	Alert  string `json:"alert,omitempty"`
	Intent struct {
		URL string `json:"url,omitempty"`
	} `json:"intent,omitempty"`
	Extras Extras `json:"extras"`
}
type Ios struct {
	Alert          string `json:"alert,omitempty"`
	Sound          string `json:"sound,omitempty"`
	Badge          string `json:"badge,omitempty"`
	Extras         Extras `json:"extras"`
	MutableContent bool   `json:"mutable-content"`
}

type Extras struct {
	ClientMsgID string `json:"clientMsgID"`
}

func (n *Notification) SetAlert(alert string) {
	n.Alert = alert
	n.Android.Alert = alert
	n.SetAndroidIntent()
	n.IOS.Alert = alert
	n.IOS.Sound = "default"
	n.IOS.Badge = "+1"
}

func (n *Notification) SetExtras(extras Extras) {
	n.IOS.Extras = extras
	n.Android.Extras = extras
}

func (n *Notification) SetAndroidIntent() {
	n.Android.Intent.URL = config.Config.Push.Jpns.PushIntent
}

func (n *Notification) IOSEnableMutableContent() {
	n.IOS.MutableContent = true
}

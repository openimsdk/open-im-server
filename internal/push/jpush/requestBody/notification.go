package requestBody

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
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
}
type Ios struct {
	Alert string `json:"alert,omitempty"`
	Sound string `json:"sound,omitempty"`
	Badge string `json:"badge,omitempty"`
}

func (n *Notification) SetAlert(alert, platform string) {
	n.Alert = alert
	switch platform {
	case constant.AndroidPlatformStr:
		n.Android.Alert = alert
		n.SetAndroidIntent()
	case constant.IOSPlatformStr:
		n.IOS.Alert = alert
		n.IOS.Sound = "default"
		n.IOS.Badge = "+1"
	default:
	}
}
func (n *Notification) SetAndroidIntent() {
	n.Android.Intent.URL = config.Config.Push.Jpns.PushIntent
}

package requestBody

import "Open_IM/pkg/common/config"

type Notification struct {
	Alert   string   `json:"alert,omitempty"`
	Android *Android `json:"android,omitempty"`
	IOS     *Ios     `json:"ios,omitempty"`
}

type Android struct {
	Intent string `json:"intent,omitempty"`
}
type Ios struct {
}

func (n *Notification) SetAlert(alert string) {
	n.Alert = alert
}
func (n *Notification) SetAndroidIntent() {
	n.Android.Intent = config.Config.Push.Jpns.PushIntent
}

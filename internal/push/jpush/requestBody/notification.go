package requestBody

const INTENT_URL = "intent:#Intent;component=io.openim.app.enterprisechat/io.openim.app.enterprisechat.MainActivity;end"

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
	n.Android.Intent = INTENT_URL
}

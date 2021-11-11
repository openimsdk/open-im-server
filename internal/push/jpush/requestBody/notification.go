package requestBody

type Notification struct {
	Alert   string   `json:"alert,omitempty"`
	Android *Android `json:"android,omitempty"`
	IOS     *Ios     `json:"ios,omitempty"`
}

type Android struct {
}

type Ios struct {
}

func (n *Notification) SetAlert(alert string) {
	n.Alert = alert
}

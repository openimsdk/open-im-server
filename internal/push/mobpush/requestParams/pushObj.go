package requestParams

type PushObj struct {
	Source      interface{} `json:"source"`
	Appkey      interface{} `json:"appkey"`
	PushTarget  interface{} `json:"pushTarget,omitempty"`
	PushNotify  interface{} `json:"pushNotify,omitempty"`
	PushForward interface{} `json:"pushForward,omitempty"`
}

func (p *PushObj) SetSource(source string) {
	p.Source = source
}

func (p *PushObj) SetAppkey(appkey string) {
	p.Appkey = appkey
}

func (p *PushObj) SetPushTarget(no *PushTarget) {
	p.PushTarget = no
}

func (p *PushObj) SetPushNotify(m *PushNotify) {
	p.PushNotify = m
}
func (p *PushObj) SetPushForward(o *PushForward) {
	p.PushForward = o
}

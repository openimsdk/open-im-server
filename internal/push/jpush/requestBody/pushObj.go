package requestBody

type PushObj struct {
	Platform     interface{} `json:"platform"`
	Audience     interface{} `json:"audience"`
	Notification interface{} `json:"notification,omitempty"`
	Message      interface{} `json:"message,omitempty"`
	Options      interface{} `json:"options,omitempty"`
}

func (p *PushObj) SetPlatform(pf *Platform) {
	p.Platform = pf.Os
}

func (p *PushObj) SetAudience(ad *Audience) {
	p.Audience = ad.Object
}

func (p *PushObj) SetNotification(no *Notification) {
	p.Notification = no
}

func (p *PushObj) SetMessage(m *Message) {
	p.Message = m
}
func (p *PushObj) SetOptions(o *Options) {
	p.Options = o
}

package body

type PushObj struct {
	Platform     any `json:"platform"`
	Audience     any `json:"audience"`
	Notification any `json:"notification,omitempty"`
	Message      any `json:"message,omitempty"`
	Options      any `json:"options,omitempty"`
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

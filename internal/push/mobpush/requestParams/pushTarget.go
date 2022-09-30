package requestParams

type PushTarget struct {
	Target interface{} `json:"target,omitempty"`
	Alias  []string    `json:"alias,omitempty"`
}

func (p *PushTarget) SetTarget(target int) {
	p.Target = target
}
func (p *PushTarget) SetAlias(alias []string) {
	p.Alias = alias
}

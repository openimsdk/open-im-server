package option

func NewOption() *Option {
	return &Option{}
}

type Option struct {
	Enable *bool
	Link   []string
}

func (o *Option) WithEnable() *Option {
	t := true
	o.Enable = &t
	return o
}

func (o *Option) WithDisable() *Option {
	f := false
	o.Enable = &f
	return o
}

func (o *Option) WithLink(key ...string) *Option {
	if len(key) > 0 {
		if len(o.Link) == 0 {
			o.Link = key
		} else {
			o.Link = append(o.Link, key...)
		}
	}
	return o
}

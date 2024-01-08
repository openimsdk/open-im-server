package option

var (
	t = true
	f = false
)

type Option struct {
	enable *bool
	key    []string
}

func (o *Option) Enable() *Option {
	o.enable = &t
	return o
}

func (o *Option) Disable() *Option {
	o.enable = &f
	return o
}

func (o *Option) DelKey(key ...string) *Option {
	if len(key) > 0 {
		if o.key == nil {
			o.key = key
		} else {
			o.key = append(o.key, key...)
		}
	}
	return o
}

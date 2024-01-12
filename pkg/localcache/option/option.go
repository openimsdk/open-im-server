package option

func NewOption() *Option {
	return &Option{}
}

type Option struct {
	Link []string
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

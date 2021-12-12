package requestBody

type Options struct {
	ApnsProduction bool `json:"apns_production"`
}

func (o *Options) SetApnsProduction(c bool) {
	o.ApnsProduction = c
}

package requestBody

type Options struct {
	ApnsProduction bool `json:"apns_production,omitempty"`
}

func (o *Options) SetApnsProduction(c bool) {
	o.ApnsProduction = c
}

package requestParams

type PushForward struct {
	NextType int    `json:"nextType"`
	Scheme   string `json:"scheme,omitempty"`
}

func (m *PushForward) SetNextType(c int) {
	m.NextType = c
}

func (m *PushForward) SetScheme(t string) {
	m.Scheme = t
}

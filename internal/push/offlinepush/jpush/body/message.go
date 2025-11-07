package body

type Message struct {
	MsgContent  string         `json:"msg_content"`
	Title       string         `json:"title,omitempty"`
	ContentType string         `json:"content_type,omitempty"`
	Extras      map[string]any `json:"extras,omitempty"`
}

func (m *Message) SetMsgContent(c string) {
	m.MsgContent = c
}

func (m *Message) SetTitle(t string) {
	m.Title = t
}

func (m *Message) SetContentType(c string) {
	m.ContentType = c
}

func (m *Message) SetExtras(key string, value any) {
	if m.Extras == nil {
		m.Extras = make(map[string]any)
	}
	m.Extras[key] = value
}

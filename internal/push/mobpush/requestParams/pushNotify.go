package requestParams

type PushNotify struct {
	Plats         []int  `json:"plats,omitempty"`
	IosProduction int    `json:"iosProduction,omitempty"`
	Content       string `json:"content,omitempty"`
	Type          int    `json:"type,omitempty"`
}

func (n *PushNotify) SetPlats(plats []int) {
	n.Plats = plats

}

func (n *PushNotify) SetIosProduction(iosProduction int) {
	n.IosProduction = iosProduction

}

func (n *PushNotify) SetContent(content string) {
	n.Content = content
}

func (n *PushNotify) SetType(Type int) {
	n.Type = Type
}

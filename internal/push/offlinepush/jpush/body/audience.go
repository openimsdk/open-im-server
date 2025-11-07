package body

const (
	TAG            = "tag"
	TAGAND         = "tag_and"
	TAGNOT         = "tag_not"
	ALIAS          = "alias"
	REGISTRATIONID = "registration_id"
)

type Audience struct {
	Object   any
	audience map[string][]string
}

func (a *Audience) set(key string, v []string) {
	if a.audience == nil {
		a.audience = make(map[string][]string)
		a.Object = a.audience
	}
	// v, ok = this.audience[key]
	// if ok {
	//	return
	//}
	a.audience[key] = v
}

func (a *Audience) SetTag(tags []string) {
	a.set(TAG, tags)
}

func (a *Audience) SetTagAnd(tags []string) {
	a.set(TAGAND, tags)
}

func (a *Audience) SetTagNot(tags []string) {
	a.set(TAGNOT, tags)
}

func (a *Audience) SetAlias(alias []string) {
	a.set(ALIAS, alias)
}

func (a *Audience) SetRegistrationId(ids []string) {
	a.set(REGISTRATIONID, ids)
}

func (a *Audience) SetAll() {
	a.Object = "all"
}

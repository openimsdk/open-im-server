package requestBody

const (
	TAG             = "tag"
	TAG_AND         = "tag_and"
	TAG_NOT         = "tag_not"
	ALIAS           = "alias"
	REGISTRATION_ID = "registration_id"
	SEGMENT         = "segment"
	ABTEST          = "abtest"
)

type Audience struct {
	Object   interface{}
	audience map[string][]string
}

func (a *Audience) set(key string, v []string) {
	if a.audience == nil {
		a.audience = make(map[string][]string)
		a.Object = a.audience
	}

	//v, ok = this.audience[key]
	//if ok {
	//	return
	//}
	a.audience[key] = v
}

func (a *Audience) SetTag(tags []string) {
	a.set(TAG, tags)
}

func (a *Audience) SetTagAnd(tags []string) {
	a.set(TAG_AND, tags)
}

func (a *Audience) SetTagNot(tags []string) {
	a.set(TAG_NOT, tags)
}

func (a *Audience) SetAlias(alias []string) {
	a.set(ALIAS, alias)
}

func (a *Audience) SetRegistrationId(ids []string) {
	a.set(REGISTRATION_ID, ids)
}

func (a *Audience) SetAll() {
	a.Object = "all"
}

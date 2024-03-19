// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

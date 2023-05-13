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

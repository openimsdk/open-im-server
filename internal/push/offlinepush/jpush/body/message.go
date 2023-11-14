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

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

type PushObj struct {
	Platform     any `json:"platform"`
	Audience     any `json:"audience"`
	Notification any `json:"notification,omitempty"`
	Message      any `json:"message,omitempty"`
	Options      any `json:"options,omitempty"`
}

func (p *PushObj) SetPlatform(pf *Platform) {
	p.Platform = pf.Os
}

func (p *PushObj) SetAudience(ad *Audience) {
	p.Audience = ad.Object
}

func (p *PushObj) SetNotification(no *Notification) {
	p.Notification = no
}

func (p *PushObj) SetMessage(m *Message) {
	p.Message = m
}

func (p *PushObj) SetOptions(o *Options) {
	p.Options = o
}

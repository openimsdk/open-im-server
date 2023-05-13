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

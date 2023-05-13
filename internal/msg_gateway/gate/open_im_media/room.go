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

package open_im_media

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:11300"
)

type Media struct {
}

func NewMedia() *Media {
	return &Media{}
}

func init() {

}

func (m *Media) CreateRoom(roomName string) (error, error) {
	return nil, nil

}

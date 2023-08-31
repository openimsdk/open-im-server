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

package common

// MySQL
// V2
const (
	UsernameV2 = "root"
	PasswordV2 = "openIM"
	IpV2       = "121.5.182.23:13306"
	DatabaseV2 = "openIM_v2"
)

// V3
const (
	UsernameV3 = "root"
	PasswordV3 = "openIM123"
	IpV3       = "43.134.63.160:13306"
	DatabaseV3 = "openIM_v3"
)

// Kafka
const (
	Topic     = "ws2ms_chat"
	KafkaAddr = "127.0.0.1:9092"
)

// Zookeeper
const (
	ZkAddr     = "43.134.63.160:2181"
	ZKSchema   = "openim"
	ZKUsername = ""
	ZKPassword = ""
	MsgRpcName = "Msg"
)

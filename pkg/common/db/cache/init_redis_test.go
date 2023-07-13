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

package cache

import (
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"testing"
)

//TestNewRedis Test redis connection
func TestNewRedis(t *testing.T) {
	err := config.InitConfig("config_folder_path")
	if err != nil {
		fmt.Println("config load error")
		return
	}
	redis, err := NewRedis()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(redis)
}

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

package fcm

import (
	"Open_IM/internal/push"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Push(t *testing.T) {
	offlinePusher := NewFcm()
	resp, err := offlinePusher.Push([]string{"test_uid"}, "test", "test", "12321", push.PushOpts{})
	assert.Nil(t, err)
	fmt.Println(resp)
}

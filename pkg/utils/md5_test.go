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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Md5(t *testing.T) {
	result := Md5("go")
	assert.Equal(t, result, "34d1f91fb2e514b8576fab1a75a89a6b")

	result2 := Md5("go")
	assert.Equal(t, result, result2)
}

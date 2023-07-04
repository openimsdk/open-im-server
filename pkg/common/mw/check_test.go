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

package mw

import (
	"fmt"
	"testing"
)

func TestCheck(t *testing.T) {
	// config.Config.TokenPolicy.Secret = "123456"

	args := []string{"1", "2", "3"}

	key := genReqKey(args)
	fmt.Println("key:", key)
	err := verifyReqKey(args, key)

	fmt.Println(err)

	args = []string{"4", "5", "6"}

	key = genReqKey(args)
	fmt.Println("key:", key)
	err = verifyReqKey(args, key)

	fmt.Println(err)

}

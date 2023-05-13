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

package token_verify

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiJvcGVuSU1BZG1pbiIsIlBsYXRmb3JtIjoiQVBhZCIsImV4cCI6MTY3NDYxNTA2MSwibmJmIjoxNjY2ODM4NzYxLCJpYXQiOjE2NjY4MzkwNjF9.l8RiIu6pR4ItwDOpNIDYA9LBzIcpk8r8n6NRtXjqOp8"
	_, err := GetClaimFromToken(token)
	assert.Nil(t, err)
}

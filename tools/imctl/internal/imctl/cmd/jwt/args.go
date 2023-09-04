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

package jwt

import (
	"fmt"
	"strings"

	"github.com/openim-sigs/component-base/pkg/json"
)

// ArgList defines a new pflag Value.
type ArgList map[string]string

// String return value of ArgList in string format.
func (l ArgList) String() string {
	data, _ := json.Marshal(l)

	return string(data)
}

// Set sets the value of ArgList.
func (l ArgList) Set(arg string) error {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid argument '%v'. Must use format 'key=value'. %v", arg, parts)
	}
	l[parts[0]] = parts[1]

	return nil
}

// Type returns the type name of ArgList.
func (l ArgList) Type() string {
	return "map"
}

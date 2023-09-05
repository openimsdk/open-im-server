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

package main

import (
	"os"
	"testing"
)

func TestGetConfigValue(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		flagValue  string
		fallback   string
		wantResult string
	}{
		{
			name:       "environment variable set",
			envKey:     "TEST_KEY",
			envValue:   "envValue",
			flagValue:  "",
			fallback:   "default",
			wantResult: "envValue",
		},
		{
			name:       "flag set and environment variable not set",
			envKey:     "TEST_KEY",
			envValue:   "",
			flagValue:  "flagValue",
			fallback:   "default",
			wantResult: "flagValue",
		},
		{
			name:       "nothing set, use fallback",
			envKey:     "TEST_KEY",
			envValue:   "",
			flagValue:  "",
			fallback:   "default",
			wantResult: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			got := getConfigValue(tt.envKey, tt.flagValue, tt.fallback)

			if got != tt.wantResult {
				t.Errorf("getConfigValue(%s, %s, %s) = %s; want %s", tt.envKey, tt.flagValue, tt.fallback, got, tt.wantResult)
			}
		})
	}
}

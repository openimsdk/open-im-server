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
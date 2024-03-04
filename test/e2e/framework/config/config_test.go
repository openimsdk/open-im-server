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

package config

import (
	"flag"
	"reflect"
	"testing"
)

func TestCopyFlags(t *testing.T) {
	type args struct {
		source *flag.FlagSet
		target *flag.FlagSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Copy empty source to empty target",
			args: args{
				source: flag.NewFlagSet("source", flag.ContinueOnError),
				target: flag.NewFlagSet("target", flag.ContinueOnError),
			},
			wantErr: false,
		},
		{
			name: "Copy non-empty source to empty target",
			args: args{
				source: func() *flag.FlagSet {
					fs := flag.NewFlagSet("source", flag.ContinueOnError)
					fs.String("test-flag", "default", "test usage")
					return fs
				}(),
				target: flag.NewFlagSet("target", flag.ContinueOnError),
			},
			wantErr: false,
		},
		{
			name: "Copy source to target with existing flag",
			args: args{
				source: func() *flag.FlagSet {
					fs := flag.NewFlagSet("source", flag.ContinueOnError)
					fs.String("test-flag", "default", "test usage")
					return fs
				}(),
				target: func() *flag.FlagSet {
					fs := flag.NewFlagSet("target", flag.ContinueOnError)
					fs.String("test-flag", "default", "test usage")
					return fs
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantErr {
					t.Errorf("CopyFlags() panic = %v, wantErr %v", r, tt.wantErr)
				}
			}()
			CopyFlags(tt.args.source, tt.args.target)

			// Verify the replicated tag
			if !tt.wantErr {
				tt.args.source.VisitAll(func(f *flag.Flag) {
					if gotFlag := tt.args.target.Lookup(f.Name); gotFlag == nil || !reflect.DeepEqual(gotFlag, f) {
						t.Errorf("CopyFlags() failed to copy flag %s", f.Name)
					}
				})
			}
		})
	}
}

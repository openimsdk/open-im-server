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

package tls

import (
	"crypto/tls"
	"reflect"
	"testing"
)

func Test_decryptPEM(t *testing.T) {
	type args struct {
		data       []byte
		passphrase []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decryptPEM(tt.args.data, tt.args.passphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("decryptPEM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decryptPEM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readEncryptablePEMBlock(t *testing.T) {
	type args struct {
		path string
		pwd  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readEncryptablePEMBlock(tt.args.path, tt.args.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readEncryptablePEMBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readEncryptablePEMBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTLSConfig(t *testing.T) {
	type args struct {
		clientCertFile string
		clientKeyFile  string
		caCertFile     string
		keyPwd         []byte
	}
	tests := []struct {
		name string
		args args
		want *tls.Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTLSConfig(tt.args.clientCertFile, tt.args.clientKeyFile, tt.args.caCertFile, tt.args.keyPwd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTLSConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

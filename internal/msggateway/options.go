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

package msggateway

import "time"

type (
	Option  func(opt *configs)
	configs struct {
		// Long connection listening port
		port int
		// Maximum number of connections allowed for long connection
		maxConnNum int64
		// Connection handshake timeout
		handshakeTimeout time.Duration
		// Maximum length allowed for messages
		messageMaxMsgLength int
		// Websocket write buffer, default: 4096, 4kb.
		writeBufferSize int
	}
)

func WithPort(port int) Option {
	return func(opt *configs) {
		opt.port = port
	}
}

func WithMaxConnNum(num int64) Option {
	return func(opt *configs) {
		opt.maxConnNum = num
	}
}

func WithHandshakeTimeout(t time.Duration) Option {
	return func(opt *configs) {
		opt.handshakeTimeout = t
	}
}

func WithMessageMaxMsgLength(length int) Option {
	return func(opt *configs) {
		opt.messageMaxMsgLength = length
	}
}

func WithWriteBufferSize(size int) Option {
	return func(opt *configs) {
		opt.writeBufferSize = size
	}
}

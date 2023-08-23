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

// Package term provides structures and helper functions to work with
// terminal (state, sizes).
package term

import (
	"io"
)

// TTY helps invoke a function and preserve the state of the terminal, even if the process is
// terminated during execution. It also provides support for terminal resizing for remote command
// execution/attachment.
type TTY struct {
	// In is a reader representing stdin. It is a required field.
	In io.Reader
	// Out is a writer representing stdout. It must be set to support terminal resizing. It is an
	// optional field.
	Out io.Writer
	// Raw is true if the terminal should be set raw.
	Raw bool
	// TryDev indicates the TTY should try to open /dev/tty if the provided input
	// is not a file descriptor.
	TryDev bool
}

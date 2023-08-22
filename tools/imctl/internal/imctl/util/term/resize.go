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

package term

import (
	"github.com/moby/term"
)

// TerminalSize represents the width and height of a terminal.
type TerminalSize struct {
	Width  uint16
	Height uint16
}

// TerminalSizeQueue is capable of returning terminal resize events as they occur.
type TerminalSizeQueue interface {
	// Next returns the new terminal size after the terminal has been resized. It returns nil when
	// monitoring has been stopped.
	Next() *TerminalSize
}

// GetSize returns the current size of the user's terminal. If it isn't a terminal,
// nil is returned.
func (t TTY) GetSize() *TerminalSize {
	outFd, isTerminal := term.GetFdInfo(t.Out)
	if !isTerminal {
		return nil
	}
	return GetSize(outFd)
}

// GetSize returns the current size of the terminal associated with fd.
func GetSize(fd uintptr) *TerminalSize {
	winsize, err := term.GetWinsize(fd)
	if err != nil {
		// runtime.HandleError(fmt.Errorf("unable to get terminal size: %v", err))
		return nil
	}

	return &TerminalSize{Width: winsize.Width, Height: winsize.Height}
}

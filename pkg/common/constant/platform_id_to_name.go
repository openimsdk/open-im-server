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

package constant

// fixme 1<--->IOS 2<--->Android  3<--->Windows
// fixme  4<--->OSX  5<--->Web  6<--->MiniWeb 7<--->Linux.
const (
	// Platform ID.
	IOSPlatformID        = 1
	AndroidPlatformID    = 2
	WindowsPlatformID    = 3
	OSXPlatformID        = 4
	WebPlatformID        = 5
	MiniWebPlatformID    = 6
	LinuxPlatformID      = 7
	AndroidPadPlatformID = 8
	IPadPlatformID       = 9
	AdminPlatformID      = 10

	// Platform string match to Platform ID.
	IOSPlatformStr        = "IOS"
	AndroidPlatformStr    = "Android"
	WindowsPlatformStr    = "Windows"
	OSXPlatformStr        = "OSX"
	WebPlatformStr        = "Web"
	MiniWebPlatformStr    = "MiniWeb"
	LinuxPlatformStr      = "Linux"
	AndroidPadPlatformStr = "APad"
	IPadPlatformStr       = "IPad"
	AdminPlatformStr      = "Admin"

	// terminal types.
	TerminalPC     = "PC"
	TerminalMobile = "Mobile"
)

var PlatformID2Name = map[int]string{
	IOSPlatformID:        IOSPlatformStr,
	AndroidPlatformID:    AndroidPlatformStr,
	WindowsPlatformID:    WindowsPlatformStr,
	OSXPlatformID:        OSXPlatformStr,
	WebPlatformID:        WebPlatformStr,
	MiniWebPlatformID:    MiniWebPlatformStr,
	LinuxPlatformID:      LinuxPlatformStr,
	AndroidPadPlatformID: AndroidPadPlatformStr,
	IPadPlatformID:       IPadPlatformStr,
	AdminPlatformID:      AdminPlatformStr,
}

var PlatformName2ID = map[string]int{
	IOSPlatformStr:        IOSPlatformID,
	AndroidPlatformStr:    AndroidPlatformID,
	WindowsPlatformStr:    WindowsPlatformID,
	OSXPlatformStr:        OSXPlatformID,
	WebPlatformStr:        WebPlatformID,
	MiniWebPlatformStr:    MiniWebPlatformID,
	LinuxPlatformStr:      LinuxPlatformID,
	AndroidPadPlatformStr: AndroidPadPlatformID,
	IPadPlatformStr:       IPadPlatformID,
	AdminPlatformStr:      AdminPlatformID,
}

var PlatformName2class = map[string]string{
	IOSPlatformStr:     TerminalMobile,
	AndroidPlatformStr: TerminalMobile,
	MiniWebPlatformStr: WebPlatformStr,
	WebPlatformStr:     WebPlatformStr,
	WindowsPlatformStr: TerminalPC,
	OSXPlatformStr:     TerminalPC,
	LinuxPlatformStr:   TerminalPC,
}

var PlatformID2class = map[int]string{
	IOSPlatformID:     TerminalMobile,
	AndroidPlatformID: TerminalMobile,
	MiniWebPlatformID: WebPlatformStr,
	WebPlatformID:     WebPlatformStr,
	WindowsPlatformID: TerminalPC,
	OSXPlatformID:     TerminalPC,
	LinuxPlatformID:   TerminalPC,
}

func PlatformIDToName(num int) string {
	return PlatformID2Name[num]
}

func PlatformNameToID(name string) int {
	return PlatformName2ID[name]
}

func PlatformNameToClass(name string) string {
	return PlatformName2class[name]
}

func PlatformIDToClass(num int) string {
	return PlatformID2class[num]
}

package constant

// fixme 1<--->IOS 2<--->Android  3<--->Windows
//fixme  4<--->OSX  5<--->Web  6<--->MiniWeb 7<--->Linux

const (
	//Platform ID
	IOSPlatformID        = 1
	AndroidPlatformID    = 2
	WindowsPlatformID    = 3
	OSXPlatformID        = 4
	WebPlatformID        = 5
	MiniWebPlatformID    = 6
	LinuxPlatformID      = 7
	AndroidPadPlatformID = 8
	IPadPlatformID       = 9

	//Platform string match to Platform ID
	IOSPlatformStr        = "IOS"
	AndroidPlatformStr    = "Android"
	WindowsPlatformStr    = "Windows"
	OSXPlatformStr        = "OSX"
	WebPlatformStr        = "Web"
	MiniWebPlatformStr    = "MiniWeb"
	LinuxPlatformStr      = "Linux"
	AndroidPadPlatformStr = "APad"
	IPadPlatformStr       = "IPad"

	//terminal types
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
}
var Platform2class = map[string]string{
	IOSPlatformStr:     TerminalMobile,
	AndroidPlatformStr: TerminalMobile,
	MiniWebPlatformStr: WebPlatformStr,
	WebPlatformStr:     WebPlatformStr,
	WindowsPlatformStr: TerminalPC,
	OSXPlatformStr:     TerminalPC,
	LinuxPlatformStr:   TerminalPC,
}

func PlatformIDToName(num int) string {
	return PlatformID2Name[num]
}
func PlatformNameToID(name string) int {
	return PlatformName2ID[name]
}
func PlatformNameToClass(name string) string {
	return Platform2class[name]
}

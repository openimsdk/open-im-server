package constant

// fixme 1<--->IOS 2<--->Android  3<--->Windows
//fixme  4<--->OSX  5<--->Web  6<--->MiniWeb 7<--->Linux

const (
	//Platform ID
	IOSPlatformID     = 1
	AndroidPlatformID = 2
	WindowsPlatformID = 3
	OSXPlatformID     = 4
	WebPlatformID     = 5
	MiniWebPlatformID = 6
	LinuxPlatformID   = 7

	//Platform string match to Platform ID
	IOSPlatformStr     = "IOS"
	AndroidPlatformStr = "Android"
	WindowsPlatformStr = "Windows"
	OSXPlatformStr     = "OSX"
	WebPlatformStr     = "Web"
	MiniWebPlatformStr = "MiniWeb"
	LinuxPlatformStr   = "Linux"

	//terminal types
	TerminalPC     = "PC"
	TerminalMobile = "Mobile"
)

var PlatformID2Name = map[int32]string{
	IOSPlatformID:     IOSPlatformStr,
	AndroidPlatformID: AndroidPlatformStr,
	WindowsPlatformID: WindowsPlatformStr,
	OSXPlatformID:     OSXPlatformStr,
	WebPlatformID:     WebPlatformStr,
	MiniWebPlatformID: MiniWebPlatformStr,
	LinuxPlatformID:   LinuxPlatformStr,
}
var PlatformName2ID = map[string]int32{
	IOSPlatformStr:     IOSPlatformID,
	AndroidPlatformStr: AndroidPlatformID,
	WindowsPlatformStr: WindowsPlatformID,
	OSXPlatformStr:     OSXPlatformID,
	WebPlatformStr:     WebPlatformID,
	MiniWebPlatformStr: MiniWebPlatformID,
	LinuxPlatformStr:   LinuxPlatformID,
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

func PlatformIDToName(num int32) string {
	return PlatformID2Name[num]
}
func PlatformNameToID(name string) int32 {
	return PlatformName2ID[name]
}
func PlatformNameToClass(name string) string {
	return Platform2class[name]
}

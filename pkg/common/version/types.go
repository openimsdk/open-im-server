package version

// Info contains versioning information.
// TODO: Add []string of api versions supported? It's still unclear
// how we'll want to distribute that information.
type Info struct {
	Major      string `json:"major,omitempty"`
	Minor      string `json:"minor,omitempty"`
	GitVersion string `json:"gitVersion"`
	GitTreeState string `json:"gitTreeState,omitempty"`
	GitCommit  string `json:"gitCommit,omitempty"`
	BuildDate  string `json:"buildDate"`
	GoVersion  string `json:"goVersion"`
	Compiler   string `json:"compiler"`
	Platform   string `json:"platform"`
}

type Output struct {
	OpenIMServerVersion Info                 `json:"OpenIMServerVersion,omitempty" yaml:"OpenIMServerVersion,omitempty"`
	OpenIMClientVersion *OpenIMClientVersion `json:"OpenIMClientVersion,omitempty" yaml:"OpenIMClientVersion,omitempty"`
}

type OpenIMClientVersion struct {
	ClientVersion string `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"`	//sdk core version
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.GitVersion
}

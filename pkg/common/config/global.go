package config

var standalone bool

func SetStandalone() {
	standalone = true
}

func Standalone() bool {
	return standalone
}

package cachekey

const (
	ApplicationLatestVersion = "APPLICATION_LATEST_VERSION:"
)

func GetApplicationLatestVersionKey(platform string, hot bool) string {
	var hotStr string
	if hot {
		hotStr = "1:"
	} else {
		hotStr = "0:"
	}
	return ApplicationLatestVersion + hotStr + platform
}

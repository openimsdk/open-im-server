package cachekey

const (
	ApplicationLatestVersion = "APPLICATION_LATEST_VERSION:"
)

func GetApplicationLatestVersionKey(platform string) string {
	return ApplicationLatestVersion + platform
}

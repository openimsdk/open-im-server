package cachekey

const ClientConfig = "CLIENT_CONFIG"

func GetClientConfigKey(userID string) string {
	if userID == "" {
		return ClientConfig
	}
	return ClientConfig + ":" + userID
}

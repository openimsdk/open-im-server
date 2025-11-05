package discovery

import (
	"os"
)

func setupTestEnvironment() {
	os.Setenv("ZOOKEEPER_SCHEMA", "openim")
	os.Setenv("ZOOKEEPER_ADDRESS", "172.28.0.1")
	os.Setenv("ZOOKEEPER_PORT", "12181")
	os.Setenv("ZOOKEEPER_USERNAME", "")
	os.Setenv("ZOOKEEPER_PASSWORD", "")
}

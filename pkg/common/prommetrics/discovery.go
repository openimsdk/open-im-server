package prommetrics

import "fmt"

const (
	APIKeyName             = "api"
	MessageTransferKeyName = "message-transfer"
	TTL                    = 300
)

type Target struct {
	Target string            `json:"target"`
	Labels map[string]string `json:"labels"`
}

type RespTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func BuildDiscoveryKeyPrefix(name string) string {
	return fmt.Sprintf("%s/%s/%s", "openim", "prometheus_discovery", name)
}

func BuildDiscoveryKey(name string, index int) string {
	return fmt.Sprintf("%s/%s/%s/%d", "openim", "prometheus_discovery", name, index)
}

func BuildDefaultTarget(host string, ip int) Target {
	return Target{
		Target: fmt.Sprintf("%s:%d", host, ip),
		Labels: map[string]string{
			"namespace": "default",
		},
	}
}

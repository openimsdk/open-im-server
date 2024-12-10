package prommetrics

import "fmt"

const (
	APIKeyName             = "api"
	MessageTransferKeyName = "message-transfer"
)

type Target struct {
	Target string            `json:"target"`
	Labels map[string]string `json:"labels"`
}

type RespTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func BuildDiscoveryKey(name string) string {
	return fmt.Sprintf("%s/%s/%s", "openim", "prometheus_discovery", name)
}

func BuildDefaultTarget(host string, ip int) Target {
	return Target{
		Target: fmt.Sprintf("%s:%d", host, ip),
		Labels: map[string]string{
			"namespace": "default",
		},
	}
}

package disco

import (
	"net/url"

	"github.com/kuritka/k8gb-discovery/internal/services/discovery/data"
)

type Disco struct {
	yaml  *url.URL
	cache map[string]data.K8gb
}

func NewDisco(yaml *url.URL) *Disco {
	return &Disco{yaml: yaml}
}

func ReloadGitHubYaml(yaml *url.URL) {

}

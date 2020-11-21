// DepResolver provides configuration for particular services.
package depresolver

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/kuritka/k8gb-discovery/internal/services/discovery"

	"github.com/AbsaOSS/gopkg/env"
	"github.com/kuritka/k8gb-discovery/internal/common"
)

const (
	envPortKey = "K8GB_DISCOVERY_PORT"
	envYamlKey = "K8GB_DISCOVERY_YAML_URL"
)

type DepResolver struct {
	listener struct {
		initOnce sync.Once
		settings discovery.Settings
	}
}

func New() *DepResolver {
	dr := new(DepResolver)
	return dr
}

func (dr *DepResolver) MustResolveDiscovery() (s discovery.Settings, err error) {
	var e *multierror.Error
	dr.listener.initOnce.Do(func() {
		// resolve and validate all inputs here!
		dr.listener.settings.Port, err = env.GetEnvAsIntOrFallback(envPortKey, common.DefaultServicePort)
		e = multierror.Append(e, err)

		yamlURL := env.GetEnvAsStringOrFallback(envYamlKey, common.InvalidYamlURL)
		if yamlURL == common.InvalidYamlURL {
			e = multierror.Append(e, fmt.Errorf("invalid %s", envYamlKey))
		}
		dr.listener.settings.YamlURL, err = url.Parse(yamlURL)
		e = multierror.Append(e, err)
	})
	return dr.listener.settings, e.ErrorOrNil()
}

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
	envPortKey = "K8GB_DISCOVERY_EXPOSED_PORT"
	envYamlKey = "K8GB_DISCOVERY_YAML_URL"
	envSealedSecret = "SECRET_INFORMATION"
)

type DepResolver struct {
	discovery struct {
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
	dr.discovery.initOnce.Do(func() {
		// resolve and validate all inputs here!
		dr.discovery.settings.Port, err = env.GetEnvAsIntOrFallback(envPortKey, common.DefaultServicePort)
		e = multierror.Append(e, err)

		yamlURL := env.GetEnvAsStringOrFallback(envYamlKey, common.InvalidYamlURL)
		if yamlURL == common.InvalidYamlURL {
			e = multierror.Append(e, fmt.Errorf("invalid %s", envYamlKey))
		}
		dr.discovery.settings.YamlURL, err = url.Parse(yamlURL)
		e = multierror.Append(e, err)
		dr.discovery.settings.SealedSecret = env.GetEnvAsStringOrFallback(envSealedSecret,"secret "+envSealedSecret+" not found!")
	})
	return dr.discovery.settings, e.ErrorOrNil()
}

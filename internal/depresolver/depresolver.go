// DepResolver provides configuration for particular services.
package depresolver

import (
	"sync"

	"github.com/AbsaOSS/gopkg/env"
	"github.com/kuritka/k8gb-discovery/internal/cmd/listener"
	"github.com/kuritka/k8gb-discovery/internal/common"
)

type DepResolver struct {
	listener struct {
		initOnce sync.Once
		settings listener.Settings
	}
}

func New() *DepResolver {
	dr := new(DepResolver)
	return dr
}

func (dr *DepResolver) MustResolveListener() (s listener.Settings, err error) {
	dr.listener.initOnce.Do(func() {
		// resolve and validate all inputs here!
		dr.listener.settings.Port, err = env.GetEnvAsIntOrFallback("K8GB_DISCOVERY_PORT", common.DefaultServicePort)
	})
	return dr.listener.settings, err
}

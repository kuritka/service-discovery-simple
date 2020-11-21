package runner

import "github.com/kuritka/k8gb-discovery/internal/common/guard"

// ServiceRunner is running all commands
type ServiceRunner struct {
	service IServiceRunner
}

// New creates new instance of ServiceRunner
func New(command IServiceRunner) *ServiceRunner {
	return &ServiceRunner{
		command,
	}
}

// MustRun runs service once and panics if service is broken
func (r *ServiceRunner) MustRun() {
	err := r.service.Run()
	guard.FailOnError(err, "command %s failed", r.service.String())
}

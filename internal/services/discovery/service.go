package discovery

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/kuritka/k8gb-discovery/internal/common/log"
	"github.com/kuritka/k8gb-discovery/internal/services/discovery/internal/controller"
)

type Listener struct {
	ctx      context.Context
	settings Settings
}

func New(ctx context.Context, settings Settings) *Listener {
	return &Listener{
		ctx:      ctx,
		settings: settings,
	}
}

func (l *Listener) Run() (err error) {
	runtime.GOMAXPROCS(4)
	c, err := controller.Startup(l.settings.YamlURL)
	if err != nil {
		return
	}
	log.Logger().Infof("listening on :%d", l.settings.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", l.settings.Port), c.Router())

	// TODO: create cache from github
	// TODO: Listener providing DATA. What about controller ?
	// TODO: Tests
	// TODO: channels ?? Look to proxy..
	// TODO: define ingress, service, pod, certmanager with kustomize ??
	// TODO: cache journal by endpoint, isn't enough one endpoint ?
	// Consult with Yury, Tim ...
	return
}

func (l *Listener) String() string {
	return "K8GB discovery listener"
}

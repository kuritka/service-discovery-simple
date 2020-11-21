package listener

import (
	"context"

	"github.com/kuritka/k8gb-discovery/internal/common/log"
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
	log.Logger().Info("hello from listener uhkhghgkj")
	// TODO: create cache from github
	// TODO: Listener providing DATA. What about controller ?
	// TODO: Tests
	// TODO: create repo with one file only (will be internal)
	// TODO: channels ?? Look to proxy..
	// TODO: define ingress, service, pod, certmanager with kustomize ??
	// TODO: heath endpoint
	// TODO: cache jurnal by endpoint
	// Consult with Yury, Tim ...
	return nil
}

func (l *Listener) String() string {
	return "K8GB discovery listener"
}

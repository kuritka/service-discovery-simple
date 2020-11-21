package listener

import (
	"context"
	"fmt"
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
	fmt.Println("hello from listener", l.settings.Port)
	return nil
}

func (l *Listener) String() string {
	return "K8GB discovery listener"
}

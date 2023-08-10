package module

import (
	"context"
	"github.com/1t-data-kit/go-kit/framework/module/signal"
)

type Module interface {
	Name() string
	Type() string
	MustRegisterNetwork() bool
	SignalHandlersMap() signal.HandlersMap
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

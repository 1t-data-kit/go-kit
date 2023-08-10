package service

import (
	"context"
	"github.com/1t-data-kit/go-kit/framework/signal"
)

type Interface interface {
	Name() string
	Type() string
	MustRegisterNetwork() bool
	SignalHandlersMap() signal.HandlersMap
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

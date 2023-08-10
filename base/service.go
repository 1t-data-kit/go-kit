package base

import (
	"context"
)

type Service interface {
	Name() string
	Type() string
	MustRegisterNetwork() bool
	SignalHandlersMap() SignalHandlersMap
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

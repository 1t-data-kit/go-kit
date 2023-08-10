package signal

import (
	"context"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

type Handler func(ctx context.Context) error
type HandlersMap map[os.Signal][]Handler

func NewHandlersMapOption(_map HandlersMap) base.Option {
	return base.NewOption(_map)
}

func (_map HandlersMap) Append(signal os.Signal, handlers ...Handler) {
	if _, exists := _map[signal]; !exists {
		_map[signal] = make([]Handler, 0, len(handlers))
	}
	_map[signal] = append(_map[signal], handlers...)
}

func (_map HandlersMap) Signals() []os.Signal {
	var signals []os.Signal
	for _signal := range _map {
		signals = append(signals, _signal)
	}
	return signals
}

func (_map HandlersMap) Invoke(ctx context.Context, signal os.Signal) error {
	handlers, exists := _map[signal]
	if !exists {
		return nil
	}

	_errors := base.NewErrors()
	for _, handler := range handlers {
		if err := handler(ctx); err != nil {
			_errors.Append(err)
		}
	}

	return _errors.Error()
}

func (_map HandlersMap) Listen(ctx context.Context) {
	if len(_map) == 0 {
		return
	}
	signals := _map.Signals()
	signalChan := make(chan os.Signal, len(signals))
	signal.Notify(signalChan, signals...)
	logrus.Infof("signal listen: %v", signals)
	go func() {
		for {
			select {
			case _signal := <-signalChan:
				logrus.Infof("signal invoke: %v", _signal)
				if err := _map.Invoke(ctx, _signal); err != nil {
					logrus.Errorf("signal invoke error: [%v]%s", _signal, err)
				}
			}
		}
	}()
}

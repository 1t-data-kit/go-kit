package base

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

type SignalHandler func(ctx context.Context) error
type SignalHandlersMap map[os.Signal][]SignalHandler

func (_map SignalHandlersMap) Append(signal os.Signal, handlers ...SignalHandler) {
	if _, exists := _map[signal]; !exists {
		_map[signal] = make([]SignalHandler, 0, len(handlers))
	}
	_map[signal] = append(_map[signal], handlers...)
}

func (_map SignalHandlersMap) Signals() []os.Signal {
	var signals []os.Signal
	for _signal := range _map {
		signals = append(signals, _signal)
	}
	return signals
}

func (_map SignalHandlersMap) Invoke(ctx context.Context, signal os.Signal) error {
	handlers, exists := _map[signal]
	if !exists {
		return nil
	}

	_errors := NewErrors()
	for _, handler := range handlers {
		if err := handler(ctx); err != nil {
			_errors.Append(err)
		}
	}

	return _errors.Error()
}

func (_map SignalHandlersMap) Listen(ctx context.Context) {
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

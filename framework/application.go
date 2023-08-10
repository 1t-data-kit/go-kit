package framework

import (
	"context"
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/1t-data-kit/go-kit/framework/registry/network"
	"github.com/1t-data-kit/go-kit/framework/registry/object"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"sync"
)

var (
	_once     sync.Once
	_instance *application
)

type application struct {
	services          []base.Service
	signalHandlersMap base.SignalHandlersMap
	networkRegistrar  *network.Registrar
	objectRegistrar   *object.Registrar

	running bool
}

func NewApplication(options ...base.Option) *application {
	_once.Do(func() {
		_instance = &application{
			signalHandlersMap: make(base.SignalHandlersMap),
		}
	})

	_instance.init(options...)
	return _instance
}

func (app *application) init(options ...base.Option) {
	_options := base.Options(options)
	if registrars := _options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(*network.Registrar); ok {
			return true
		}
		return false
	}); len(registrars) > 0 {
		app.networkRegistrar = registrars[len(registrars)-1].Value().(*network.Registrar)
	}
	if registrars := _options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(*object.Registrar); ok {
			return true
		}
		return false
	}); len(registrars) > 0 {
		app.objectRegistrar = registrars[len(registrars)-1].Value().(*object.Registrar)
	}
	app.appendServices(_options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(base.Service); ok {
			return true
		}
		return false
	}).Values()...)
	app.appendSignalHandlersMap(_options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(base.SignalHandlersMap); ok {
			return true
		}
		return false
	}).Values()...)
}

func (app *application) appendServices(services ...interface{}) *application {
	for _, service := range services {
		if _service, ok := service.(base.Service); ok {
			app.services = append(app.services, _service)
		}
	}
	return app
}

func (app *application) appendSignalHandlersMap(signalsMaps ...interface{}) *application {
	for _, signalsMap := range signalsMaps {
		if _signalsMap, ok := signalsMap.(base.SignalHandlersMap); ok {
			for _signal, handlers := range _signalsMap {
				app.signalHandlersMap.Append(_signal, handlers...)
			}
		}
	}
	return app
}

func (app *application) ObjectRegistrar() *object.Registrar {
	return app.objectRegistrar
}

func (app *application) NetworkRegistrar() *network.Registrar {
	return app.networkRegistrar
}

func (app *application) Start(ctx context.Context, options ...base.Option) error {
	if app.running {
		return fmt.Errorf("application has be running")
	}

	app.init(options...)
	if len(app.services) == 0 {
		return fmt.Errorf("application has no service to run")
	}

	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)
	for _, service := range app.services {
		_service := service
		group.Go(func() error {
			logrus.Infof("application start service: %s[%s]", _service.Name(), _service.Type())
			if err := _service.Start(ctx); err != nil {
				return err
			}
			logrus.Infof("application stop service: %s[%s]", _service.Name(), _service.Type())
			return nil
		})

		app.appendSignalHandlersMap(_service.SignalHandlersMap())
		if app.networkRegistrar != nil && _service.MustRegisterNetwork() {
			app.networkRegistrar.Register(ctx, _service)
			logrus.Infof("application.NetworkRegistrar register service %s[%s]", _service.Name(), _service.Type())
		}
	}
	app.signalHandlersMap.Listen(ctx)
	app.running = true

	return group.Wait()
}

func (app *application) Stop(ctx context.Context) error {
	if !app.running {
		return fmt.Errorf("application has not be running")
	}

	_errors := base.NewErrors()
	for _, service := range app.services {
		if err := service.Stop(ctx); err != nil {
			_errors.Append(err)
		}
	}
	app.running = false

	err := _errors.Error()
	if err == nil {
		return nil
	}
	return errors.Wrap(err, "application stop error")
}

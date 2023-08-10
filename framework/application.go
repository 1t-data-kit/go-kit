package framework

import (
	"context"
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/1t-data-kit/go-kit/framework/command"
	"github.com/1t-data-kit/go-kit/framework/registry/network"
	"github.com/1t-data-kit/go-kit/framework/registry/object"
	"github.com/1t-data-kit/go-kit/framework/service"
	"github.com/1t-data-kit/go-kit/framework/signal"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"sync"
)

var (
	_once     sync.Once
	_instance *application
)

type application struct {
	services          []service.Interface
	signalHandlersMap signal.HandlersMap
	networkRegistrar  *network.Registrar
	objectRegistrar   *object.Registrar
	commandAgent      *command.Agent

	running bool
}

func NewApplication(options ...base.Option) *application {
	_once.Do(func() {
		_instance = &application{
			signalHandlersMap: make(signal.HandlersMap),
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

	if agents := _options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(*command.Agent); ok {
			return true
		}
		return false
	}); len(agents) > 0 {
		app.commandAgent = agents[len(agents)-1].Value().(*command.Agent)
	}

	app.appendServices(_options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(service.Interface); ok {
			return true
		}
		return false
	}).Values()...)

	app.appendSignalHandlersMap(_options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(signal.HandlersMap); ok {
			return true
		}
		return false
	}).Values()...)
}

func (app *application) appendServices(services ...interface{}) *application {
	for _, _service := range services {
		if __service, ok := _service.(service.Interface); ok {
			app.services = append(app.services, __service)
		}
	}
	return app
}

func (app *application) appendSignalHandlersMap(signalsMaps ...interface{}) *application {
	for _, signalsMap := range signalsMaps {
		if _signalsMap, ok := signalsMap.(signal.HandlersMap); ok {
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

	if app.commandAgent != nil && app.commandAgent.MustRun() {
		return app.commandAgent.Run(ctx, os.Args)
	}

	if len(app.services) == 0 {
		return fmt.Errorf("application has no service to run")
	}

	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)
	for _, _service := range app.services {
		__service := _service
		group.Go(func() error {
			logrus.Infof("application start service: %s[%s]", __service.Name(), __service.Type())
			if err := __service.Start(ctx); err != nil {
				return err
			}
			logrus.Infof("application stop service: %s[%s]", __service.Name(), __service.Type())
			return nil
		})

		app.appendSignalHandlersMap(__service.SignalHandlersMap())
		if app.networkRegistrar != nil && __service.MustRegisterNetwork() {
			if endpoint, ok := __service.(network.Endpoint); ok {
				app.networkRegistrar.Register(ctx, endpoint)
				logrus.Infof("application.NetworkRegistrar register service %s[%s]", __service.Name(), __service.Type())
			}
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

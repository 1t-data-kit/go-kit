package framework

import (
	"context"
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/1t-data-kit/go-kit/framework/module/registry/network"
	"github.com/1t-data-kit/go-kit/framework/module/registry/object"
	signalLib "github.com/1t-data-kit/go-kit/framework/module/signal"
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
	modules           []base.Module
	signalHandlersMap signalLib.HandlersMap
	networkRegistrar  *network.Registrar
	objectRegistrar   *object.Registrar

	running bool
}

func NewApplication(options ...base.Option) *application {
	_once.Do(func() {
		_instance = &application{
			signalHandlersMap: make(signalLib.HandlersMap),
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
	app.appendModels(_options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(base.Module); ok {
			return true
		}
		return false
	}).Values()...)
	app.appendSignalHandlersMap(_options.Filter(func(item base.Option) bool {
		if _, ok := item.Value().(signalLib.HandlersMap); ok {
			return true
		}
		return false
	}).Values()...)
}

func (app *application) appendModels(modules ...interface{}) *application {
	for _, _module := range modules {
		if __module, ok := _module.(base.Module); ok {
			app.modules = append(app.modules, __module)
		}
	}
	return app
}

func (app *application) appendSignalHandlersMap(signalsMaps ...interface{}) *application {
	for _, signalsMap := range signalsMaps {
		if _signalsMap, ok := signalsMap.(signalLib.HandlersMap); ok {
			for _signal, handlers := range _signalsMap {
				if _, exists := app.signalHandlersMap[_signal]; !exists {
					app.signalHandlersMap[_signal] = make([]signalLib.Handler, 0)
				}
				app.signalHandlersMap[_signal] = append(app.signalHandlersMap[_signal], handlers...)
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
	if len(app.modules) == 0 {
		return fmt.Errorf("application has no module to run")
	}

	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)
	for _, module := range app.modules {
		_module := module
		group.Go(func() error {
			logrus.Infof("application start module: %s[%s]", _module.Name(), _module.Type())
			if err := _module.Start(ctx); err != nil {
				return err
			}
			logrus.Infof("application stop module %s[%s]", _module.Name(), _module.Type())
			return nil
		})

		app.appendSignalHandlersMap(_module.SignalHandlersMap())
		if app.networkRegistrar != nil && _module.MustRegisterNetwork() {
			app.networkRegistrar.Register(ctx, _module)
			logrus.Infof("application.NetworkRegistrar register module %s[%s]", _module.Name(), _module.Type())
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
	for _, module := range app.modules {
		if err := module.Stop(ctx); err != nil {
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

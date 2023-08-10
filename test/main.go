package main

import (
	"context"
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/1t-data-kit/go-kit/framework"
	"github.com/1t-data-kit/go-kit/framework/signal"
	"github.com/sirupsen/logrus"
	"syscall"
	"time"
)

type a struct {
	name string
}

func newA(name string) *a {
	return &a{
		name: name,
	}
}

func (_a *a) Name() string {
	return _a.name
}

func (_a *a) Type() string {
	return "a module"
}

func (_a *a) Start(ctx context.Context) error {
	time.Sleep(time.Minute)
	return nil
}

func (_a *a) Stop(ctx context.Context) error {
	return nil
}

func (_a *a) MustRegisterNetwork() bool {
	return true
}

func (_a *a) SignalHandlersMap() signal.HandlersMap {
	return signal.HandlersMap{
		syscall.SIGTERM: []signal.Handler{
			func(ctx context.Context) error {
				fmt.Printf("%s[%s] sigterm\n", _a.Name(), _a.Type())
				return nil
			},
		},
		syscall.SIGQUIT: []signal.Handler{
			func(ctx context.Context) error {
				fmt.Printf("%s[%s] siqquit\n", _a.Name(), _a.Type())
				return nil
			},
		},
	}
}

func main() {
	app := framework.NewApplication(base.NewOption(newA("a1")), base.NewOption(newA("a2")), base.NewOption(newA("a3")), base.NewOption(signal.HandlersMap{
		syscall.SIGTERM: []signal.Handler{
			func(ctx context.Context) error {
				logrus.Info("application outside sigterm invoke")
				return nil
			},
		},
	}))
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	app.Stop(context.Background())
}

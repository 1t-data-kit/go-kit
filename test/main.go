package main

import (
	"context"
	"fmt"
	"github.com/1t-data-kit/go-kit/base"
	"github.com/1t-data-kit/go-kit/framework"
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

func (_a *a) GetName() string {
	return _a.name
}

func (_a *a) GetType() string {
	return "a module"
}

func (_a *a) Start(ctx context.Context) error {
	time.Sleep(time.Minute)
	return nil
}

func (_a *a) Stop(ctx context.Context) error {
	return nil
}

func (_a *a) SignalHandlersMap() base.HandlersMap {
	return base.HandlersMap{
		syscall.SIGTERM: []base.Handler{
			func(ctx context.Context) error {
				fmt.Printf("%s[%s] sigterm\n", _a.GetName(), _a.GetName())
				return nil
			},
		},
		syscall.SIGQUIT: []base.Handler{
			func(ctx context.Context) error {
				fmt.Printf("%s[%s] siqquit\n", _a.GetName(), _a.GetName())
				return nil
			},
		},
	}
}

func main() {
	app := framework.Application(framework.ModuleOption(newA("a1")), framework.ModuleOption(newA("a2")), framework.ModuleOption(newA("a3")), framework.SignalOption(syscall.SIGTERM, func(ctx context.Context) error {
		logrus.Info("application outside sigterm invoke")
		return nil
	}))
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	app.Stop(context.Background())
}

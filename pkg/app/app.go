package app

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

func New(name, version string, servers []Server, opts ...Option) *App {
	a := &App{
		ctx:     context.Background(),
		name:    name,
		version: version,
		servers: servers,
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}

	if id, err := uuid.NewUUID(); err == nil {
		a.id = id.String()
	}

	for _, o := range opts {
		o(a)
	}

	a.ctx, a.cancel = context.WithCancel(a.ctx)
	return a
}

type App struct {
	id      string
	ctx     context.Context
	name    string
	version string
	servers []Server
	sigs    []os.Signal
	cancel  func()
}

func (a *App) Run() error {
	eg, ctx := errgroup.WithContext(a.ctx)

	for _, srv := range a.servers {
		srv := srv
		// 注册关闭
		eg.Go(func() error {
			<-ctx.Done()
			return srv.Stop(ctx)
		})

		// 开启服务
		eg.Go(func() error {
			return srv.Start(ctx)
		})
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				err := a.Stop()
				if err != nil {
					return err
				}
			}
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func (a *App) ID() string {
	return a.id
}

func (a *App) Name() string {
	return a.name
}

func (a *App) Version() string {
	return a.version
}

type Option func(a *App)

func SetID(id string) Option {
	return func(a *App) {
		a.id = id
	}
}

func WithContext(ctx context.Context) Option {
	return func(a *App) {
		a.ctx = ctx
	}
}

func WithSigs(sigs []os.Signal) Option {
	return func(a *App) {
		a.sigs = sigs
	}
}

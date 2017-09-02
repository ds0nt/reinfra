package components

import (
	"context"
	"fmt"

	"github.com/ds0nt/reinfra/config"
	"github.com/ds0nt/reinfra/readymanager"
	"github.com/ds0nt/reinfra/service"
	"github.com/jackc/pgx"
)

type Postgreser interface {
	Pg() *pgx.ConnPool
}

type Postgres struct {
	Pool         *pgx.ConnPool
	config       pgx.ConnPoolConfig
	customConfig bool
	afterConnect func(*pgx.Conn) error
	readymanager.ReadyManager
}

func (a *Postgres) Run(s *service.Service) (err error) {
	if !a.customConfig {
		a.config = config.EnvPostgresConfig()
	}

	if a.afterConnect != nil {
		a.config.AfterConnect = func(c *pgx.Conn) error {
			err := a.afterConnect(c)
			if err != nil {
				s.Log().Error(err)
				return err
			}
			a.ReadyManager.SetReady()
			return nil
		}
	}

	fmt.Println("Creating Postgres Pool", a.config)

	a.Pool, err = pgx.NewConnPool(a.config)
	if a.afterConnect == nil {
		a.ReadyManager.SetReady()
	}
	return
}

func (a *Postgres) SetConfig(c pgx.ConnPoolConfig) {
	a.config = c
	a.customConfig = true
}

func (a *Postgres) SetAfterConnect(fn func(*pgx.Conn) error) {
	a.afterConnect = fn
}

func (a *Postgres) Stop(*service.Service) error {
	a.Pool.Close()
	return nil
}

func (a *Postgres) Pg() *pgx.ConnPool {
	return a.Pool
}

func (s *Postgres) Ready() bool {
	return s.ReadyManager.Ready()
}

func (s *Postgres) WaitForReady(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-s.ReadyManager.ReadyCh():
		return
	}
}

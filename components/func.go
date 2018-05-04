package components

import (
	"context"

	"github.com/ds0nt/reinfra/readymanager"
	"github.com/ds0nt/reinfra/service"
)

type Func struct {
	readymanager.ReadyManager
	runFn  func(s *service.Service) error
	stopFn func(s *service.Service) error
}

func (a *Func) Run(s *service.Service) (err error) {
	err = a.runFn(s)
	if err != nil {
		return err
	}
	a.ReadyManager.SetReady()
	return nil
}

func (a *Func) Stop(s *service.Service) error {
	return a.stopFn(s)
}

func (s *Func) Ready() bool {
	return s.ReadyManager.Ready()
}

func (s *Func) WaitForReady(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-s.ReadyManager.ReadyCh():
		return
	}
}

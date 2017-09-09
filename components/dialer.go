package components

import (
	"context"

	"github.com/ds0nt/reinfra/readymanager"
	"github.com/ds0nt/reinfra/service"
)

type DialerComponent struct {
	dialer GRPCDialer
	readymanager.ReadyManager
}

func (s *DialerComponent) Run(*service.Service) error {
	defer s.ReadyManager.SetReady()
	return s.dialer.DialReinfra()
}

func (s *DialerComponent) Stop(*service.Service) error {
	return s.dialer.Close()
}

func (s *DialerComponent) Ready() bool {
	return s.ReadyManager.Ready()
}

func (s *DialerComponent) WaitForReady(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-s.ReadyManager.ReadyCh():
		return
	}
}

// WrapDialer turns a dialer into a service component :)
func WrapDialer(d GRPCDialer) service.ServiceComponent {
	return &DialerComponent{
		dialer: d,
	}
}

// GRPCDialer is basically just a grpc client that is compatible via Init and Run
type GRPCDialer interface {
	DialReinfra() (err error)
	Close() error
}

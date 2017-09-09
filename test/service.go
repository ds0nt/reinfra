package test

import (
	"golang.org/x/net/context"

	"github.com/ds0nt/reinfra/components"
	"github.com/ds0nt/reinfra/service"
)

type TestService struct {
	*service.Service
	*components.GRPCServer
}

func (s *TestService) Create(ctx context.Context, in *TestMessage) (*TestMessage, error) {
	return in, nil
}

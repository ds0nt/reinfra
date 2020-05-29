package components

import (
	"context"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/sirupsen/logrus"

	"github.com/ds0nt/reinfra/config"
	logmw "github.com/ds0nt/reinfra/logger-middleware"
	"github.com/ds0nt/reinfra/readymanager"
	"github.com/ds0nt/reinfra/service"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	Addr   string
	readymanager.ReadyManager
	log         *logrus.Entry
	GRPCOptions []grpc.ServerOption
}

func (s *GRPCServer) Init(svc *service.Service) {
	s.log = svc.Log().WithField("component", s)
}

func (s *GRPCServer) Server() *grpc.Server {
	if s.server != nil {
		return s.server
	}

	// serve
	s.server = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			logmw.UnaryServerInterceptor(s.log),
		),
	)

	return s.server
}

func (s *GRPCServer) Run(svc *service.Service) error {
	if len(s.Addr) == 0 {
		s.Addr = config.GRPCAddr
	}

	s.log.Println("listening on", s.Addr)
	defer s.log.Println("stopped")

	// create tcp C
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return errors.Wrap(err, "tcp for grpc failed to listen")
	}

	// handle readiness
	go func() {
		for {
			select {
			case <-svc.Context().Done():
				return
			default:
				c, err := grpc.Dial(s.Addr, grpc.WithInsecure())
				if err == nil {
					c.Close()
					s.ReadyManager.SetReady()
					return
				}
				fmt.Println(err)
			}
		}
	}()

	return s.Server().Serve(lis)
}

func (s *GRPCServer) Stop(*service.Service) error {
	if s.server != nil {
		s.server.Stop()
	}
	return nil
}

func (s *GRPCServer) Ready() bool {
	return s.ReadyManager.Ready()
}

func (s *GRPCServer) WaitForReady(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-s.ReadyManager.ReadyCh():
		return
	}
}

func (s *GRPCServer) String() string {
	return "grpc-" + s.Addr
}

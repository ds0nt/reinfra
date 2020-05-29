package components

import (
	"context"
	"net/http"
	"time"

	"github.com/ds0nt/reinfra/config"
	"github.com/ds0nt/reinfra/readymanager"
	"github.com/ds0nt/reinfra/service"
)

type HTTPServer struct {
	Addr        string
	server      *http.Server
	HTTPHandler http.Handler
	readymanager.ReadyManager
}

func (s *HTTPServer) Run(svc *service.Service) error {
	log := svc.Log().WithField("component", s)
	if len(s.Addr) == 0 {
		s.Addr = config.HTTPAddr
	}
	log.Println("listening on", s.Addr)
	defer log.Println("stopped")
	s.server = &http.Server{Addr: s.Addr}
	s.server.Handler = s.HTTPHandler

	go func() {
		time.Sleep(time.Millisecond * 10)
		s.ReadyManager.SetReady()
	}()
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(*service.Service) error {
	return s.server.Close()
}

func (s *HTTPServer) Ready() bool {
	return s.ReadyManager.Ready()
}

func (s *HTTPServer) WaitForReady(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-s.ReadyManager.ReadyCh():
		return
	}
}

func (s *HTTPServer) String() string {
	return "http-" + s.Addr
}

package components

import (
	"context"
	"net/http"
	"time"

	"github.com/ds0nt/reinfra/config"
	"github.com/ds0nt/reinfra/readymanager"
	"github.com/ds0nt/reinfra/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	Addr   string
	server *http.Server
	readymanager.ReadyManager
}

func (s *MetricsServer) Run(svc *service.Service) error {
	log := svc.Log().WithField("component", s)
	if len(s.Addr) == 0 {
		s.Addr = config.MetricsAddr
	}
	log.Println("listening on", s.Addr)
	defer log.Println("stopped")
	s.server = &http.Server{Addr: s.Addr}
	s.server.Handler = promhttp.Handler()

	go func() {
		time.Sleep(time.Millisecond * 10)
		s.ReadyManager.SetReady()
	}()
	return s.server.ListenAndServe()
}

func (s *MetricsServer) Stop(*service.Service) error {
	return s.server.Close()
}

func (s *MetricsServer) Ready() bool {
	return s.ReadyManager.Ready()
}

func (s *MetricsServer) WaitForReady(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-s.ReadyManager.ReadyCh():
		return
	}
}

func (s *MetricsServer) String() string {
	return "metrics-" + s.Addr
}

package components

import (
	"context"
	"net/http"
	"time"

	"github.com/ds0nt/reinfra/config"
	"github.com/ds0nt/reinfra/readymanager"
	"github.com/ds0nt/reinfra/service"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsServer struct {
	Addr   string
	server *http.Server
	mux    *http.ServeMux
	readymanager.ReadyManager
}

func (s *MetricsServer) Run(svc *service.Service) error {
	log := svc.Log().WithField("component", s)
	if len(s.Addr) == 0 {
		s.Addr = config.MetricsAddr
	}
	log.Println("listening on", s.Addr)
	s.server = &http.Server{Addr: s.Addr}
	s.mux = http.NewServeMux()
	s.mux.Handle("/metrics", prometheus.Handler())
	s.server.Handler = s.mux

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

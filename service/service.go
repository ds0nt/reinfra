package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ds0nt/reinfra/readymanager"
)

/*
	Creating a Service
	1. Create a Service
	2. Add any components, readyable components


	Testing
	1. Instantiate Service
	2. go infra.Run()
	3. WaitForReady
		- service is ready if all components are ready
	4. Run tests
	5. infra.Stop()

	Running
	1. infra.Run the service
	  - Runs all Components
	  - use service IsReady to power readiness probes
	  - blocks until stopped by an error
	  - tears down service if error occurs
*/

type Service struct {
	parentContext context.Context
	components    []ServiceComponent
	readymanager.ReadyManager
	logger *logrus.Logger
	fields logrus.Fields
	entry  *logrus.Entry
}

func (s *Service) Log() *logrus.Entry {
	return s.entry.WithFields(s.fields)
}
func (s *Service) SetLogFields(fields logrus.Fields) {
	s.fields = fields
}

// Init initializes a service. It is run before components init, and before any Run methods are called.
func (s *Service) Init() {
	s.logger = logrus.New()
	s.entry = logrus.NewEntry(s.logger)
	s.logger.Info("Service Initialized")
}

// Run starts all of a services components. It returns if a service errors, or the context is cancelled
func (s *Service) Run(ctx context.Context) chan error {
	s.parentContext = ctx
	stopCh := make(chan error) //, len(s.components))
	errors := make(chan error) //, len(s.components)+1)

	for _, c := range s.components {
		go func(c ServiceComponent) {
			err := c.Run(s)
			if err != nil {
				stopCh <- err
			}
		}(c)

		// silly way to wait for above goroutine to start blocking
		time.Sleep(time.Millisecond * 10)
	}

	s.ReadyManager.SetReady()

	go func() {
		select {
		case x := <-stopCh:
			errors <- x
		case <-ctx.Done():
			errors <- ctx.Err()
		}
		s.ReadyManager.SetUnready()

		for _, c := range s.components {
			err := c.Stop(s)
			if err != nil {
				errors <- err
			}
		}
		// close(stopCh)
		close(errors)
	}()

	return errors
}

// Context returns the services context. All components should stop when this context
// is done.
func (s *Service) Context() context.Context {
	return s.parentContext
}

// WaitForReady waits for the service, and all of it's readyable components to become ready
func (s *Service) WaitForReady(ctx context.Context) {
	<-s.ReadyManager.ReadyCh()

	for _, c := range s.components {
		c.WaitForReady(ctx)
	}
}

// Ready determines if a service and all of a services components are ready
func (s *Service) Ready() bool {
	if !s.ReadyManager.Ready() {
		return false
	}
	for _, c := range s.components {
		if !c.Ready() {
			return false
		}
	}
	return true
}

func (s *Service) RegisterComponent(component ServiceComponent) {
	s.components = append(s.components, component)
}

type ServiceComponent interface {
	Run(*Service) error
	Stop(*Service) error
	Ready() bool
	WaitForReady(context.Context)
}

// Initer powers the Init method, which will run before the Run() if it defined
// in the method set of a service component
type Initer interface {
	Init(*Service)
}

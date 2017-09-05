package reinfra_test

import (
	"context"
	"fmt"

	"github.com/ds0nt/reinfra"
	"github.com/ds0nt/reinfra/components"
	"github.com/ds0nt/reinfra/readymanager"
	"github.com/ds0nt/reinfra/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("readymanager", func() {
	It("readymanager should work", func() {
		r := readymanager.ReadyManager{}
		Expect(r.Ready()).To(BeFalse())
		Expect(r.Ready()).To(BeFalse())

		go func() {
			r.SetReady()
			Expect(r.Ready()).To(BeTrue())
		}()

		<-r.ReadyCh()
		Expect(r.Ready()).To(BeTrue())
		Expect(r.Ready()).To(BeTrue())

		r.SetUnready()
		Expect(r.Ready()).To(BeFalse())
		Expect(r.Ready()).To(BeFalse())

		go func() {
			r.SetReady()
			Expect(r.Ready()).To(BeTrue())
		}()

		<-r.ReadyCh()
		Expect(r.Ready()).To(BeTrue())
		Expect(r.Ready()).To(BeTrue())
	})
})

var _ = Describe("Infra", func() {
	It("should find and initialize it's service", func(done Done) {
		type svc struct {
			*service.Service
			*components.GRPCServer
		}
		s := svc{}

		Expect(s.Service).To(BeNil())

		reinfra.Init(&s)

		Expect(s.Service).NotTo(BeNil())
		Expect(s.Service.Ready()).To(BeFalse())

		ctx, cancel := context.WithCancel(context.Background())

		errCh := reinfra.Run(ctx, &s)

		s.Service.WaitForReady(ctx)

		Expect(s.Service.Ready()).To(BeTrue())

		cancel()
		for err := range errCh {
			fmt.Println("Service Err Ch", err)
		}
		Expect(s.Service.Ready()).To(BeFalse())
		close(done)
	}, 0.5)
})

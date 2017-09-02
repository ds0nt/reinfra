package reinfra_test

import (
	"context"

	"github.com/ds0nt/reinfra"
	"github.com/ds0nt/reinfra/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Infra", func() {
	It("should find and initialize it's service", func() {

		type svc struct {
			*service.Service
		}
		s := svc{}

		Expect(s.Service).To(BeNil())

		reinfra.Init(&s)

		Expect(s.Service).NotTo(BeNil())
		Expect(s.Service.Ready()).To(BeFalse())

		reinfra.Run(context.Background(), &s)
		Expect(s.Service.Ready()).To(BeFalse())
		s.Service.WaitForReady(context.Background())
		Expect(s.Service.Ready()).To(BeTrue())

	})
})

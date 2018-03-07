package logger

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("logger", func() {

	Describe("populateRequestContext", func() {
		It("returns a context including a RequestContext based on the given request", func() {
			const (
				remoteAddr = "127.0.0.1"
				reqID      = "123"
			)
			r, err := http.NewRequest("GET", "https://homes.co.nz", http.NoBody)
			Expect(err).ToNot(HaveOccurred())
			r.RemoteAddr = remoteAddr
			r.Header.Set("X-Request-Id", reqID)
			ctx := populateRequestContext(context.Background(), r)

			rctx, ok := retrieveRequestContext(ctx)
			Expect(ok).To(BeTrue())
			Expect(rctx.RemoteAddr).To(Equal(remoteAddr))
			Expect(rctx.RequestID).To(Equal(reqID))
		})
		Context("using X-Forwarded-For instead of remote addr", func() {
			It("returns a context including a RequestContext based on the given request", func() {
				const (
					remoteAddr = "127.0.0.1"
					reqID      = "123"
				)
				r, err := http.NewRequest("GET", "https://homes.co.nz", http.NoBody)
				Expect(err).ToNot(HaveOccurred())
				r.Header.Set("X-Forwarded-For", remoteAddr)
				r.Header.Set("X-Request-Id", reqID)
				ctx := populateRequestContext(context.Background(), r)

				rctx, ok := retrieveRequestContext(ctx)
				Expect(ok).To(BeTrue())
				Expect(rctx.RemoteAddr).To(Equal(remoteAddr))
				Expect(rctx.RequestID).To(Equal(reqID))
			})
		})
	})

	Describe("retrieveRequestContext", func() {
		It("returns an empty (but valid) context if not initialized", func() {
			rctx, ok := retrieveRequestContext(context.Background())
			Expect(ok).To(BeFalse())
			Expect(rctx).To(Equal(requestContext{}))
		})
	})

})

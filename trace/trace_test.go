package trace

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTrace(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trace Suite")
}

var _ = Describe("Trace", func() {
	Context(".FromCtx", func() {
		It("should return empty Trace", func() {
			ctx := context.Background()
			trace := FromCtx(ctx)

			Expect(trace.EventID).To(Equal(""))
			Expect(trace.CorrelationID).To(Equal(""))
			Expect(trace.CausationID).To(Equal(""))
		})
		It("should return Trace", func() {
			ctx := context.Background()
			trace := New()
			ctx = SetToCtx(ctx, trace)
			v, ok := ctx.Value(key).(Trace)
			Expect(ok).To(Equal(true))
			Expect(v).To(Equal(trace))
		})
	})
	Context(".SetToCtx", func() {
		It("should return a new valid Trace", func() {
			ctx := context.Background()
			ctx = SetToCtx(ctx, Trace{})
			trace, ok := ctx.Value(key).(Trace)
			Expect(ok).To(Equal(true))
			Expect(trace.EventID).NotTo(Equal(""))
			Expect(trace.CorrelationID).NotTo(Equal(""))
			Expect(trace.CausationID).NotTo(Equal(""))
		})
		It("should return Trace", func() {
			ctx := context.Background()
			trace := New()
			ctx = SetToCtx(ctx, trace)
			v, ok := ctx.Value(key).(Trace)
			Expect(ok).To(Equal(true))
			Expect(v).To(Equal(trace))
		})
	})
	Context(".SetToCtxFromHeader", func() {
		It("should return a new valid Trace", func() {
			ctx := context.Background()
			ctx = SetToCtxFromHeader(ctx, "")
			trace := FromCtx(ctx)
			Expect(trace.EventID).NotTo(Equal(""))
			Expect(trace.CorrelationID).NotTo(Equal(""))
			Expect(trace.CausationID).NotTo(Equal(""))
		})
		It("should return Trace", func() {
			ctx := context.Background()
			trace := New()
			ctx = SetToCtxFromHeader(ctx, trace.ToJSON())
			ctxTrace := FromCtx(ctx)
			Expect(ctxTrace.EventID).To(Equal(trace.EventID))
			Expect(ctxTrace.CorrelationID).To(Equal(trace.CorrelationID))
			Expect(ctxTrace.CausationID).To(Equal(trace.CausationID))
		})
	})
	Context(".ToJSON", func() {
		It("should return a valid JSON", func() {
			trace := New()
			json := trace.ToJSON()
			Expect(json).NotTo(Equal(""))
		})
		It("should return an empty string", func() {
			trace := Trace{}
			json := trace.ToJSON()
			Expect(json).To(Equal(""))
		})
	})
	Context(".FromJSON", func() {
		It("should return a valid Trace", func() {
			trace := New()
			json := trace.ToJSON()
			traceFromJSON := FromJSON(json)
			Expect(traceFromJSON.EventID).To(Equal(trace.EventID))
			Expect(traceFromJSON.CorrelationID).To(Equal(trace.CorrelationID))
			Expect(traceFromJSON.CausationID).To(Equal(trace.CausationID))
		})
		It("should return empty Trace", func() {
			trace := FromJSON("")
			Expect(trace.EventID).To(Equal(""))
			Expect(trace.CorrelationID).To(Equal(""))
			Expect(trace.CausationID).To(Equal(""))
		})
	})
	Context(".ToJSONFromCtx", func() {
		It("should return a valid JSON", func() {
			ctx := context.Background()
			trace := New()
			ctx = SetToCtx(ctx, trace)
			json := ToJSONFromCtx(ctx)
			Expect(json).NotTo(Equal(""))
		})
		It("should return an empty string", func() {
			ctx := context.Background()
			json := ToJSONFromCtx(ctx)
			Expect(json).To(Equal(""))
		})
	})
	Context(".LinkCtxFromJSON", func() {
		It("should return a valid Trace", func() {
			ctx := context.Background()
			trace := New()
			ctx = SetToCtx(ctx, trace)
			json := trace.ToJSON()
			newCtx := LinkCtxFromJSON(ctx, json)
			newTrace, ok := newCtx.Value(key).(Trace)
			Expect(ok).To(Equal(true))

			Expect(newTrace.EventID).NotTo(Equal(trace.EventID))
			Expect(newTrace.CorrelationID).To(Equal(trace.EventID))
			Expect(newTrace.CausationID).To(Equal(trace.EventID))
		})
		It("should return a valid Trace with different causation id", func() {
			ctx := context.Background()
			trace := New()
			json := trace.ToJSON()
			newCtx := LinkCtxFromJSON(ctx, json)
			newTrace := FromCtx(newCtx)
			Expect(newTrace.EventID).NotTo(Equal(trace.EventID))
			Expect(newTrace.CorrelationID).To(Equal(trace.EventID))
			Expect(newTrace.CausationID).To(Equal(trace.EventID))

			newTraceJSON := newTrace.ToJSON()
			newCtx = LinkCtxFromJSON(ctx, newTraceJSON)
			newTrace = FromCtx(newCtx)
			Expect(newTrace.EventID).NotTo(Equal(trace.EventID))
			Expect(newTrace.CorrelationID).To(Equal(trace.CorrelationID))
			Expect(newTrace.CausationID).NotTo(Equal(trace.EventID))
		})
		It("should return a new valid Trace, if header is not valid Trace", func() {
			ctx := context.Background()
			newCtx := LinkCtxFromJSON(ctx, "")
			newTrace := FromCtx(newCtx)
			Expect(newTrace.EventID).NotTo(Equal(""))
			Expect(newTrace.CorrelationID).NotTo(Equal(""))
			Expect(newTrace.CausationID).NotTo(Equal(""))
		})
		It("should return a new valid Trace, if header is not valid Trace", func() {
			ctx := context.Background()
			newCtx := LinkCtxFromJSON(ctx, "{}")
			newTrace := FromCtx(newCtx)
			Expect(newTrace.EventID).NotTo(Equal(""))
			Expect(newTrace.CorrelationID).NotTo(Equal(""))
			Expect(newTrace.CausationID).NotTo(Equal(""))
		})
	})
	Context(".LinkCtxFromCtx", func() {
		It("should return a valid Trace", func() {
			ctx := context.Background()
			trace := New()
			ctx = SetToCtx(ctx, trace)
			newCtx := LinkCtxFromCtx(ctx)
			newTrace, ok := newCtx.Value(key).(Trace)
			Expect(ok).To(Equal(true))
			Expect(newTrace.EventID).NotTo(Equal(trace.EventID))
			Expect(newTrace.CorrelationID).To(Equal(trace.EventID))
			Expect(newTrace.CausationID).To(Equal(trace.EventID))
		})
		It("should return a valid Trace with different causation id", func() {
			ctx := context.Background()
			trace := New()
			ctx = SetToCtx(ctx, trace)
			newCtx := LinkCtxFromCtx(ctx)
			newTrace := FromCtx(newCtx)
			Expect(newTrace.EventID).NotTo(Equal(trace.EventID))
			Expect(newTrace.CorrelationID).To(Equal(trace.EventID))
			Expect(newTrace.CausationID).To(Equal(trace.EventID))

			newCtx = LinkCtxFromCtx(newCtx)
			newTrace = FromCtx(newCtx)
			Expect(newTrace.EventID).NotTo(Equal(trace.EventID))
			Expect(newTrace.CorrelationID).To(Equal(trace.CorrelationID))
			Expect(newTrace.CausationID).NotTo(Equal(trace.EventID))
		})
		It("should return a new valid Trace, if trace is not presented in ctx", func() {
			ctx := context.Background()
			newCtx := LinkCtxFromCtx(ctx)
			newTrace := FromCtx(newCtx)
			Expect(newTrace.EventID).NotTo(Equal(""))
			Expect(newTrace.CorrelationID).NotTo(Equal(""))
			Expect(newTrace.CausationID).NotTo(Equal(""))
		})
	})
})

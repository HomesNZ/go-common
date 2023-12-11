package gin

import (
	"context"
	"github.com/HomesNZ/go-common/trace"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrace(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trace Suite")
}

type log struct{}

func (l *log) Info(ctx context.Context, msg string, args ...any)  {}
func (l *log) Error(ctx context.Context, msg string, args ...any) {}
func (l *log) Warn(ctx context.Context, msg string, args ...any)  {}
func (l *log) Debug(ctx context.Context, msg string, args ...any) {}

var testTrace = trace.New()
var traceJSON = trace.ToJSON(testTrace)

var traceWithDiffCausation = &trace.Trace{
	EventID:       uuid.NewString(),
	CorrelationID: uuid.NewString(),
	CausationID:   uuid.NewString(),
}
var traceJSONWithDiffCausation = trace.ToJSON(traceWithDiffCausation)

func newServerWithTraceHeader() *gin.Engine {
	router := gin.Default()
	l := &log{}
	router.Use(LoggerMiddleware(l, "/health", "/metrics"))
	router.GET("/", func(c *gin.Context) {
		c.String(200, "response")
	})
	return router
}

var _ = Describe("LoggerMiddleware", func() {
	Describe("Trace", func() {
		Context("without trace header", func() {
			It("should set trace header with Trace", func() {
				req, _ := http.NewRequest("GET", "/", nil)
				w := httptest.NewRecorder()
				r := newServerWithTraceHeader()
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(200))
				traceHeader := w.Header().Get("X-Homes-Trace")
				Expect(traceHeader).NotTo(Equal(""))
				newTrace := trace.FromJSON(traceHeader)
				Expect(newTrace).NotTo(BeNil())
				Expect(newTrace.EventID).NotTo(Equal(""))
				Expect(newTrace.EventID).To(Equal(newTrace.CorrelationID))
				Expect(newTrace.EventID).To(Equal(newTrace.CausationID))
				Expect(newTrace.CorrelationID).To(Equal(newTrace.CausationID))

				body, _ := ioutil.ReadAll(w.Body)
				Expect(string(body)).To(Equal("response"))
			})
		})
		Context("with trace header", func() {
			It("should set a new event id to the trace header", func() {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("X-Homes-Trace", traceJSON)
				w := httptest.NewRecorder()
				r := newServerWithTraceHeader()
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(200))
				traceHeader := w.Header().Get("X-Homes-Trace")
				Expect(traceHeader).NotTo(Equal(""))
				newTrace := trace.FromJSON(traceHeader)
				Expect(newTrace).NotTo(BeNil())
				Expect(newTrace.EventID).NotTo(Equal(""))
				Expect(newTrace.EventID).NotTo(Equal(newTrace.CorrelationID))
				Expect(newTrace.EventID).NotTo(Equal(newTrace.CausationID))
				Expect(newTrace.CorrelationID).To(Equal(newTrace.CausationID))

				body, _ := ioutil.ReadAll(w.Body)
				Expect(string(body)).To(Equal("response"))
			})
			It("should set a new event id to the trace header with different causation id", func() {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("X-Homes-Trace", traceJSONWithDiffCausation)
				w := httptest.NewRecorder()
				r := newServerWithTraceHeader()
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(200))
				traceHeader := w.Header().Get("X-Homes-Trace")
				Expect(traceHeader).NotTo(Equal(""))
				newTrace := trace.FromJSON(traceHeader)
				Expect(newTrace).NotTo(BeNil())
				Expect(newTrace.EventID).NotTo(Equal(""))
				Expect(newTrace.EventID).NotTo(Equal(newTrace.CorrelationID))
				Expect(newTrace.EventID).NotTo(Equal(newTrace.CausationID))
				Expect(newTrace.CorrelationID).NotTo(Equal(newTrace.CausationID))
				Expect(newTrace.CausationID).To(Equal(traceWithDiffCausation.EventID))

				body, _ := ioutil.ReadAll(w.Body)
				Expect(string(body)).To(Equal("response"))
			})
		})
	})
})

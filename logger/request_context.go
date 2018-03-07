package logger

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type contextKey int

const (
	contextKeyRequestContext contextKey = iota
)

type requestContext struct {
	RemoteAddr string
	RequestID  string
}

func FieldsFromContext(ctx context.Context) logrus.Fields {
	rctx, ok := retrieveRequestContext(ctx)
	if !ok {
		return logrus.Fields{}
	}

	return logrus.Fields{
		"remote":     rctx.RemoteAddr,
		"request_id": rctx.RequestID,
	}
}

func populateRequestContext(ctx context.Context, r *http.Request) context.Context {
	remoteAddr := r.RemoteAddr
	if realIP := getFirst(r.Header, "X-Real-IP", "X-Forwarded-For"); realIP != "" {
		remoteAddr = realIP
	}

	rctx := requestContext{
		RemoteAddr: remoteAddr,
		RequestID:  getFirst(r.Header, "X-Request-Id", "X-Amzn-Trace-Id"),
	}

	return context.WithValue(ctx, contextKeyRequestContext, rctx)
}

func retrieveRequestContext(ctx context.Context) (requestContext, bool) {
	rctx, ok := ctx.Value(contextKeyRequestContext).(requestContext)
	return rctx, ok
}

func getFirst(h http.Header, names ...string) string {
	for _, name := range names {
		if v := h.Get(name); v != "" {
			return v
		}
	}
	return ""
}

package logger

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type contextKey int

const (
	ContextKeyRequestContext contextKey = iota
)

type RequestContext struct {
	RemoteAddr string
	RequestID  string
}

func FieldsFromContext(ctx context.Context) logrus.Fields {
	rctx := retrieveRequestContext(ctx)

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

	requestContext := RequestContext{
		RemoteAddr: remoteAddr,
		RequestID:  getFirst(r.Header, "X-Request-Id", "X-Amzn-Trace-Id"),
	}

	return context.WithValue(ctx, ContextKeyRequestContext, requestContext)
}

func retrieveRequestContext(ctx context.Context) RequestContext {
	rctx := ctx.Value(ContextKeyRequestContext)
	if rctx == nil {
		return RequestContext{}
	}

	return rctx.(RequestContext)
}

func getFirst(h http.Header, names ...string) string {
	for _, name := range names {
		if v := h.Get(name); v != "" {
			return v
		}
	}
	return ""
}

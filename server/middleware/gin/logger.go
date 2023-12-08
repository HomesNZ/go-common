package gin

import (
	"bytes"
	"context"
	"github.com/HomesNZ/go-common/trace"
	"github.com/gin-gonic/gin"

	"net/http"
	"strings"
	"time"
)

const (
	homesTraceHeader = "X-Homes-Trace" // => {event_id: aa, correlation_id: aa, causation_id: aa}

)

type logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
}

// LoggerMiddleware
// log - logger instance to use for logging
// endpoints - endpoints to skip logging
func LoggerMiddleware(log logger, endpoints ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := 0
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		host := c.Request.Host
		method := c.Request.Method

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// trace logic
		reqCtx := c.Request.Context()
		traceHeader := c.Request.Header.Get(homesTraceHeader)

		// check if trace header is present in the request header
		if traceHeader == "" {
			// set a new event id for the request and set it to the context
			// or create a new trace if the trace is nil
			tracedCtx := trace.LinkCtx(reqCtx)
			c.Request = c.Request.WithContext(tracedCtx)
			reqCtx = tracedCtx
		} else {
			// set a new event id for the request and set it to the context
			tracedCtx := trace.LinkCtxFromJSON(reqCtx, traceHeader)
			reqCtx = tracedCtx
			c.Request = c.Request.WithContext(tracedCtx)
		}

		// Process request
		c.Next()

		// add trace header to the response
		c.Header(homesTraceHeader, trace.ToJSONFromCtx(reqCtx))

		endpoints = append(endpoints, "health") // skip health endpoint
		if skip(path, endpoints...) {
			statusCode = c.Writer.Status()
			errMsg := ""
			//stack := ""
			//type stackTracer interface {
			//	StackTrace() errors.StackTrace
			//}
			//extract message error from the response body
			if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError || statusCode >= http.StatusInternalServerError {
				if len(c.Errors) > 0 {
					//for _, e := range c.Errors {
					//	if err, ok := e.Err.(stackTracer); ok {
					//		stack = fmt.Sprintf("%+v", err.StackTrace())
					//	}
					//}
					errMsg = c.Errors.String()
				}
			}

			end := time.Now()
			latency := end.Sub(start)
			timeStamp := time.Now()

			switch {
			case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
				{
					log.Warn(
						reqCtx, errMsg,
						"time", timeStamp,
						"status", statusCode,
						"method", method,
						"host", host,
						"path", path,
						"latency", latency,
					)

				}
			case statusCode >= http.StatusInternalServerError:
				{
					log.Error(
						reqCtx, errMsg,
						"time", timeStamp,
						"status", statusCode,
						"method", method,
						"host", host,
						"path", path,
						"latency", latency,
					)
				}
			default:
				log.Info(
					reqCtx, errMsg,
					"time", timeStamp,
					"status", statusCode,
					"method", method,
					"host", host,
					"path", path,
					"latency", latency,
				)
			}
		}
	}
}

func skip(str string, arr ...string) bool {
	for idx := range arr {
		if strings.EqualFold(arr[idx], str) {
			return false
		}
	}
	return true
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

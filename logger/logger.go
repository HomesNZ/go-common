package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	config *Config
	once   sync.Once
)

type Logger struct {
	handler slog.Handler
}

// New constructs a new log for application use.
func New(opts ...Option) *Logger {
	config = &Config{
		level:       LevelDebug,
		writer:      os.Stdout,
		hooks:       Hooks{},
		serviceName: "unknown",
	}

	for _, opt := range opts {
		opt(config)
	}

	return new(config)
}

// Debug logs at LevelDebug with the given context.
// ctx can be used to add key/value pairs to the log. Can be nil.
// msg is the message to log.
// args are key/value pairs. Where key is a string and value is any type.
// If value is an error, the stack trace is added to the log only if the log is configured with WithStackTrace().
// Example: log.Debug(ctx, "message", "id(Int)", 12, "request", Request{}, "error_stack", err)
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.write(ctx, LevelDebug, 6, msg, args...)
}

// Info logs at LevelInfo with the given context.
// ctx can be used to add key/value pairs to the log. Can be nil.
// msg is the message to log.
// args are key/value pairs. Where key is a string and value is any type.
// If value is an error, the stack trace is added to the log only if the log is configured with WithStackTrace().
// Example: log.Debug(ctx, "message", "id(Int)", 12, "request", Request{}, "error_stack", err)
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.write(ctx, LevelInfo, 3, msg, args...)
}

// Warn logs at LevelWarn with the given context.
// ctx can be used to add key/value pairs to the log. Can be nil.
// msg is the message to log.
// args are key/value pairs. Where key is a string and value is any type.
// If value is an error, the stack trace is added to the log only if the log is configured with WithStackTrace().
// Example: log.Debug(ctx, "message", "id(Int)", 12, "request", Request{}, "error_stack", err)
func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.write(ctx, LevelWarn, 3, msg, args...)
}

// Error logs at LevelError with the given context.
// ctx can be used to add key/value pairs to the log. Can be nil.
// msg is the message to log.
// args are key/value pairs. Where key is a string and value is any type.
// If value is an error, the stack trace is added to the log only if the log is configured with WithStackTrace().
// Example: log.Debug(ctx, "message", "id(Int)", 12, "request", Request{}, "error_stack", err)
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.write(ctx, LevelError, 4, msg, args...)
}

// private
// =============================================================================

func (l *Logger) write(ctx context.Context, level Level, caller int, msg string, args ...any) {
	slogLevel := slog.Level(level)

	if !l.handler.Enabled(ctx, slogLevel) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])

	r := slog.NewRecord(time.Now(), slogLevel, msg, pcs[0])

	//if l.traceIDFunc != nil {
	//	args = append(args, "trace_id", log.traceIDFunc(ctx))
	//}
	r.Add(args...)

	l.handler.Handle(ctx, r)
}

func new(cfg *Config) *Logger {
	// Convert the file name to just the name.ext when this key/value will
	// be logged.
	f := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				v := fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line)
				return slog.Attr{Key: "file", Value: slog.StringValue(v)}
			}
		} else if a.Key == slog.TimeKey {
			if t, ok := a.Value.Any().(time.Time); ok {
				v := t.Format(time.RFC3339)
				return slog.Attr{Key: "time", Value: slog.StringValue(v)}
			}

		}

		if cfg.withStackTrace { // add stack trace
			switch a.Value.Kind() {
			case slog.KindAny:
				switch v := a.Value.Any().(type) {
				case error:
					a.Value = errToValue(v)
				}
			}
		}

		return a
	}

	// Construct the slog JSON handler for use.
	handler := slog.Handler(slog.NewJSONHandler(cfg.writer, &slog.HandlerOptions{AddSource: true, Level: slog.Level(cfg.level), ReplaceAttr: f}))

	// If events are to be processed, wrap the JSON handler around the custom
	// log handler.
	if cfg.hooks.Debug != nil || cfg.hooks.Info != nil || cfg.hooks.Warn != nil || cfg.hooks.Error != nil {
		handler = newLogHandler(handler, cfg.hooks)
	}

	// Attributes to add to every log.
	attrs := []slog.Attr{
		{Key: "service", Value: slog.StringValue(cfg.serviceName)},
	}

	// Add those attributes and capture the final handler.
	handler = handler.WithAttrs(attrs)

	return &Logger{
		handler: handler,
	}
}

// fmtErr returns a slog.GroupValue with keys "msg" and "trace". If the error
// does not implement interface { StackTrace() errors.StackTrace }, the "trace"
// key is omitted.
func errToValue(err error) slog.Value {
	var attr []slog.Attr
	attr = append(attr, slog.String("msg", err.Error()))

	type StackTracer interface {
		StackTrace() errors.StackTrace
	}

	// Find the trace to the location of the first errors.New,
	// errors.Wrap, or errors.WithStack call.
	var st StackTracer
	for err := err; err != nil; err = errors.Unwrap(err) {
		if x, ok := err.(StackTracer); ok {
			st = x
		}
	}

	if st != nil {
		attr = append(attr,
			slog.Any("trace", traceLines(st.StackTrace())),
		)
	}

	return slog.GroupValue(attr...)
}

func traceLines(frames errors.StackTrace) []string {
	traceLines := make([]string, len(frames))

	// Iterate in reverse to skip uninteresting, consecutive runtime frames at
	// the bottom of the trace.
	var skipped int
	skipping := true
	for i := len(frames) - 1; i >= 0; i-- {
		// Adapted from errors.Frame.MarshalText(), but avoiding repeated
		// calls to FuncForPC and FileLine.
		pc := uintptr(frames[i]) - 1
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			traceLines[i] = "unknown"
			skipping = false
			continue
		}

		name := fn.Name()

		// Skip runtime frames that occur before the first frame of the
		if skipping && strings.HasPrefix(name, "runtime.") {
			skipped++
			continue
		} else {
			skipping = false
		}

		filename, lineNr := fn.FileLine(pc)

		traceLines[i] = fmt.Sprintf("%s %s:%d", name, filename, lineNr)
	}

	return traceLines[:len(traceLines)-skipped]
}

package logger

import (
	"context"
	"log/slog"
)

// logHandler wraps a slog.Handler and adds hook functionality.
type logHandler struct {
	handler slog.Handler
	hooks   Hooks
}

// newLogHandler creates a new logHandler with the provided slog.Handler and Hooks.
func newLogHandler(handler slog.Handler, hooks Hooks) *logHandler {
	return &logHandler{
		handler: handler,
		hooks:   hooks,
	}
}

// Enabled reports whether the handler handles records at the given level.
// the handler ignores records below the given level.
func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns a new JSONHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{handler: h.handler.WithAttrs(attrs), hooks: h.hooks}
}

// WithGroup returns a new Handler with the given group appended to the receiver's
// existing groups. The keys of all subsequent attributes, whether added by With
// or in a Record, should be qualified by the sequence of group names.
func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{handler: h.handler.WithGroup(name), hooks: h.hooks}
}

// Handle looks to see if an event function needs to be executed for a given
// log level and then formats its argument Record.
func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelDebug:
		if len(h.hooks.Debug) > 0 {
			for _, hook := range h.hooks.Debug {
				hook(ctx, toRecord(r))
			}
		}

	case slog.LevelError:
		if len(h.hooks.Error) > 0 {
			for _, hook := range h.hooks.Error {
				hook(ctx, toRecord(r))
			}
		}

	case slog.LevelWarn:
		if len(h.hooks.Warn) > 0 {
			for _, hook := range h.hooks.Warn {
				hook(ctx, toRecord(r))
			}
		}

	case slog.LevelInfo:
		if len(h.hooks.Warn) > 0 {
			for _, hook := range h.hooks.Info {
				hook(ctx, toRecord(r))
			}
		}
	}
	return h.handler.Handle(ctx, r)
}

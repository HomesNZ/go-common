package logger

import (
	"context"
	"io"
	"os"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
	bugsnagErrors "github.com/bugsnag/bugsnag-go/v2/errors"
)

type Config struct {
	serviceName    string
	level          Level
	withBugsnag    bool
	writer         io.Writer
	hooks          Hooks
	withStackTrace bool
}

type Option func(*Config) error

// WithLevel sets the log level. Defaults to "info".
func WithLevel(level string) Option {
	return func(l *Config) error {
		l.level = toLevel(level)
		return nil
	}
}

// WithEnv sets the log level based on the environment. Defaults to "debug".
func WithEnv() Option {
	return func(l *Config) error {
		env := os.Getenv("ENV")
		if env == "production" || env == "prod" {
			l.level = LevelInfo
		} else {
			l.level = LevelDebug
		}
		level := os.Getenv("LOG_LEVEL")
		if level != "" {
			l.level = toLevel(level)
		}

		return nil
	}
}

// WithStackTrace adds a stack trace to the log. Defaults to false.
func WithStackTrace() Option {
	return func(l *Config) error {
		l.withStackTrace = true
		return nil
	}
}

const skipFrames = 4

// WithBugsnag adds a bugsnag hook to the logger. Defaults to false.
// It sends messages to bugsnag when an error occurs with stack traces.
func WithBugsnag() Option {
	return func(l *Config) error {
		l.withBugsnag = true
		l.hooks.Error = append(l.hooks.Error, func(ctx context.Context, r Record) {
			bugsnag.Notify(bugsnagErrors.New(r.ToError(), skipFrames))
		})
		return nil
	}
}

// WithServiceName sets the service name for the logger. Defaults to "service".
func WithServiceName(name string) Option {
	return func(l *Config) error {
		l.serviceName = name
		return nil
	}
}

// WithWriter sets the writer for the logger. Defaults to os.Stdout.
// This is useful for testing. For example, you can set the writer to a bytes.Buffer
// and then assert on the contents of the buffer.
func WithWriter(w io.Writer) Option {
	return func(l *Config) error {
		l.writer = w
		return nil
	}
}

// WithHooks adds hooks to the logger. Hooks are executed in the order they are added.
// Hooks can be used to execute custom logic when a log event occurs (e.g. send an email).
func WithHooks(events Hooks) Option {
	return func(l *Config) error {
		l.hooks.Error = append(l.hooks.Error, events.Error...)
		l.hooks.Warn = append(l.hooks.Warn, events.Warn...)
		l.hooks.Debug = append(l.hooks.Debug, events.Debug...)
		l.hooks.Info = append(l.hooks.Info, events.Info...)
		return nil
	}
}

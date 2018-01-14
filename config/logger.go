package config

import (
	"os"

	"github.com/HomesNZ/go-common/env"
	"github.com/Sirupsen/logrus"
	bugsnag "github.com/bugsnag/bugsnag-go"
	bugsnagErrors "github.com/bugsnag/bugsnag-go/errors"
	"github.com/pkg/errors"
)

// InitLogger initializes the logger by setting the log level to the env var LOG_LEVEL, or defaulting to `info`.
func InitLogger() {
	// If running in the production environment, output the logs as JSON format for parsing by Logstash.
	if env.IsProd() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	logrus.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(env.GetString("LOG_LEVEL", "info"))
	// No need to handle the error here, just don't update the log level
	if err == nil {
		logrus.SetLevel(level)
	}

	// Hooks
	logrus.AddHook(bugsnagHook{})

	logrus.Infof("Log level: %s", logrus.GetLevel().String())
}

type bugsnagHook struct{}

// Levels returns the logging levels that this hook will be fired for.
func (b bugsnagHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// skipFrames is the number of stack frames to skip in the error given to bugsnag
const skipFrames = 4

// Fire sends the logrus entry to bugsnag.
func (b bugsnagHook) Fire(entry *logrus.Entry) error {
	var err error
	switch er := entry.Data[logrus.ErrorKey].(type) {
	case stackTracer:
		err = stackError{er}
	case error:
		err = er
	default:
		err = errors.New(entry.Message)
	}
	notify := bugsnagErrors.New(err, skipFrames)
	meta := bugsnag.MetaData{}
	for field, value := range entry.Data {
		if field == logrus.ErrorKey {
			continue
		}
		meta.Add("logrus", field, value)
	}
	return bugsnag.Notify(notify, meta)
}

type stackError struct {
	stackTracer
}

var _ bugsnagErrors.ErrorWithCallers = stackError{}

func (e stackError) Callers() []uintptr {
	trace := e.StackTrace()
	callers := make([]uintptr, len(trace))
	for i, frame := range e.StackTrace() {
		callers[i] = uintptr(frame)
	}
	return callers
}

type stackTracer interface {
	error
	StackTrace() errors.StackTrace
}

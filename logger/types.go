package log

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type Level slog.Level

const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

func (l Level) String() string {
	return slog.Level(l).String()
}

func toLevel(l string) Level {
	switch strings.ToLower(l) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

type Record struct {
	Time       time.Time
	Level      Level
	Message    string
	Attributes map[string]any
}

func (r Record) ToError() error {
	var b strings.Builder
	b.WriteString("\"time\": ")
	b.WriteString(r.Time.String())
	b.WriteString(", ")
	b.WriteString("\"level\": ")
	b.WriteString(r.Level.String())
	b.WriteString(", ")
	if r.Message != "" {
		b.WriteString("\"message\": ")
		b.WriteString(r.Message)
		b.WriteString(", ")
	}

	if len(r.Attributes) > 0 {
		for k, v := range r.Attributes {
			switch v.(type) {
			case error:
				res := errToValue(v.(error))
				b.WriteString("\"" + k + "\": ")
				b.WriteString(res.String())
			default:
				b.WriteString("\"" + k + "\": ")
				b.WriteString(fmt.Sprint(v))
			}
		}
	}
	return fmt.Errorf(b.String())
}

func toRecord(r slog.Record) Record {
	atts := make(map[string]any, r.NumAttrs())
	f := func(attr slog.Attr) bool {
		atts[attr.Key] = attr.Value.Any()
		return true
	}
	r.Attrs(f)

	return Record{
		Time:       r.Time,
		Message:    r.Message,
		Level:      Level(r.Level),
		Attributes: atts,
	}
}

type HookFunc func(ctx context.Context, r Record)

type Hooks struct {
	Debug []HookFunc
	Info  []HookFunc
	Warn  []HookFunc
	Error []HookFunc
}

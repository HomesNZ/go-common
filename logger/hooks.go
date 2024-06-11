package logger

import "context"

type HookFunc func(ctx context.Context, r Record)

type Hooks struct {
	Debug []HookFunc
	Info  []HookFunc
	Warn  []HookFunc
	Error []HookFunc
}

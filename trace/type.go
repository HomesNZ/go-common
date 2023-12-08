package trace

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
)

// case 1:
// http request to server with trace header
// middleware extracts the trace header and sets it to the context
// server sends a message to sns topic

// case 2:
//  http request to server with trace header
// middleware extracts the trace header and sets it to the context
// server sends a message to another service by http

// case 3:
// sqs message received
// trace is in the message, set it to the context
// server sends a message to another service by http

// case 4:
// sqs message received
// trace is in the message, set it to the context
// server sends a message to sns topic

// ctxKey is the key for the trace in the context. Should be unexported to prevent collisions with context keys defined in other packages.
type ctxKey int

// key is used to store the trace in the context. Should be unexported to prevent collisions with context keys defined in other packages.
const key ctxKey = 61

// Trace represent state for each request.
type Trace struct {
	EventID       string `json:"event_id,omitempty"`       // event ID is used to uniquely identify an event.
	CorrelationID string `json:"correlation_id,omitempty"` // correlation ID is used to group related interactions that are part of the same higher-level operation. It remains constant across a transaction or a process involving multiple steps and possibly multiple microservices.
	CausationID   string `json:"causation_id,omitempty"`   // causation ID is used to track the cause of an event. It directly links an event to the action that caused it.
}

// ToJSON returns the trace header from the trace. If the trace is nil, it will return an empty string.
// trace - Trace
func (t Trace) ToJSON() string {
	b, err := json.Marshal(t)
	if err != nil || string(b) == "{}" {
		return ""
	}

	return string(b)
}

func (t Trace) IsEmpty() bool {
	return t.EventID == "" && t.CorrelationID == "" && t.CausationID == ""
}

// New creates a new trace.
func New() Trace {
	id := uuid.NewString()
	return Trace{
		EventID:       id,
		CorrelationID: id,
		CausationID:   id,
	}
}

// FromCtx returns the trace from the context.
func FromCtx(ctx context.Context) Trace {
	v, ok := ctx.Value(key).(Trace)
	if !ok {
		return Trace{}
	}
	return v
}

// SetToCtx sets the trace to the context. If the trace is nil, it will create a new trace.
func SetToCtx(ctx context.Context, trace Trace) context.Context {
	if trace.IsEmpty() {
		trace = New()
	}
	return context.WithValue(ctx, key, trace)
}

// SetToCtxFromHeader sets the trace from the header. If the header is empty, it will create a new trace.
// ctx - context
// traceHeader - trace header, should be in the format of {event_id: aa, correlation_id: aa, causation_id: aa}
func SetToCtxFromHeader(ctx context.Context, traceHeader string) context.Context {
	if traceHeader == "" {
		return SetToCtx(ctx, Trace{})
	}

	trace := Trace{}
	err := json.Unmarshal([]byte(traceHeader), &trace)
	if err != nil {
		return SetToCtx(ctx, trace)
	}

	return SetToCtx(ctx, trace)
}

// ToJSONFromCtx returns the trace header from the context. If the trace is nil, it will return an empty string.
// ctx - context
func ToJSONFromCtx(ctx context.Context) string {
	trace := FromCtx(ctx)
	if trace.IsEmpty() {
		return ""
	}

	b, err := json.Marshal(trace)
	if err != nil {
		return ""
	}

	return string(b)
}

// FromJSON returns the trace from the trace header. If the trace header is empty, it will return nil.
func FromJSON(traceJSON string) Trace {
	trace := Trace{}
	err := json.Unmarshal([]byte(traceJSON), &trace)
	if err != nil {
		return trace
	}
	return trace
}

// LinkFromTrace chain the trace. If the trace is nil, it will create a new trace.
// if the trace is not nil, it will create a new event id and set it to the trace.
// it will also set the causation id to the previous event id.
// tr - trace
func LinkFromTrace(tr Trace) Trace {
	tr = link(tr)
	return tr
}

func LinkFromJSON(traceJSON string) Trace {
	trace := FromJSON(traceJSON)
	trace = link(trace)
	return trace
}

// LinkCtxFromJSON chain the trace. If the trace header is empty, it will create a new trace.
// if the trace header is not empty, it will create a new event id and set it to the trace header.
// it will also set the causation id to the previous event id.
// ctx - context
// traceHeader - trace header, should be in the format of {event_id: aa, correlation_id: aa, causation_id: aa}
func LinkCtxFromJSON(ctx context.Context, traceHeader string) context.Context {
	trace := Trace{}
	err := json.Unmarshal([]byte(traceHeader), &trace)
	if err != nil {
		return SetToCtx(ctx, trace)
	}

	trace = link(trace)
	return SetToCtx(ctx, trace)
}

// LinkCtxFromCtx chain the trace. If the trace is nil, it will create a new trace.
// if the trace is not nil, it will create a new event id and set it to the trace.
// it will also set the causation id to the previous event id.
// ctx - context
func LinkCtxFromCtx(ctx context.Context) context.Context {
	trace := FromCtx(ctx)
	trace = link(trace)
	return SetToCtx(ctx, trace)
}

// LinkCtxFromTrace chain the trace. If the trace is nil, it will create a new trace.
// if the trace is not nil, it will create a new event id and set it to the trace.
// it will also set the causation id to the previous event id.
// ctx - context
// trace - trace
func LinkCtxFromTrace(ctx context.Context, tr Trace) context.Context {
	tr = link(tr)
	return SetToCtx(ctx, tr)
}

func link(tr Trace) Trace {
	if tr.IsEmpty() {
		return New()
	}

	tr.CausationID = tr.EventID
	tr.EventID = uuid.NewString()
	return tr
}

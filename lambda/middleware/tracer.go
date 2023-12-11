package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	homesEvents "github.com/HomesNZ/events"
	"github.com/HomesNZ/go-common/trace"
	"github.com/aws/aws-lambda-go/events"
)

const (
	attrHomesTrace = "homes_trace"
)

type SQSHandler func(ctx context.Context, event interface{}) error

func Trace(next SQSHandler) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		records := sqsEvent.Records
		var errors []error
		for idx := range records {
			fmt.Printf("MessageAttributes: %+v\n", records[idx].MessageAttributes)
			ev := struct {
				Message string `json:"Message"`
			}{}
			err := json.Unmarshal([]byte(records[idx].Body), &ev)
			if err != nil {
				return err
			}

			var msgTrace trace.Trace
			if records[idx].MessageAttributes != nil {
				if traceAttr, ok := records[idx].MessageAttributes[attrHomesTrace]; ok {
					msgTrace = trace.FromJSON(*traceAttr.StringValue) // only extract trace if it exists
				} else {
					msgTrace = trace.New()
				}
			}

			tracedCtx := trace.LinkCtxFromTrace(ctx, msgTrace)

			event, err := homesEvents.Parse([]byte(ev.Message))
			if err != nil {
				return err
			}

			if err := next(tracedCtx, event); err != nil {
				errors = append(errors, err)
			}

		}

		return errors[0]
	}
}

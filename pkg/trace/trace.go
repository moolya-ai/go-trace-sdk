package trace

import (
	"context"

	"github.com/google/uuid"
)

type key string

const TraceIDKey key = "traceID"
const TraceIDHeader = "X-Trace-ID"

// GenerateTraceID generates a new unique trace ID
func GenerateTraceID() string {
	return uuid.New().String()
}

// GetTraceID retrieves the trace ID from the context
func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	if !ok {
		return ""
	}
	return traceID
}

// SetTraceIDToContext adds the trace ID to the context
func SetTraceIDToContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

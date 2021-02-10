package otelhealth

import (
	"context"

	"go.opentelemetry.io/otel/label"
)

// Int64Counter is the interface for otel `metric.Int64Counter` struct.
type Int64Counter interface {
	Add(ctx context.Context, value int64, labels ...label.KeyValue)
}

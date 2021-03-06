package otelhealth

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// Int64Counter is the interface for otel `metric.Int64Counter` struct.
type Int64Counter interface {
	Add(ctx context.Context, value int64, labels ...attribute.KeyValue)
}

package otelhealth

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

// Int64Counter is the interface for otel `metric.Int64Counter` interface.
type Int64Counter interface {
	Add(ctx context.Context, value int64, opts ...metric.AddOption)
}

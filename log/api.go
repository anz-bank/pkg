package log

import (
	"context"

	"github.com/anz-bank/pkg/log/loggers"
)

type loggerContextKey int

const loggerKey loggerContextKey = iota

// With adds a copy of the logger to the context
func With(ctx context.Context, logger loggers.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger.Copy())
}

// WithFields adds multiple fields in the scope of the context, fields will be logged alphabetically
func WithFields(ctx context.Context, fields MultipleFields) context.Context {
	return context.WithValue(ctx, loggerKey,
		getCopiedLogger(ctx).PutFields(fromFields([]Fields{fields})))
}

// From is a way to access the API of the logger inside the context and add log-specific fields
func From(ctx context.Context, fields ...Fields) loggers.Logger {
	logger := getCopiedLogger(ctx)
	if len(fields) == 0 {
		return logger
	}
	return logger.PutFields(fromFields(fields))
}

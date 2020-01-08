package log

import (
	"context"

	"github.com/anz-bank/pkg/log/loggers"
)

type loggerContextKey int

const loggerKey loggerContextKey = iota

// WithLogger adds a copy of the logger to the context
func WithLogger(ctx context.Context, logger loggers.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger.Copy())
}

// The fields setup in WithField and WithFields are for context-specific fields
// Fields will be logged alphabetically

// WithField adds a single field in the scope of the context
func WithField(ctx context.Context, key string, val interface{}) context.Context {
	return context.WithValue(ctx, loggerKey,
		getCopiedLogger(ctx).PutField(key, val))
}

// WithFields adds multiple fields in the scope of the context
func WithFields(ctx context.Context, fields MultipleFields) context.Context {
	return context.WithValue(ctx, loggerKey,
		getCopiedLogger(ctx).PutFields(fromFields([]Fields{fields})))
}

// Logger is a way to access the API of the logger inside the context
func Logger(ctx context.Context, fields ...Fields) loggers.Logger {
	logger := getCopiedLogger(ctx)
	if len(fields) == 0 {
		return logger
	}
	return logger.PutFields(fromFields(fields))
}

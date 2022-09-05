package logging

import (
	"context"
	"errors"
	"time"

	"github.com/arr-ai/frozen"
	"github.com/rs/zerolog"
)

var (
	// ErrNoLoggerInContext is used when a log call occurs on a context with no logger
	ErrNoLoggerInContext = errors.New("log: context has no logger")
)

type loggerkeytype struct{}

var loggerkey = &loggerkeytype{}

// Context creates loggers from context
//
// This type acts as a go-between between loggers and context.
// Context is able to create a logger with values extracted from context
// according to the stored ContextFuncs.
//
// Context can also be added to context itself to support the main log
// functions (eg: logging.Info(ctx))
type Context struct {
	logger *Logger
	funcs  []ContextFunc
}

func (c Context) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerkey, &c)
}

// FromContext creates a logger from the given context and Context
func (c Context) FromContext(ctx context.Context) *Logger {
	logger := *c.logger
	fields := logger.internal.With()
	for _, fun := range c.funcs {
		fields = fun.Function(ctx, fields)
	}
	logger.internal = fields.Logger()
	return &logger
}

// FromContext returns a logger with values taken from context
//
// If no log context is defined, this will return a logger created from the global log context
func FromContext(ctx context.Context) *Logger {
	logCtx, ok := ctx.Value(loggerkey).(*Context)
	if !ok {
		return gblContext.FromContext(ctx)
	}
	return logCtx.FromContext(ctx)
}

// ContextFromContext returns a Context from context
//
// If no log context is defined, this will return the global log context
func ContextFromContext(ctx context.Context) *Context {
	logCtx, ok := ctx.Value(loggerkey).(*Context)
	if !ok {
		return gblContext
	}
	return logCtx
}

// With adds context funcs to the log context
func (c Context) With(funcs ...ContextFunc) *Context {
	for _, fun := range funcs {
		if c.logger.keys.Intersection(frozen.NewSet[string](fun.Keys...)).Count() == 0 {
			c.logger.keys = c.logger.keys.Union(frozen.NewSet[string](fun.Keys...))
			c.funcs = append(c.funcs, fun)
		} else {
			// warn of duplicate keys
			c.logger.Debug().Msg("duplicate keys detected in logger setup")
		}
	}
	return &c
}

// WithStr creates a static string field on the log context
func (c Context) WithStr(key string, val string) *Context {
	c.logger = c.logger.WithStr(key, val)
	return &c
}

// WithInt creates a static int field on the log context
func (c Context) WithInt(key string, val int) *Context {
	c.logger = c.logger.WithInt(key, val)
	return &c
}

// WithDict creates a static dictionary field (json object) on the log context
func (c Context) WithDict(key string, val *zerolog.Event) *Context {
	c.logger = c.logger.WithDict(key, val)
	return &c
}

// WithBool creates a static boolean field on the logger
func (c Context) WithBool(key string, val bool) *Context {
	c.logger = c.logger.WithBool(key, val)
	return &c
}

// WithArray creates a static array field on the logger
func (c Context) WithArray(key string, val zerolog.LogArrayMarshaler) *Context {
	c.logger = c.logger.WithArray(key, val)
	return &c
}

// WithDur creates a static duration field on the logger
func (c Context) WithDur(key string, val time.Duration) *Context {
	c.logger = c.logger.WithDur(key, val)
	return &c
}

// WithTime creates a static time field on the logger
func (c Context) WithTime(key string, val time.Time) *Context {
	c.logger = c.logger.WithTime(key, val)
	return &c
}

// Wrapper functions to provide nice single call logs with context.

// Info extracts the logger from context and logs a message at the Info level
func Info(ctx context.Context) *zerolog.Event {
	return FromContext(ctx).info(1)
}

// Error extracts the logger from context and logs a message at the Error level
func Error(ctx context.Context, err error) *zerolog.Event {
	return FromContext(ctx).error(err, 1)
}

// Debug extracts the logger from context and logs a message at the Debug level
func Debug(ctx context.Context, args ...interface{}) *zerolog.Event {
	return FromContext(ctx).debug(1)
}

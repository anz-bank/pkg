package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

// Fields is a struct that contains all the fields data to log.
type Fields struct{ m frozen.Map }

// Debug logs from context at the debug level
func Debug(ctx context.Context, args ...interface{}) {
	Fields{}.Debug(ctx, args...)
}

// Debugf logs from context at the debug level
func Debugf(ctx context.Context, format string, args ...interface{}) {
	Fields{}.Debugf(ctx, format, args...)
}

// From returns a copied logger from the context that you can use to access logger API.
func From(ctx context.Context) Logger {
	f := getFields(ctx)
	return f.configureLogger(ctx, f.getCopiedLogger().(fieldSetter))
}

// Info logs from context at the debug level
func Info(ctx context.Context, args ...interface{}) {
	Fields{}.Info(ctx, args...)
}

// Infof logs from context at the debug level
func Infof(ctx context.Context, format string, args ...interface{}) {
	Fields{}.Infof(ctx, format, args...)
}

// Suppress will ensure that suppressed keys are not logged.
func Suppress(keys ...string) Fields {
	return Fields{}.Suppress(keys...)
}

// With creates a field with a single key value pair.
func With(key string, val interface{}) Fields {
	return Fields{}.With(key, val)
}

// WithConfigs adds extra configuration for the logger.
func WithConfigs(configs ...Config) Fields {
	return Fields{}.WithConfigs(configs...)
}

// WithCtxRef creates a field with a key that refers to the provided context key,
// fields will use key as the fields property and take the value that corresponds
// to ctxKey.
func WithCtxRef(key string, ctxKey interface{}) Fields {
	return Fields{}.WithCtxRef(key, ctxKey)
}

// WithFunc creates a field with a string key and a callback value.
func WithFunc(key string, f func(context.Context) interface{}) Fields {
	return Fields{}.WithFunc(key, f)
}

// WithLogger adds logger which will be used for the log operation.
func WithLogger(logger Logger) Fields {
	return Fields{}.WithLogger(logger)
}

// Chain merges all the fields and returns the merged fields, the precedence
// of fields in case of overlapping gets higher from left to right.
func (f Fields) Chain(fieldses ...Fields) Fields {
	merged := f.m
	for _, fields := range fieldses {
		merged = merged.Update(fields.m)
	}
	return Fields{merged}
}

// Debug logs from context at the debug level
func (f Fields) Debug(ctx context.Context, args ...interface{}) {
	f.From(ctx).Debug(args...)
}

// Debugf logs from context at the debug level
func (f Fields) Debugf(ctx context.Context, format string, args ...interface{}) {
	f.From(ctx).Debugf(format, args...)
}

// From returns a logger with the new fields which is the fields from the context
// merged with the current fields were current fields replaces value from
// the context fields.
func (f Fields) From(ctx context.Context) Logger {
	return From(f.Onto(ctx))
}

// Info logs from context at the debug level
func (f Fields) Info(ctx context.Context, args ...interface{}) {
	f.From(ctx).Info(args...)
}

// Infof logs from context at the debug level
func (f Fields) Infof(ctx context.Context, format string, args ...interface{}) {
	f.From(ctx).Infof(format, args...)
}

// Onto finishes fields operation, merge them all with the precedence of fields
// in case overlapping gets higher from left to right, and puts the merged fields
// in the context.
func (f Fields) Onto(ctx context.Context) context.Context {
	return context.WithValue(ctx, fieldsContextKey{}, getFields(ctx).Chain(f).m)
}

// Suppress ensures that the keys will not be logger.
func (f Fields) Suppress(keys ...string) Fields {
	return f.Chain(Fields{
		frozen.NewMapFromKeys(
			frozen.NewSetFromStrings(keys...),
			func(_ interface{}) interface{} {
				return suppress{}
			},
		),
	})
}

// With adds to the fields a single key value pair.
func (f Fields) With(key string, val interface{}) Fields {
	return f.with(key, val)
}

// WithConfigs adds extra configuration for the logger.
func (f Fields) WithConfigs(configs ...Config) Fields {
	return f.Chain(Fields{
		createConfigMap(configs...),
	})
}

// WithCtxRef adds key and the context key to the fields.
func (f Fields) WithCtxRef(key string, ctxKey interface{}) Fields {
	return f.with(key, ctxRef{ctxKey})
}

// WithFunc adds key and the function to the fields.
func (f Fields) WithFunc(key string, val func(context.Context) interface{}) Fields {
	return f.with(key, val)
}

// WithLogger adds logger which will be used for the log operation.
func (f Fields) WithLogger(logger Logger) Fields {
	return f.with(loggerKey{}, logger.(copyable).Copy())
}

// String returns a string that represent the current fields
func (f Fields) String(ctx context.Context) string {
	fields := &fieldsCollector{}
	f.configureLogger(ctx, fields)
	return fields.fields.String()
}

// MergedString returns a string that represents the current fields merged by fields in context
func (f Fields) MergedString(ctx context.Context) string {
	return getFields(ctx).Chain(f).String(ctx)
}

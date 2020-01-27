package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

type fieldsContextKey struct{}
type loggerKey struct{}
type suppress struct{}
type ctxRef struct{ ctxKey interface{} }

func (f Fields) getCopiedLogger() Logger {
	logger, exists := f.m.Get(loggerKey{})
	if !exists {
		panic("Logger has not been added")
	}
	return logger.(internalLoggerOps).Copy()
}

func (f Fields) resolveFields(ctx context.Context) frozen.Map {
	fields := f.m
	var toSuppress frozen.SetBuilder
	toSuppress.Add(loggerKey{})
	for i := fields.Range(); i.Next(); {
		switch k := i.Value().(type) {
		case ctxRef:
			if val := ctx.Value(k.ctxKey); val != nil {
				fields = fields.With(i.Key(), val)
			} else {
				toSuppress.Add(i.Key())
			}
		case func(context.Context) interface{}:
			if val := k(ctx); val != nil {
				fields = fields.With(i.Key(), val)
			} else {
				toSuppress.Add(i.Key())
			}
		case suppress:
			toSuppress.Add(i.Key())
		}
	}
	return fields.Without(toSuppress.Finish())
}

func (f Fields) with(key, val interface{}) Fields {
	return Fields{f.m.With(key, val)}
}

func getFields(ctx context.Context) Fields {
	fields, exists := ctx.Value(fieldsContextKey{}).(frozen.Map)
	if !exists {
		return Fields{}
	}
	return Fields{fields}
}

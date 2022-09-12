package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

const errMsgKey = "error_message"

type fieldsContextKey struct{}
type loggerKey struct{}
type suppress struct{}
type ctxRef struct{ ctxKey interface{} }

type fieldsCollector struct {
	fields frozen.Map
}

func (f *fieldsCollector) PutFields(fields frozen.Map) Logger {
	f.fields = fields
	return NewNullLogger()
}

func (f Fields) getCopiedLogger() Logger {
	logger, exists := f.m.Get(loggerKey{})
	if !exists {
		return NewStandardLogger()
	}
	return logger.(copyable).Copy()
}

func (f Fields) configureLogger(ctx context.Context, logger fieldSetter) Logger {
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
		case Config:
			toSuppress.Add(i.Key())
			err := k.Apply(logger.(Logger))
			if err != nil {
				//TODO: should decide whether it should panic or not
				panic(err)
			}
		}
	}
	return logger.PutFields(fields.Without(toSuppress.Finish()))
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

func createConfigMap(configs ...Config) frozen.Map {
	var mb frozen.MapBuilder
	for _, c := range configs {
		mb.Put(c.TypeKey(), c)
	}
	return mb.Finish()
}

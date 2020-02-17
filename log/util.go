package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

const errMsgKey = "error_message"

type fieldsContextKey struct{}
type canonicalFieldsKey struct{}
type listenerKey struct {}
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

func configureLogger(logger fieldSetter, fields frozen.Map, configs frozen.Set) Logger {
	if configs.Count() > 0 {
		logger = applyConfiguration(logger.(Logger), configs).(fieldSetter)
	}
	return logger.PutFields(fields)
}

func applyConfiguration(logger Logger, configs frozen.Set) Logger {
	for c := configs.Range(); c.Next(); {
		err := c.(Config).Apply(logger.(Logger))
		if err != nil {
			//TODO: should decide whether it should panic or not
			panic(err)
		}
	}
	return logger
}

func (f Fields) getResolvedFields(ctx context.Context) (frozen.Map, frozen.Set) {
	fields := f.m
	var toSuppress, configBuilder frozen.SetBuilder
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
			configBuilder.Add(k)
		}
	}
	return fields.Without(toSuppress.Finish()), configBuilder.Finish()
}

func (f Fields) with(key, val interface{}) Fields {
	return Fields{f.m.With(key, val)}
}

func getCanonicalFields(ctx context.Context) *frozen.MapBuilder {
	mb, exists := ctx.Value(canonicalFieldsKey{}).(*frozen.MapBuilder)
	if !exists {
		return frozen.NewMapBuilder(0)
	}
	return mb
}

func addCanonicalFields(mb *frozen.MapBuilder, fields Fields) {
	for i := fields.m.Range(); i.Next(); {
		switch k := i.Key().(type) {
		case Config, Logger:
			mb.Put(k, i.Value())
		default:
			if f, exists := mb.Get(k); exists {
				if	val, isList := f.([]interface{}); isList {
					mb.Put(k, append(val, i.Value()))
				} else {
					mb.Put(k, []interface{}{f, i.Value()})
				}
			} else {
				mb.Put(k, i.Value())
			}
		}
	}
}

func doCallbacks(ctx context.Context, fields Fields) {
	callbacks := ctx.Value(listenerKey{})
	if callbacks != nil {
		for _, c := range callbacks.([]func(context.Context, Fields)) {
			c(ctx, fields)
		}
	}
}

func from(ctx context.Context, f Fields) Logger {
	fields, configs := f.getResolvedFields(ctx)
	doCallbacks(ctx, Fields{fields})
	return configureLogger(f.getCopiedLogger().(fieldSetter), fields, configs)
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

package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

const errMsgKey = "error_message"

type fieldsContextKey struct{}
type canonicalFieldsKey struct{}
type canonicalListenerKey struct {}
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
		panic("Logger has not been added")
	}
	return logger.(copyable).Copy()
}

func (f Fields) configureLogger(ctx context.Context, logger fieldSetter, configs frozen.Set) Logger {
	for c := configs.Range(); c.Next(); {
		err := c.(Config).Apply(logger.(Logger))
		if err != nil {
			//TODO: should decide whether it should panic or not
			panic(err)
		}
	}
	return logger.PutFields(f.m)
}

func (f Fields) applyConfiguration(logger Logger, configs frozen.Set) Logger {
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
		if f, exists := mb.Get(i.Key()); exists {
			if	val, isList := f.([]interface{}); isList {
				mb.Put(i.Key(), append(val, i.Value()))
			} else {
				mb.Put(i.Key(), []interface{}{f, i.Value()})
			}
		} else {
			mb.Put(i.Key(), i.Value())
		}
	}
}

func doCallbackIfRegistered(ctx context.Context, fields Fields, cb func(context.Context, Fields)) {
	callbacks, exists := ctx.Value(canonicalListenerKey{}).(frozen.Set)
	if !exists {
		return
	}
	if callbacks.Has(cb) {
		cb(ctx, fields)
	}
}

func from(ctx context.Context, f Fields, callback ...func(context.Context, Fields)) Logger {
	fields, configs := f.getResolvedFields(ctx)
	if len(callback) == 1 {
		doCallbackIfRegistered(ctx, Fields{fields}, callback[0])
	}
	return f.configureLogger(ctx, f.getCopiedLogger().(fieldSetter), configs)
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

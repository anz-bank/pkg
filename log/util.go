package log

import (
	"context"

	"github.com/arr-ai/frozen"
)

type fieldsContextKey struct{}
type loggerKey struct{}
type suppress struct{}
type ctxRef struct{ ctxKey interface{} }

type fieldsCollector struct{
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
		case configKey:
			toSuppress.Add(i.Key())
			setConfigToLogger(&logger, i.Key().(configKey), k)
		}
	}
	return logger.PutFields(fields.Without(toSuppress.Finish()))
}

func setConfigToLogger(logger *fieldSetter, configType, configSpec configKey) {
	switch configType {
	case formatter:
		(*logger).(formattable).SetFormatter(configSpec)
	default:
		panic("unknown configuration")
	}
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

func createConfigMap(configs ...config) frozen.Map {
	var mb frozen.MapBuilder
	for _, c := range configs {
		mb.Put(c.getConfigType(), c.getConfig())
	}
	return mb.Finish()
}

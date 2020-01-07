package log

import (
	"context"

	"github.com/anz-bank/pkg/log/loggers"
	"github.com/marcelocantos/frozen"
)

func getCopiedLogger(ctx context.Context) loggers.Logger {
	logger, exists := ctx.Value(loggerKey).(loggers.Logger)
	if !exists {
		panic("Logger does not exist in context")
	}
	return logger.Copy()
}

func fromFields(fields []Fields) frozen.Map {
	// set capacity to 0 because length of fields is undetermined
	frozenMapBuilder := frozen.NewMapBuilder(0)

	for _, f := range fields {
		for _, kv := range f.GetKeyValues() {
			frozenMapBuilder.Put(kv.Key, kv.Value)
		}
	}

	return frozenMapBuilder.Finish()
}

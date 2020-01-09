package main

import (
	"context"

	"github.com/anz-bank/pkg/log"
	"github.com/anz-bank/pkg/log/loggers"
)

func main() {
	// User is expected to choose a logger and add it to the context using the library's API
	ctx := context.Background()

	// this is a logger based on the logrus standard logger
	logger := loggers.NewStandardLogger()

	// With returns a new context
	ctx = log.With(ctx, logger)

	logLevelsDemo(ctx)
	logWithFieldsDemo(ctx)
}

func logLevelsDemo(ctx context.Context) {
	// logging requires the context variable so it must be given to any function that requires it
    log.From(ctx).Debug("Debug")
    log.From(ctx).Trace("Trace")
    log.From(ctx).Warn("Warn")
	log.From(ctx).Error("Error")

	// commented so it does not block the demo
    // log.From(ctx).Fatal("Fatal")
    // log.From(ctx).Panic("Panic")

    // Expected to log
    // (time in RFC3339Nano Format) (Level) (Message)
    //
    // Example:
    // 2019-12-12T08:23:59.210878+11:00 PRINT Hello There
    //
    // Each API also has its Format counterpart (Debugf, Printf, Tracef, Warnf, Errorf, Fatalf, Panicf)
}

func logWithFieldsDemo(ctx context.Context) {
	// context-level field adds fields to the context and creates a new context
    ctx = log.WithFields(ctx, map[string]interface{}{
        "just": "stuff",
        "stuff": 1,
    })

    // or alternatively
    // ctx = log.WithFields(ctx, MultipleFields{
    //     "just": "stuff",
    //     "stuff": 1
	// })

	// Any log at this point will have fields and to any function that uses the same context
    // just=stuff random=stuff stuff=1
    contextLevelField(ctx)
    logLevelField(ctx)
}

func contextLevelField(ctx context.Context) {
	// This is expected to log something like
    // 2019-12-12T08:23:59.210878+11:00 just=stuff random=stuff stuff=1 WARN Warn
    log.From(ctx).Warn("Warn")
}

func logLevelField(ctx context.Context) {
	// Log level fields are fields that are not stored into the context logger
    // Log level fields will add fields on top of the existing context level fields
    // If an existing key exists in the stored field, it will replace the value
    // This is expected to log something like
    // 2019-12-12T08:23:59.210878+11:00 just=stuff more=random stuff random=stuff stuff=1 very=random WARN Warn

    // You can add multiple fields at once through FromMap API or just use MultipleFields struct
    log.From(ctx,
        log.FromMap(map[string]interface{}{
            "more": "random stuff",
            "very": "random",
        },
    )).Warn("Warn")

    // or alternatively
    // log.From(
    //     ctx,
    //     log.MultipleFields{
    //         "more": "random stuff",
    //         "very": "random",
    //     },
    // ).Warn("Warn")

    // You can also add single field through the API NewField
    // log.From(ctx,
    //     log.NewField("more", "random stuff"),
    //     log.NewField("very", "random"),
    // ).Warn("Warn")

    // since context logger is not modified, it will log again only the context level fields
    contextLevelField(ctx)
}
package logging_test

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/anz-bank/pkg/logging"
	"github.com/rs/zerolog"
)

func testContext(out io.Writer) context.Context {
	logger := logging.New(out).WithLevel(logging.DebugLevel)
	return logger.With().ToContext(context.Background())
}

type keyStruct struct{}

var key = &keyStruct{}

func ExampleLogger_With() {
	// With adds a ContextFunc to the logger. This is invoked on every log call
	// The ID is a unique identifier for the function.
	// This function will log a 'context_value' if it exists in context
	ctx := logging.New(os.Stdout).With(
		logging.ContextFunc{
			Keys: []string{"context_value"},
			Function: func(ctx context.Context, logCtx zerolog.Context) zerolog.Context {
				ctxValue, ok := ctx.Value(key).(string)
				if ok {
					logCtx = logCtx.Str("context_value", ctxValue)
				}
				return logCtx
			},
		},
	).ToContext(context.Background())

	// this context does not have a value set yet so it won't be logged
	logging.Info(ctx).Msg("Hello World")

	// This context does have it. We don't need to manually extract the value
	ctx = context.WithValue(ctx, key, "I am contextual")
	logging.Info(ctx).Msg("Hello World")
	// Output: {"level":"info","message":"Hello World"}
	// {"level":"info","context_value":"I am contextual","message":"Hello World"}
}

//nolint:lll
func ExampleLogger_staticFields() {
	// logging provides a few convenience functions for static fields
	logger := logging.New(os.Stdout).
		WithStr("string", "foo").
		WithInt("some_int", 1).
		WithDict("metadata", logging.Dict().
			Str("service", "service_name").
			Str("version", "v1.0.0"),
		)

	logger.Info().Msg("Hello World")
	//Output: {"level":"info","string":"foo","some_int":1,"metadata":{"service":"service_name","version":"v1.0.0"},"message":"Hello World"}
}

func ExampleLogger_WithTimeDiff() {
	logger := logging.New(os.Stdout).WithTimeDiff("time_since_request_start", time.Now())

	// This log time_since_request_start field will be fairly small
	logger.Info().Msg("Hello World")

	// This log time_since_request_start field will have progressed by approx 10ms
	time.Sleep(10 * time.Millisecond)
	logger.Info().Msg("Hello World")
}

func ExampleInfo() {
	// Standard log calls will look like this.
	// MAKE SURE your context has a logger in it, otherwise it will panic
	logging.Info(testContext(os.Stdout)).Msg("Hello World")
	// Output: {"level":"info","message":"Hello World"}
}

func ExampleInfo_extraFields() {
	// You can add log specific fields
	logging.Info(testContext(os.Stdout)).Int("count", 1).Msg("Hello World")
	// Output: {"level":"info","count":1,"message":"Hello World"}
}

func ExampleInfo_embeddedObject() {
	// You can also embed an object
	logging.Info(testContext(os.Stdout)).
		Dict("subfields", logging.Dict().
			Str("field1", "foo").
			Int("field2", 123),
		).
		Msg("Hello World")
	// Output: {"level":"info","subfields":{"field1":"foo","field2":123},"message":"Hello World"}
}

func ExampleLogger_duplicateKeys() {
	// logger fields are protected from duplication
	logger := logging.New(os.Stdout).
		WithStr("foo", "bar").
		WithStr("foo", "baz")

	// log contains a single 'foo' field
	logger.Info().Msg("Hello World")

	// Unfortunately there is no duplication protection after context has been applied
	logger.Info().Str("foo", "baz").Msg("Hello World")
	// Output: {"level":"info","foo":"bar","message":"Hello World"}
	// {"level":"info","foo":"bar","foo":"baz","message":"Hello World"}
}

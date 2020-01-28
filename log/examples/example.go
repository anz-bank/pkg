package main

import (
	"context"
	"fmt"

	"github.com/anz-bank/pkg/log"
)

type contextKey struct{}
type contextKey2 struct{}

func main() {
	ctx := context.Background()
	fieldsDemo(ctx)
	loggingDemo(ctx)
	configLogDemo(ctx)
}

func fieldsDemo(ctx context.Context) {
	// In this library, fields are treated as an object that you can
	// manipulate and add to the context when ready, there are three
	// kinds of fields you can add.

	// With adds a regular key value pair.
	fields := log.With("hello", "world")

	// WithCtxRef adds a key whose value will be taken from the context
	// before logging. You have to define an alias to the key which will
	// be used during logging as context key are usually a struct or iota
	// which has no information about it when logged. If the key does not
	// exist in the context, it will not be logged.
	fields = fields.WithCtxRef("my alias", contextKey{})

	// WithFunc adds a key and a function with a context argument which
	// will be called before logging. If the result of the function is
	// nil, it will not be logged.
	ctx = context.WithValue(ctx, "bar", 42)
	fields = fields.WithFunc("foo", func(ctx context.Context) interface{} {
		return ctx.Value("bar")
	})
	fmt.Printf("Current Fields after using the With API: %s\n", fields.String(ctx))
	ctx = context.WithValue(ctx, contextKey{}, "now exist in context")
	fmt.Printf("Current Fields after using the With API and adding my alias to context: %s\n", fields.String(ctx))

	// Fields operation can also be chained either by the With APIs or Chain API.
	// An important thing to note is that fields operation always merge with the
	// previous fields, but when key overlaps, the newest fields will always replace
	// the overlapping keys. This means, in a chain operation or Chain API, the precedence
	// gets higher from left to right.

	// In this example, the final fields will have ("test": "test four") instead of ("test": "test too").
	fields = fields.
		With("test", "test too").
		WithCtxRef("test three", contextKey2{}).
		WithFunc("doesn't", func(context.Context) interface{} {
			return "matter"
		}).
		With("test", "test four")

	// The final fields will have ("out of": "things to write")
	fields1 := log.With("I'm", "running")
	fields2 := log.With("out of", "examples")
	fields3 := log.With("out of", "things to write")
	fields = fields.Chain(fields1, fields2, fields3)
	fmt.Printf("Current Fields after chaining: %s\n", fields.String(ctx))

	// Fields are also immutable so it will be thread-safe. Adding fields
	// to the context can be done using the Onto API which returns a context
	// that contains the fields. Onto will merge instead of replacing the fields.
	ctx = fields.Onto(ctx)
	fmt.Printf("Current Fields in context: %s\n", log.Fields{}.MergedString(ctx))

	// If you would like to suppress certain fields, you can use the Suppress API.
	// Just give it the keys that you would like to not be logged and it will ensure
	// that fields will be ignored during logging.
	fields = fields.Suppress("test", "doesn't")
	fmt.Printf("Suppressing: %s\n", fields.String(ctx))

	// Remember to finalize it with Onto if you want to remove it from the context. If you
	// would like to use it just for specific log, you don't need to finalize it.
	// If you want to unsuppress a value, you would have to add the value back using the With API.
	fmt.Printf("Current fields in context before finalizing the suppressed fields: %s\n", log.Fields{}.MergedString(ctx))
	ctx = fields.Onto(ctx)
	fmt.Printf("Current Fields in context after suppressing: %s\n", log.Fields{}.MergedString(ctx))

	// If you do not want to put the fields in the context, you can also just generate
	// a logger that contains all the fields with the From method. The logger will always
	// merge your current fields with the fields the context contains. But you need to make sure
	// that logger has been added. Logger is treated as another fields, so you can just chain it
	// with the WithLogger API. Trying to generate a logger without setting up the logger will make the
	// program panic. More about logging in another function.
	logger1, logger2 := log.NewNullLogger(), log.NewStandardLogger()

	// Creating a context that contains logger1.
	ctx = log.WithLogger(logger1).Onto(ctx)

	// Since logger is treated like fields, this log will use logger2 instead due to higher precedence.
	fields.With("additional", "fields").WithLogger(logger2).From(ctx).Info("Logging additional fields with logger2")

	// You can also just log directly from context. Logger will always take the fields
	// from context before logging.
	log.From(ctx).Info("Since this uses logger1 which is a null logger, this will not be logged")
}

func loggingDemo(ctx context.Context) {
	// You can choose different loggers or implement your own.
	// For this example, we are using the standard logger.
	logger := log.NewStandardLogger()

	// Adding logger can be done through the WithLogger API. Do not forget to finalize it with
	// the Onto API if you want the logger to be contained in the context.
	ctx = log.WithLogger(logger).Onto(ctx)

	// Logging from the context can be done using the From API which will take the logger
	// from the context and take any fields contained in the context and get them logged
	// by the logger. There are two levels of logging, Debug and Info. Each will also have
	// the format counterpart, Debugf and Infof.
	// Logging will log in the following format:
	// (time in RFC3339Nano Format) (Fields) (Level) (Message)
	// Fields themselves are logged as a space separated list of key=value.
	log.From(ctx).Debug("This does not have any fields")

	ctx = log.With("this", "one").With("have", "fields").Onto(ctx)
	log.From(ctx).Info("log with fields")

	// Sometimes, you might like to log with additional fields and not have to add them to the
	// context. This can be done with the From method which generates Logger with the
	// context fields merged with additional fields. Additional fields of course have higher
	// precedence and will replace context fields when keys overlap. But since the fields
	// is not added to the context, context fields will be untouched.
	log.With("have", "additional fields").With("log-specific", "fields").From(ctx).Debug("log-specific fields")
	log.From(ctx).Info("context fields are still untouched as long as the context is unchanged")
}

func configLogDemo(ctx context.Context) {
	// Adding configuration can be done by adding the correct struct. Configurations are once again
	// treated as fields, which means it will replace old configurations when a configuration
	// of the same type is added. For example, if before you added StandardFormatter, calling WithConfig
	// with JSONFormatter will replace StandardFormatter. Just like Fields, it will also be stored
	// in the context.
	ctx = log.WithConfigs(log.JSONFormatter{}).Onto(ctx)

	// You can also have a log-specific configs by not saving it to the context.
	log.WithConfigs(log.StandardFormatter{}, log.JSONFormatter{}).
		WithLogger(log.NewStandardLogger()).
		With("yeet", map[string]interface{}{"foo": "bar", "doesn't": "matter"}).
		From(ctx).
		Info("json formatted log")
}

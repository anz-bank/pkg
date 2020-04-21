# Examples of Usage

## Setup

Adding logger to the context

```go
package main

import (
    "context"

    "github.com/anz-bank/pkg/log"
    "github.com/anz-bank/pkg/log/loggers"
)

func main() {
    ctx := context.Background()
    logger := loggers.NewStandardLogger()

    // This adds logger to the context.
    // Trying to log without a logger inside the context will make it panic.
    ctx = log.WithLogger(logger).Onto(ctx)

    // Since logger is treated as fields you can also do this while adding fields.
    // You need to ensure that logger exist in the context or the fields
    // you are about to add to context (more about fields in the next section).
    ctx = log.With("key", "value").WithLogger(ctx).Onto(ctx)
}
```

## Fields

Fields are key value pair information that will be logged along with the log
message.

```go
func fieldsDemo(ctx context.Context) {
	// In this library, fields are treated as an object that you can
	// manipulate and add to the context when ready, there are three
	// kinds of fields you can add.

	// With adds a regular key value pair.
	fields := log.With("hello", "world")

	// WithContextRef adds a key whose value will be taken from the context
	// before logging. You have to define an alias to the key which will
	// be used during logging as context key are usually a struct or iota
	// which has no information about it when logged. If the key does not
	// exist in the context, it will not be logged.
	fields = fields.WithContextRef("my alias", contextKey{})

	// Fields operation can also be chained either by the With APIs or Chain API.
	// An important thing to note is that fields operation always merge with the
	// previous fields, but when key overlaps, the newest fields will always replace
	// the overlapping keys. This means, in a chain operation or Chain API, the precedence
	// gets higher from left to right.

	// In this example, the final fields will have ("test": "test four") instead of ("test": "test too").
	fields = fields.
		With("test", "test too").
		WithContextRef("test three", contextKey2{}).
		With("test", "test four")

	// The final fields will have ("out of": "things to write")
	fields1 := log.With("I'm", "running")
	fields2 := log.With("out of", "examples")
	fields3 := log.With("out of", "things to write")
	fields = fields.Chain(fields1, fields2, fields3)

	// Fields are also immutable so it will be thread-safe. Adding fields
	// to the context can be done using the Onto API which returns a context
	// that contains the fields. Onto will merge instead of replacing the fields.
	ctx = fields.Onto(ctx)

	// If you would like to suppress certain fields, you can use the Suppress API.
	// Just give it the keys that you would like to not be logged and it will ensure
	// that fields will be ignored during logging.
	fields = fields.Suppress("test", "doesn't")

	// Remember to finalize it with Onto if you want to remove it from the context. If you
	// would like to use it just for specific log, you don't need to finalize it.
	// If you want to unsuppress a value, you would have to add the value back using the With API.
	ctx = fields.Onto(ctx)

	// If you do not want to put the fields in the context, you can also just generate
	// a logger that contains all the fields with the From method. The logger will always
	// merge your current fields with the fields the context contains. But you need to make sure
	// that logger has been added. Logger is treated as another fields, so you can just chain it
	// with the WithLogger API. Trying to generate a logger without setting up the logger will make the
	// program panic. More about logging in another function.
	logger1, logger2 := loggers.NewNullLogger(), loggers.NewStandardLogger()

	// Creating a context that contains logger1.
	ctx = log.WithLogger(logger1).Onto(ctx)

	// Since logger is treated like fields, this log will use logger2 instead due to higher precedence.
	fields.With("additional", "fields").WithLogger(logger2).From(ctx).Info("Logging additional fields with logger2")

	// You can also just log directly from context. Logger will always take the fields
	// from context before logging.
	log.From(ctx).Info("Since this uses logger1 which is a null logger, this will not be logged")
}
```

## Logging

Logging can be done through the `From` API.

```go
func loggingDemo(ctx context.Context) {
	// You can choose different loggers or implement your own.
	// For this example, we are using the standard logger.
	logger := loggers.NewStandardLogger()

	// Adding logger can be done through the WithLogger API. Do not forget to finalize it with
	// the Onto API if you want the logger to be contained in the context.
	ctx = log.WithLogger(logger).Onto(ctx)

	// Logging from the context can be done using the From API which will take the logger
	// from the context and take any fields contained in the context and get them logged
	// by the logger. There are two levels of logging, Debug and Info. Each will also have
	// the format counterpart, Debugf and Infof.
	// Logging will log in the following format:
	// (time in RFC3339Nano Format) (Fields) (Level) (Message)
	// Fields themselves are logged as a space separated list of key=value
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

```

### Logging in different formats

You can specify JSON output by adding `WithConfig(NewJSONFormat())`. It will log in the following format:
```js
{
	"fields": {
		"key1": "value1", // value can be any data types
		"key2": "value2",
	},
	"level": "log level", // string, either INFO or DEBUG
	"message": "log message", // string,
	"timestamp": "log time", // timestamp in RFC3339Nano format
}
```

## Configuring logger

Logger can be configured using the `WithConfig` API and giving it the correct configuration struct.
```go
func configLogDemo(ctx context.Context) {
	// Adding configuration can be done by adding the correct struct. Configurations are once again
	// treated as fields, which means it will replace old configurations when a configuration
	// of the same type is added. For example, if before you added StandardFormatter, calling WithConfig
	// with JSONFormatter will replace StandardFormatter. Just like Fields, it will also be stored
	// in the context.
	ctx = log.WithConfigs(log.NewJSONFormat(), log.NewStderrOut(), log.SetVerboseMode(true), log.SetOutput(os.Stdout)).Onto(ctx)

	// You can also have a log-specific configs by not saving it to the context.
	log.WithConfigs(log.NewStandardFormat(), log.NewBufferOut(), log.SetVerboseMode(false)).
		WithLogger(log.NewStandardLogger()).
		With("key", map[string]interface{}{"foo": "bar", "doesn't": "matter"}).
		From(ctx).
		Info("json formatted log")
}
```

Code snippets can be run in the [example file](examples/example.go)

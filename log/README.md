# Contextual Logging

## Log
This library is a contextual logging library that makes use of context as part of the logging process. It is designed to makes development easier by using the context variable to log instead of passing a logger to every functions.

## Main Features
### Fields
Fields are key value data that are logged along with the log message. This library makes manipulating Fields easier and more flexible. With this library, everything is treated as Fields, that includes Fields themselves, logger configuration, and even the logger itself. This makes it possible for you to create everything in one chained operation.

```go
package fields

import(
    "context"
    
    "github.com/anz-bank/pkg/log"
)

type ctxKey struct{}

func fieldsDemo(ctx context.Context) {
    // This is adding a regular field.
    f := log.With("key1", "value1")
    
    // There is also the context fields where it will take values that correspond
    // to the given context key. The context key can be any object but you have
    // to provide the alias for the key. If the key does not have any value in
    // the context, it will not be logged.
    f = log.WithCtxRef("alias", ctxKey{})
    
    //TODO: add WithFunc or not?
    
    // You can also add multiple Fields by chaining the operation.
    f = f.
    	  With("another", "key").
          With("more", 1).
    	  With("more key", 'q')

    // One thing to remember, since fields are key value data, in the event of
    // overlapping keys, values will be replaced based on the order of operation.
    // In a chain operation, the later operations have higher precedence and will
    // replace the keys. At this example, the value that corresponds to "another"
    // is now "fields" instead of "key".
    f = f.With("another", "fields")

    // As mentioned before, everything is treated as fields and that includes
    // the logger and its configuration. Only one logger can be in fields and
    // one of each type of configurations (e.g. only one rule for format etc).
    // Because they are fields they will also follow the precedence rule, which
    // means adding another logger or a configuration type will replace the older
    // values.
    f = f.
    	  WithLogger(log.NewStandardLogger()).
          WithConfigs(log.NewJSONFormat())
    
    // The fields then can be used to log directly (example in later section) or
    // they can be saved in the context for later use by using the Onto API.
    newCtx := f.Onto(ctx)

    // A couple more useful APIs to know.

    // Suppress will ensure that the provided keys will not be logged.
    // In this example, the key "another", "more key", and "alias" will not 
    // be logged. For context reference fields, you have to refer to them
    // by their alias.
    f = f.Suppress("another", "more key", "alias")

    // ￿Chain provides a way of merging multiple fields. Just like before,
    // precedence gets higher from left to right.
    f1 := log.With("key1", "value1")
    f2 := log.With("key2", "value2")
    f3 := log.With("key3", "value3")
    f ￿= f.Chain(f1, f2, f3)
}
```

A very important thing to note is that Fields are immutable which makes them thread-safe but it also means that you need to receive the returned value of fields operation as they do not mutate themselves.

### Logging
There are three levels of logging, they are `Debug`, `Info`, and `Error`. Each of them also have their format function counterpart which are `Debugf`, `Infof`, and `Errorf`.

```go
package logging

import (
    "context"
    "errors"
    
    "github.com/anz-bank/pkg/log"
)

func logDemo(ctx context.Context){
    // Logging can be accessed through the Debug, Info, and Error API.
    // Each of the log functions require a context to be passed in.
    // If the context contains fields, that fields will be logged along
    // with the message given. If the context does not contain a logger
    // a standard logger will be provided as the default.
    log.Debug(ctx, "this is debug")
    log.Debugf(ctx, "%s with format", "this is debug")
    log.Info(ctx, "this is info")
    log.Infof(ctx, "%s with format", "this is info")

    // For Error and Errorf, the error variable is required. The error
    // message will be logged as a field with the key of "error_message"
    log.Error(ctx, errors.New("error"), "this is error")
    log.Errorf(ctx, errors.New("error"), "%s with format", "this is error")

    // If you would like to log certain fields without adding additional
    // fields to the context, you can do so by using the same API on the
    // additional fields. Additional fields are merged with the fields
    // in context if the context contains fields and it also has higher 
    // precedence but they do not mutate the fields in context.
    log.With("additional", "fields").With("more", "fields").Debug(ctx, "debug")
    log.With("additional", "fields").With("more", "fields").Debugf(ctx, "formatted %s", "debug")
    
    // This log will only log fields inside the context.
    log.Debug(ctx, "no additional fields")

    // Should you require the logger object itself, you can do so by using the
    // From API which will extract the logger in the context. If context does not
    // have any logger, it will returns a new standard logger. The returned logger
    // is copied for immutability. The logger returned by From have all the fields
    // and configuration applied to it. The fields are also resolved, meaning
    // any context reference will use any value in the context at the time of call.
    logger := log.From(ctx)
}
```

### Configuring logger
Logger configurations are treated as fields. This can be done through the `WithConfigs` API. You can add multiple configurations in a single `WithConfigs` operation. The configurations can also be saved in a context along with other fields. Even if you replace the logger, the configurations stay and will always be applied to the logger.
```go
package configs

import (
    "context"
    
    "github.com/anz-bank/pkg/log"
)

func configLog(ctx context.Context) {
    
}
```
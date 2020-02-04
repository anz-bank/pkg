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
    // to provide the alias for the key.
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
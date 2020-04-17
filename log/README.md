# Contextual Logging

## Intro

This library is a contextual logging library that makes use of context as part of the logging 
process. It is designed to make development easier by using the context variable to log instead of 
using a single global logger or passing a logger to every function.

## Getting Started

```go
package main
import (
    "context"

    "github.com/anz-bank/pkg/log"
)

func main() {
    ctx := context.Background()
    logger := log.NewStandardLogger()

    // Setup with context fields.
    ctx = log.WithLogger(logger).With("key1", "val1").With("key2", "val2").Onto(ctx)

    // This is how you log.
    log.Debug(ctx, "Hello There!")
    // This is how you log with temporary extra fields.
    log.With("temporary", "fields").Info(ctx, "What's poppin?")
}
```

## Why use this library

### Do a lot of things in one operation

The library focuses on doing multiple operations, whether it is adding fields or configuring the 
logger, in one chained operation. This makes setup and logging very simple, especially when they 
involve `Fields`.

### Shallow Context Tree

`Fields` are stored in the context tree when you use the `Onto` method. Doing many things in one 
operation allows you to produce a shallow context tree as you do not need to add `Fields` 
one-by-one. By finalizing the `Fields` operation using `Onto`, it ensures that it will only add all 
the provided `Fields` once. This is extremely beneficial when several fields must be logged, which 
is common in large and complex codebases.

### Greater control over `Fields` and `Logger`

There are also many operations you can do on `Fields` as the library allows you to store fields in a 
variable for finer control. The `With` methods allow many different types of `Fields` to be entered 
and APIs like [`Chain`](https://godoc.org/github.com/anz-bank/pkg/log#Fields.Chain) and 
[`Suppress`](https://godoc.org/github.com/anz-bank/pkg/log#Suppress) make `Fields` a lot more 
customizable.

### Immutability

The library ensures that `Fields` are immutable and the real `Logger` is never exposed.
Any access to the logger will return a copy. This is very beneficial in programs with concurrent 
processes.

### Customisable

The library provides a lot of ways to customize your logger to meet your needs. You can create your 
own configuration or even an entirely different logger. The provided interfaces are small which 
makes it really easy to create your own configurations to the library.

### Compared to other solutions

A very popular solution for logging in open source is the [logrus](https://github.com/sirupsen/logrus) 
library. While it is a great logging library, it does not provide a built-in `Fields` solution and a 
very high level of `Fields` manipulation. It also does not implement context properly as it requires 
you to create a custom format even after using their `WithContext` API. Finally, logrus has a large 
API, and while it provides a great amount of features, users can find it intimidating and confusing 
to use.

Compared to logrus, the library provides a built-in solution in implementing context, the provided 
default formats ensure that context values that you need are logged. The library also provides a 
simple set of APIs that are easy to use. Everything you need that involves a logger, you can easily 
find it.

## Main Features

### Fields

Fields are key value data that are logged along with the log message. This library makes 
manipulating Fields easier and more flexible. With this library, everything is treated a Field, 
that includes Fields themselves, logger configuration, and even the logger itself. This makes it 
possible for you to create everything in one chained operation.

```go
    f := log.With("key1", "value1")
```

There is also the context fields where it will take values that correspond to the given context key. 
The context key can be any object but you have to provide the alias for the key. If the key does not 
have any value in the context, it will not be logged.

```go
    f = log.WithCtxRef("alias", ctxKey{})
```

You can also add multiple Fields by chaining the operation.

```go
    f = f.
       With("another", "key").
          With("more", 1).
       With("more key", 'q')
```

One thing to remember, since fields are key value data, in the event of overlapping keys, values 
will be replaced based on the order of operation. In a chain operation, the later operations have 
higher precedence and will replace the keys. At this example, the value that corresponds to 
`another` is now `fields` instead of `key`.

```go
    f = f.With("another", "fields")
```

As mentioned before, everything is treated as fields and that includes the logger and its 
configuration. Only one logger can be in fields and one of each type of configurations (e.g. only 
one rule for format etc). Because they are fields they will also follow the precedence rule, which 
means adding another logger or a configuration type will replace the older values.

```go
    f = f.
       WithLogger(log.NewStandardLogger()).
          WithConfigs(log.NewJSONFormat())
```

The fields then can be used to log directly (example in later section) or they can be saved in the 
context for later use by using the `Onto` API.

```go
    newCtx := f.Onto(ctx)
```

A couple more useful APIs to know.

`Suppress` will ensure that the provided keys will not be logged. In this example, the key 
`another`, `more key`, and `alias` will not be logged. For context reference fields, you have to 
refer to them by their alias.

```go
    f = f.Suppress("another", "more key", "alias")
```

`Chain` provides a way of merging multiple fields. Just like before, precedence gets higher from 
left to right.

```go
    f1 := log.With("key1", "value1")
    f2 := log.With("key2", "value2")
    f3 := log.With("key3", "value3")
    f = f.Chain(f1, f2, f3)
```

A very important thing to note is that Fields are immutable which makes them thread-safe but it also 
means that you need to receive the returned value of fields operation as they do not mutate 
themselves.

### Logging

Logging can be accessed through the `Debug`, `Info`, and `Error` API.  Each of them also have their 
format function counterpart which are `Debugf`, `Infof`, and `Errorf`. `Debug` and `Debugf` is only 
logged when the logger is in the verbose mode while the others will always be logged. Each of the 
log functions require a context to be passed in. If the context contains fields, that fields will be 
logged along with the message given. If the context does not contain a logger a standard logger will 
be provided as the default.

```go
    log.Debug(ctx, "this is debug")
    log.Debugf(ctx, "%s with format", "this is debug")
    log.Info(ctx, "this is info")
    log.Infof(ctx, "%s with format", "this is info")
```

For `Error` and `Errorf`, the error variable is required. The error message will be logged as a 
field with the key of `error_message`.

```go
    log.Error(ctx, errors.New("error"), "this is error")
    log.Errorf(ctx, errors.New("error"), "%s with format", "this is error")
```

If you would like to log certain fields without adding additional fields to the context, you can do 
so by using the same API on the additional fields. Additional fields are merged with the fields in 
context if the context contains fields and it also has higher precedence but they do not mutate the 
fields in context.

```go
    log.With("additional", "fields").With("more", "fields").Debug(ctx, "debug")
    log.With("additional", "fields").With("more", "fields").Debugf(ctx, "formatted %s", "debug")

     // This log will only log fields inside the context.
    log.Debug(ctx, "no additional fields")
```

Should you require the logger object itself, you can do so by using the `From` API which will 
extract the logger in the context. If context does not have any logger, it will returns a new 
standard logger. The returned logger is copied for immutability. The logger returned by From have 
all the fields and configuration applied to it. The fields are also resolved, meaning any context 
reference will use any value in the context at the time of call.

```go
    logger := log.From(ctx)

    // This one will return a logger with the additional fields merged with context fields
    logger := log.With("extra", "fields").From(ctx)
```

### Configuring logger

Logger configurations are treated as fields. This can be done through the `WithConfigs` API. You can 
add multiple configurations in a single `WithConfigs` operation. The configurations can also be 
saved in a context along with other fields. Even if you replace the logger, the configurations stay 
and will always be applied to the logger. Only one type of each configuration type can exist in a 
fields. If another config of the same type is added, it will replace the old one.

```go
    // This adds the JSON formatter to the logger.
    f = log.WithConfigs(log.NewJSONFormat())

    // This will replace JSON formatter.
    f = log.WithConfigs(log.NewStandardFormat())

    // You can add multiple configurations
    f = log.WithConfigs(log.NewJSONFormat(), log.NewStandardFormat())
```

#### Logging Format

Currently there is only one logger which is the `StandardLogger` which uses 
[logrus](https://github.com/sirupsen/logrus). The provided formatter implements logrus formatter 
system. There are two formatters, the JSON formatter and the Standard formatter (which is the 
default formatter when no configuration is added).

##### JSON format

JSON formatter will log in the following format:

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

Fields will be logged as an object of the attribute `fields`. One thing to remember is that, for 
context reference, the key will use the provided alias.

##### Standard format

The standard formatter will log in the following format without the parentheses:

```text
(time in RFC3339Nano Format) (Fields) (Level) (Message)
```

For example:

```log
2020-02-05T09:05:11.041651+11:00 this=one have=fields INFO log with fields
```

In the current implementation, the fields are logged in a random order.

#### Verbosity

Setting the verbose mode of the logger will log debug entries:

```go
    ctx := context.Background()

    // By default, the logger will not log debug entries
    log.Info(ctx, "not logged")

    // Make the logger log debug level entries
    ctx = log.WithConfigs(log.SetVerboseMode(true)).Onto(ctx)
    
    // With verbose mode enabled, the logger will log debug entries
    log.Info(ctx, "logged")
```

#### Log Caller

Setting the logger to log the caller will include the a reference to the source from which the log
was called:

```go
    ctx := context.Background()

    // By default, the logger will not log the caller
    // 2020-02-05T09:05:11.041651+11:00 INFO one
    log.Info(ctx, "one")

    // Make the logger log the caller
    ctx = log.WithConfigs(log.SetLogCaller(true)).Onto(ctx)
    
    // With caller log enabled, the logger will log the caller
    // 2020-02-05T09:05:11.041651+11:00 INFO two [/path/to/example.go:42]
    log.Info(ctx, "two")
```

#### Hook

Hooks can be added to the logger that are notified when an entry is logged:

```go
    type myHook struct { }
    func (h *myHook) OnLogged(entry *LogEntry) error { ... }

    ctx = log.WithConfigs(log.AddHook(myHook{})).Onto(context.Background())
    log.Info(ctx, "message") // log entry sent to hook
```

#### Custom configuration

It is possible to create your own configuration. You will have to create an object that implements 
the provided interface.

```go
type Config interface {
 TypeKey() interface{}
 Apply(logger Logger) error
}
```

`TypeKey()` returns the type of the configuration and `Apply()` will apply the configuration to the 
logger. For formatters, use the `FormatterType` type key provided by the library to ensure that it 
is recognized as a formatter.

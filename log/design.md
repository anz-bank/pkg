# Design Choices

The library design has gone through several iterations.

## Iteration 1 &mdash; logrus-style, hidden in context

The first design of the library tried to hide the logger from the user as much as possible. The user would only interact with the logger directly during operations involving log-specific fields. The `Field` and `Fields` methods returned the copied `Logger` with the new fields (without storing it in the context) which they can use to add more fields through the `Logger` interface or perform a log operation through the `Logger` method. The main API can be found [here](https://github.com/anz-bank/sysl/blob/1edd0489cfff9673ab6cd4d2160e1274151af2cb/pkg/log/api.go).

```go
func WithLogger(ctx context.Context, logger loggers.Logger) context.Context {...}
// Context level Fields
func WithField(ctx context.Context, key string, val interface{}) context.Context {...}
func WithFields(ctx context.Context, fields map[string]interface{}) context.Context {...}
// Log-specific Fields
func Field(ctx context.Context, key string, value interface{}) loggers.Logger {...}
func Fields(ctx context.Context, fields map[string]interface{}) loggers.Logger {...}
func Debug(ctx context.Context, args ...interface{}) {...}
func Error(ctx context.Context, args ...interface{}) {...}
func Exit(ctx context.Context, code int) {...}
func Fatal(ctx context.Context, args ...interface{}) {...}
func Panic(ctx context.Context, args ...interface{}) {...}
func Print(ctx context.Context, args ...interface{}) {...}
func Trace(ctx context.Context, args ...interface{}) {...}
func Warn(ctx context.Context, args ...interface{}) {...}

func Debugf(ctx context.Context, format string, args ...interface{}) {...}
func Errorf(ctx context.Context, format string, args ...interface{}) {...}
⋮
```

Context-level Fields are stored in the context. Log-specific fields, which are not stored in the context, are logged in addition to context-level fields.

This design conformed closely to [logrus](https://github.com/sirupsen/logrus) in order to offer a familiar interface. However, it proved to be unnecessarily complex as the functionality overlapped with the `Logger` interface itself. This lead to a second design.

## Iteration 2 &mdash; simplified top-level API

The second design focused on cleaning up the top level APIs. It requires the user to access the `Logger`'s methods directly. The source is available [here](https://github.com/anz-bank/pkg/tree/cfb66c107eebe67ff73d9c071bdc190807295c10).

```go
func With(ctx context.Context, logger loggers.Logger) context.Context {...}

// WithFields adds multiple fields in the scope of the context, fields will be logged alphabetically
func WithFields(ctx context.Context, fields MultipleFields) {...}

// From is a way to access the API of the logger inside the context and add log-specific fields
func From(ctx context.Context, fields ...Fields) loggers.Logger {...}
```

Using the `From` API, the user then can do their log operations through the `Logger` interface which is the following:

```go
// Logger is the underlying logger that is to be added to a context
type Logger interface {
    Debug(args ...interface{})
    Debugf(format string, args ...interface{})
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    Exit(code int)
    Fatal(args ...interface{})
    Fatalf(format string, args ...interface{})
    Panic(args ...interface{})
    Panicf(format string, args ...interface{})
    Trace(args ...interface{})
    Tracef(format string, args ...interface{})
    Warn(args ...interface{})
    Warnf(format string, args ...interface{})

    // PutField returns the Logger with the new field added
    PutField(key string, val interface{}) Logger
    // PutFields returns the Logger with the new fields added
    PutFields(fields frozen.Map) Logger
    // Copy returns a logger whose data is copied from the caller
    Copy() Logger
}
```

Under this design, the `Fields` type is also preserved through the `WithFields` and `From` API. The `examples.md` file in that repo explains this well, but essentially, context level fields can be done by using the `WithFields` API and providing the fields in a `MultipleFields` data structure which is a `map[string]interface{}`. The `From` API provides an optional argument for `Fields` which is meant as a way to log with log-specific fields. `Fields` type was an interface that `MultipleFields` and the removed `SingleField` struct implements. `SingleFields` was removed because it would encourage bad usage and because of that only `MultipleFields` was allowed.

This design still had some problems. Firstly, there are too many unnecessary log levels. This [article by Dave Cheney](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) argued for far fewer levels. Secondly, the `Fields` operation is too rigid as we can not add a different type of fields and this became apparent when the concept of global fields came up. Global fields do not hold their own values. They are actually references to context values. This lead to a third design &mdash; very similar to the second design &mdash; as follows (source unavailable):

```go
func With(ctx context.Context, logger loggers.Logger) context.Context {...}

// WithFields adds multiple fields in the scope of the context, fields will be logged alphabetically
func WithFields(ctx context.Context, fields MultipleFields) context.Context {...}

// WithGlobals adds keys whose values will be taken directly from the context
func WithGlobals(ctx context.Context, keys ...interface{}) context.Context {...}

// From is a way to access the API of the logger inside the context and add log-specific fields
func From(ctx context.Context, fields ...Fields) loggers.Logger {...}

// Logger is the underlying logger that is to be added to a context
type Logger interface {
    Debug(args ...interface{})
    Debugf(format string, args ...interface{})
    Info(args ...interface{})
    Infof(format string, args ...interface{})
    PutFields(fields frozen.Map, replaceExistingValues bool) Logger
    // Copy returns a logger whose data is copied from the caller
    Copy() Logger
}
```

This design still had problems. Firstly, `PutFields` complicated the API and was strictly redundant. Secondly, with more types of `Fields` appearing, adding `Fields` became too much work as it would require adding the `Fields` to the context one by one which would create a deep context chain. `Fields` are important and as systems get bigger, the concern was that the growing number of `Fields` needed to be handled.

This led to a fourth design that treats `Fields` more consistently. We also hid some of the logger interface which makes the interface more focused on logging operations. The initial fourth redesign can be found [here](https://github.com/anz-bank/pkg/blob/9e479508561c3381a6b24165d11e548e976e69ef).

```go
func From(ctx context.Context) {...}
func Suppress(keys ...string) {...}
func With(key string, val interface{}) {...}
func WithContextRef(key string, ctxKey interface{}) {...}
func WithLogger(logger Logger) {...}

func (f Fields) Chain(fieldses ...Fields) {...}
func (f Fields) From(ctx context.Context) {...}
func (f Fields) Onto(ctx context.Context) context.Context {...}
func (f Fields) Suppress(keys ...string) Fields {...}
func (f Fields) With(key string, val interface{}) {...}
func (f Fields) WithContextRef(key string, ctxKey interface{}) {...}
func (f Fields) WithLogger(logger Logger) Fields {...}
```

This design favors chaining the operations so that users can do many things in a single statement. While chaining is not idiomatic, it does offer a very clean coding style while retaining the ability to add multiple `Fields` to the context in a single operation. This design also preserves support for special `Fields` types such as global fields via `WithContextRef`.

To simplify logging operations, the methods are also hidden in a private interface like the following.

```go
// Logger is the underlying logger that is to be added to a context
type Logger interface {
    // Debug logs the message at the Debug level
    Debug(args ...interface{})
    // Debugf logs the message at the Debug level
    Debugf(format string, args ...interface{})
    // Info logs the message at the Info level
    Info(args ...interface{})
    // Infof logs the message at the Info level
    Infof(format string, args ...interface{})
}

type internalLoggerOps interface {
    // PutFields returns the Logger with the new fields added
    PutFields(fields frozen.Map) Logger
    // Copy returns a logger whose data is copied from the caller
    Copy() Logger
}
```

This design ensures that users only have what they need while still providing flexibility over `Fields`.

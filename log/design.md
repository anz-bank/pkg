# Design Choices
The library itself has gone through many redesigns.

The first design of the library tried to hide the logger from the user as much as possible. The only time when user would interact with the logger directly is during operations involving log-specific fields. The `Field` and `Fields` return the copied `Logger` with the new fields without storing it in the context which they can use to add additional fields through the `Logger` interface or do a log operation through the `Logger` method. The main API can be found [here](https://github.com/anz-bank/sysl/blob/1edd0489cfff9673ab6cd4d2160e1274151af2cb/pkg/log/api.go).

```go
func WithLogger(ctx context.Context, logger loggers.Logger) context.Context {...}
// Context level Fields
func WithField(ctx context.Context, key string, val interface{}) context.Context {...} 
func WithFields(ctx context.Context, fields map[string]interface{}) context.Context {...} 
// Log-specific Fields
func Field(ctx context.Context, key string, value interface{}) loggers.Logger {...} 
func Fields(ctx context.Context, fields map[string]interface{}) loggers.Logger {...}
func Debug(ctx context.Context, args ...interface{}) {...} 
func Debugf(ctx context.Context, format string, args ...interface{}) {...} 
func Error(ctx context.Context, args ...interface{}) {...} 
func Errorf(ctx context.Context, format string, args ...interface{}) {...}
func Exit(ctx context.Context, code int) {...} 
func Fatal(ctx context.Context, args ...interface{}) {...}
func Fatalf(ctx context.Context, format string, args ...interface{}) {...}
func Panic(ctx context.Context, args ...interface{}) {...}
func Panicf(ctx context.Context, format string, args ...interface{}) {...}
func Print(ctx context.Context, args ...interface{}) {...}
func Printf(ctx context.Context, format string, args ...interface{}) {...}
func Trace(ctx context.Context, args ...interface{}) {...}
func Tracef(ctx context.Context, format string, args ...interface{}) {...}
func Warn(ctx context.Context, args ...interface{}) {...}
func Warnf(ctx context.Context, format string, args ...interface{}) {...}
```

Context level Fields are fields that are meant to be stored in the context and will be logged as long as the Fields in context is not changed. The log-specific fields are additional fields on top of the context level fields that are not stored in the context.

The API design at this stage was meant to copy [logrus](https://github.com/sirupsen/logrus) so that usage would feel more similar to it. But this API proved to be way too complex as the functionality overlap with the `Logger` interface itself. This would lead to a second redesign.

The second redesign focuses on cleaning up the top level APIs. This design relies on just exposing the logger library which then user can use the `Logger` method to do their log operations. You can see the repo of this design [here](https://github.com/anz-bank/pkg/tree/cfb66c107eebe67ff73d9c071bdc190807295c10).

In this design, the main API becomes like the following:
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
	// Debug logs the message at the Debug level
	Debug(args ...interface{})
	// Debugf logs the message at the Debug level
	Debugf(format string, args ...interface{})
	// Error logs the message at the Error level
	Error(args ...interface{})
	// Errorf logs the message at the Error level
	Errorf(format string, args ...interface{})
	// Exit exits the program with specified code
	Exit(code int)
	// Fatal logs the message at the Fatal level and exits the program with code 1
	Fatal(args ...interface{})
	// Fatalf logs the message at the Fatal level and exits the program with code 1
	Fatalf(format string, args ...interface{})
	// Panic logs the message at the Panic level
	Panic(args ...interface{})
	// Panicf logs the message at the Panic level
	Panicf(format string, args ...interface{})
	// Trace logs the message at the Trace level
	Trace(args ...interface{})
	// Tracef logs the message at the Trace level
	Tracef(format string, args ...interface{})
	// Warn logs the message at the Warn level
	Warn(args ...interface{})
	// Warnf logs the message at the Warn level
	Warnf(format string, args ...interface{})
	// PutField returns the Logger with the new field added
	PutField(key string, val interface{}) Logger
	// PutFields returns the Logger with the new fields added
	PutFields(fields frozen.Map) Logger
	// Copy returns a logger whose data is copied from the caller
	Copy() Logger
}
```

Through this API design, the Fields type is also preserved through the `WithFields` and `From` API. The `examples.md` file in that repo explains this well, but essentially, context level fields can be done by using the `WithFields` API and providing the fields in a `MultipleFields` data structure which is a `map[string]interface{}`. The `From` API provides an optional argument for `Fields` which is meant as a way to log with log-specific fields. `Fields` type was an interface that `MultipleFields` and the removed `SingleField` struct implements. `SingleFields` was removed because it would encourage bad usage and because of that only `MultipleFields` was allowed.

This design however presented another set of problems. Firstly, there are too many unnecessary log levels, this [article by Dave Cheney](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) was our reason. Secondly, the `Fields` operation is too rigid as we can not add a different type of fields and this became apparent when a different type of fields was needed which was the global fields. The global fields are fields whose values are taken directly from the context variable. We wanted to be able to add keys whose values are then taken from the context directly before every logs. This would lead to a third redesign which is very similar to the second redesign.

Unfortunately, the commit for the third redesign was squashed with the newest design, but the API design is the following:
```go
func With(ctx context.Context, logger loggers.Logger) context.Context {...}

// WithFields adds multiple fields in the scope of the context, fields will be logged alphabetically
func WithFields(ctx context.Context, fields MultipleFields) {...}

// WithGlobals adds keys whose values will be taken directly from the context
func WithGlobals(ctx context.Context, keys ...interface{}) context.Context {...}

// From is a way to access the API of the logger inside the context and add log-specific fields
func From(ctx context.Context, fields ...Fields) loggers.Logger {...}
``` 

And the logger interface is the following:
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
	// PutFields returns the Logger with the new fields added
	PutFields(fields frozen.Map, replaceExistingValues bool) Logger
	// Copy returns a logger whose data is copied from the caller
	Copy() Logger
}
```

This design presented a lot of problem under the hood, mainly the precedence of the fields. It was then decided that global fields would have the lowest precedence, the context level fields in the middle, and the log-specific fields have the highest.

In this design however, there are still problems. Firstly, `PutFields` became slightly more complicated which is unnecessary and since this method is exposed to the user, the complexity might confuse users. Clearly, certain methods need to be hidden. Secondly, with the more type of `Fields` appearing, adding `Fields` became too much work as it would require adding the `Fields` to the context one by one which would create a deep context tree. `Fields` are important and as systems get bigger, it became apparent that the growing number of `Fields` is something that needs to be handled.

This led to a fourth design that tries to prioritise and simplifies how `Fields` work. `Fields` priority uses a simple principle, the order of the operation dictates the `Fields` precedence where the later operations get higher precedence. We also hid some of the logger interface which makes the interface more focused on logging operations. The initial fourth redesign can be found [here](https://github.com/anz-bank/pkg/blob/9e479508561c3381a6b24165d11e548e976e69ef).

```go
func From(ctx context.Context) {...}
func Suppress(keys ...string) {...}
func With(key string, val interface{}) {...}
func WithCtxRef(key string, ctxKey interface{}) {...}
func WithFunc(key string, f func(context.Context) interface{}) {...}
func WithLogger(logger Logger) {...}
func (f Fields) Chain(fieldses ...Fields) {...}
func (f Fields) From(ctx context.Context) {...}
func (f Fields) Onto(ctx context.Context) context.Context {...}
func (f Fields) Suppress(keys ...string) Fields {...}
func (f Fields) With(key string, val interface{}) {...}
func (f Fields) WithCtxRef(key string, ctxKey interface{}) {...}
func (f Fields) WithFunc(key string, val func(context.Context) interface{}) {...}
func (f Fields) WithLogger(logger Logger) Fields {...}
```

The design favors chaining the operations so that users can do many things in one operations. While chaining is not idiomatic, it does present a very clean code and allows us to add multiple `Fields` in one operation to the context. In this model, it allows users to have better control over their `Fields`. This design also preserves the multiple types of `Fields` needed. Global fields can be added through the `WithCtxRef` API while the other `Fields` can just use the `With` APIs. Whether a fields becomes log-specific or context-specific depends on whether or not user stores them in the context with the `Onto` API.

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

This design ensures that users have what they need.
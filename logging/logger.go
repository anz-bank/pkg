package logging

import (
	"context"
	"io"
	"runtime"
	"time"

	"github.com/anz-bank/pkg/logging/codelinks"
	"github.com/arr-ai/frozen"
	"github.com/rs/zerolog"
)

// ContextFunc defines the function signature used by the logger to extract fields from context
//
// Includes an ID that loggers use to prevent duplication
type ContextFunc struct {
	Keys     []string
	Function func(ctx context.Context, event zerolog.Context) zerolog.Context
}

// Logger is responsible for executing log calls
//
// This logger is a wrapper on zerolog.Logger that adds contextual logging and
// code linking. Invoke a log with one of Info, Debug or Error, and use the
// zerolog api to add fields and execute the log.
//
// eg: logger.Info(ctx).Str("key", "value").Msg("Hello World")
//
// Add contextual logging with the With method. This takes an array of
// ContextFunc function instances to call on every log call. This should be
// done to set up the logger before it is used, usually in application init.
//
// Loggers are immutable. Mutations return a new logger with the mutation
// applied
type Logger struct {
	internal  zerolog.Logger
	keys      frozen.Set
	codeLinks bool
	linker    codelinks.CodeLinker
	timeDiffs []*timeDiff
	discard   bool
}

// New returns a new logger
func New(out io.Writer) *Logger {
	logger := zerolog.New(out).Level(zerolog.InfoLevel)
	return &Logger{
		internal:  logger,
		keys:      frozen.NewSetFromStrings("level", "message"), // prevents these fields from being duplicated
		codeLinks: false,
	}
}

// ToContext adds the logger to context
func (l Logger) ToContext(ctx context.Context) context.Context {
	return l.With().ToContext(ctx)
}

// With adds ContextFuncs to the logger
//
// Use context funcs to add any field to the logger, static or contextual.
func (l Logger) With(funcs ...ContextFunc) *Context {
	logContext := Context{
		logger: &l,
		funcs:  []ContextFunc{},
	}
	return logContext.With(funcs...)
}

// WithStr creates a static string field on the logger
func (l Logger) WithStr(key string, val string) *Logger {
	if !l.keys.Has(key) {
		l.keys = l.keys.With(key)
		l.internal = l.internal.With().Str(key, val).Logger()
	} else {
		// warn of duplicate keys
		l.Debug().Msg("duplicate keys in logger, ignoring duplication")
	}
	return &l
}

// WithInt creates a static int field on the logger
func (l Logger) WithInt(key string, val int) *Logger {
	if !l.keys.Has(key) {
		l.keys = l.keys.With(key)
		l.internal = l.internal.With().Int(key, val).Logger()
	} else {
		l.Debug().Msg("duplicate keys in logger, ignoring duplication")
	}
	return &l
}

// WithDict creates a static dictionary field (json object) on the logger
func (l Logger) WithDict(key string, val *zerolog.Event) *Logger {
	if !l.keys.Has(key) {
		l.keys = l.keys.With(key)
		l.internal = l.internal.With().Dict(key, val).Logger()
	} else {
		l.Debug().Msg("duplicate keys in logger, ignoring duplication")
	}
	return &l
}

// Dict defines the value for a dict field
//
// returns a zerolog.logCtx. Use the logCtx methods to add sub fields to a dict entry
func Dict() *zerolog.Event {
	return zerolog.Dict()
}

// WithBool creates a static boolean field on the logger
func (l Logger) WithBool(key string, val bool) *Logger {
	if !l.keys.Has(key) {
		l.keys = l.keys.With(key)
		l.internal = l.internal.With().Bool(key, val).Logger()
	} else {
		l.Debug().Msg("duplicate keys in logger, ignoring duplication")
	}
	return &l
}

// WithArray creates a static array field on the logger
func (l Logger) WithArray(key string, val zerolog.LogArrayMarshaler) *Logger {
	if !l.keys.Has(key) {
		l.keys = l.keys.With(key)
		l.internal = l.internal.With().Array(key, val).Logger()
	} else {
		l.Debug().Msg("duplicate keys in logger, ignoring duplication")
	}
	return &l
}

// Array defines the value for an array field
//
// returns a LogArrayMarshaller. Use methods on this to add array elements to an array entry
func Array() *zerolog.Array {
	return zerolog.Arr()
}

// WithDur creates a static duration field on the logger
func (l Logger) WithDur(key string, val time.Duration) *Logger {
	if !l.keys.Has(key) {
		l.keys = l.keys.With(key)
		l.internal = l.internal.With().Dur(key, val).Logger()
	} else {
		l.Debug().Msg("duplicate keys in logger, ignoring duplication")
	}
	return &l
}

// WithTime creates a static time field on the logger
func (l Logger) WithTime(key string, val time.Time) *Logger {
	if !l.keys.Has(key) {
		l.keys = l.keys.With(key)
		l.internal = l.internal.With().Time(key, val).Logger()
	} else {
		l.Debug().Msg("duplicate keys in logger, ignoring duplication")
	}
	return &l
}

// WithTimeDiff adds a time diff field to the logger
func (l Logger) WithTimeDiff(key string, start time.Time) *Logger {
	l.timeDiffs = append(l.timeDiffs, &timeDiff{
		key:   key,
		start: start,
	})
	return &l
}

// WithOutput creates a copy of the logger with a new writer to write logs to
func (l Logger) WithOutput(out io.Writer) *Logger {
	l.internal = l.internal.Output(out)
	return &l
}

// WithLevel creates a copy of the logger with a new log level
func (l Logger) WithLevel(level Level) *Logger {
	l.internal = l.internal.Level(zerolog.Level(level))
	return &l
}

// WithCodeLinks sets the code linker
//
// set 'links' to false to turn links off
func (l Logger) WithCodeLinks(links bool, linker codelinks.CodeLinker) *Logger {
	l.codeLinks = links
	l.linker = linker
	return &l
}

// Discard causes any log call on the logger to be ignored
//
// This is useful for creating a logger that shhould not log depending on some
// criteria
//
//  if cond { logger = logger.Discard() }
func (l Logger) Discard() *Logger {
	l.discard = true
	return &l
}

// Main Log functions
//
// The unexported functions allow this package to pass the call depth to withInfo
// to ensure any code links points to the caller, not some function in this library

// Info logs a message at the info level
func (l *Logger) Info() *zerolog.Event {
	return l.info(1)
}
func (l *Logger) info(depth int) *zerolog.Event {
	return l.withInfo(l.internal.Info(), depth+1)
}

// Debug logs a message at the debug level
func (l *Logger) Debug() *zerolog.Event {
	return l.debug(1)
}
func (l *Logger) debug(depth int) *zerolog.Event {
	return l.withInfo(l.internal.Debug(), depth+1)
}

// Error logs a message at the error level
func (l *Logger) Error(err error) *zerolog.Event {
	return l.error(err, 1)
}
func (l *Logger) error(err error, depth int) *zerolog.Event {
	return l.withInfo(l.internal.Err(err), depth+1)
}

// This function adds fields defined by the logger ContextFunc array and codelinks if enabled
func (l *Logger) withInfo(event *zerolog.Event, depth int) *zerolog.Event {
	if l.codeLinks {
		_, file, line, _ := runtime.Caller(depth + 1)
		event = event.Str("source_code", l.linker.Link(file, line))
	}
	for _, td := range l.timeDiffs {
		event = event.TimeDiff(td.key, time.Now(), td.start)
	}
	return event
}

type timeDiff struct {
	key   string
	start time.Time
}

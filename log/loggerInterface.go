package log

import (
	"io"
	"time"

	"github.com/arr-ai/frozen"
)

// Logger is the underlying logger that is to be added to a context.
type Logger interface {
	// Debug logs the message at the Debug level.
	Debug(args ...interface{})
	// Debugf logs the message at the Debug level.
	Debugf(format string, args ...interface{})
	// Error logs the message at the Error level
	Error(errMsg error, args ...interface{})
	// Errorf logs the message at the Error level
	Errorf(errMsg error, format string, args ...interface{})
	// Info logs the message at the Info level
	Info(args ...interface{})
	// Infof logs the message at the Info level.
	Infof(format string, args ...interface{})
}

// LogEntry describes an entry to log.
type LogEntry struct { // nolint:golint // log.LogEntry stutters but is public API now.

	// Time at which the log entry was created
	Time time.Time

	// Message passed to Debug, Info or Error
	Message string

	// Data set by the user.
	Data frozen.Map

	// Source code reference of the calling function.
	// Initialised if SetLogCaller is true.
	Caller CodeReference

	// True if the log is verbose (Debug), false otherwise (Info or Error)
	Verbose bool
}

// CodeReference describes a reference to a point within a source code file.
type CodeReference struct {

	// Path of the file (within the local file system) where the source code is found
	File string

	// Line number (1-indexed) within the source code file
	Line int
}

// Hook describes a callback to receive notice when an entry is logged
type Hook interface {
	OnLogged(*LogEntry) error
}

type copyable interface {
	// Copy returns a logger whose data is copied from the caller.
	Copy() Logger
}

type fieldSetter interface {
	// PutFields returns the Logger with the new fields added.
	PutFields(fields frozen.Map) Logger
}

type entryLogger interface {
	// Log the given log entry.
	Log(entry *LogEntry)
}

type Formattable interface {
	// SetFormatter sets the formatter for the logger.
	// The formatter provided must also implement the Formatter interface.
	SetFormatter(formatter Config) error
}

type Formatter interface {
	// Format translates a log entry into a string representation
	Format(*LogEntry) (string, error)
}

type SettableVerbosity interface {
	// SetVerbose sets the verbosity of the logger.
	SetVerbose(on bool) error
}

type SettableOutput interface {
	// SetOutput sets where the logger outputs to.
	SetOutput(w io.Writer) error
}

type AddableHooks interface {
	// AddHooks adds the given hooks to the logger.
	AddHooks(hooks ...Hook) error
}

type SettableLogCaller interface {
	// SetLogCaller sets whether or not a reference to the calling function is logged.
	SetLogCaller(on bool) error
}

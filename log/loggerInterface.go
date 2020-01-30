package log

import "github.com/arr-ai/frozen"

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

type copyable interface {
	// Copy returns a logger whose data is copied from the caller
	Copy() Logger
}

type fieldSetter interface {
	// PutFields returns the Logger with the new fields added
	PutFields(fields frozen.Map) Logger
}

type formattable interface {
	// SetFormatter sets the formatter for the logger
	SetFormatter(formatter Config) error
}

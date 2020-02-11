package log

import (
	"bytes"

	"github.com/arr-ai/frozen"
)

type nullLogger struct {
	internal *standardLogger
}

// Create a null logger that outputs to a buffer, for benchmarking
func NewNullLogger() Logger {
	return &nullLogger{internal: setUpLogger()}
}

func (n *nullLogger) Debug(args ...interface{}) {
	n.internal.Debug(args...)
}

func (n *nullLogger) Debugf(format string, args ...interface{}) {
	n.internal.Debugf(format, args...)
}

func (n *nullLogger) Info(args ...interface{}) {
	n.internal.Info(args...)
}

func (n *nullLogger) Infof(format string, args ...interface{}) {
	n.internal.Infof(format, args...)
}

func (n *nullLogger) PutFields(fields frozen.Map) Logger {
	n.internal.fields = fields
	return n
}

func (n *nullLogger) Copy() Logger {
	return &nullLogger{
		setUpLogger().PutFields(n.internal.fields).(*standardLogger),
	}
}

func (n *nullLogger) SetFormatter(formatter Config) error {
	return n.internal.SetFormatter(formatter)
}

func (n *nullLogger) SetVerbosity(on bool) error {
	return n.internal.SetVerbosity(on)
}

func setUpLogger() *standardLogger {
	logger := NewStandardLogger().(*standardLogger)
	logger.internal.SetOutput(&bytes.Buffer{})
	return logger
}

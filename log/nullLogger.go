package log

import (
	"bytes"
	"io"

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

func (n *nullLogger) Error(errMsg error, args ...interface{}) {
	n.internal.Error(errMsg, args...)
}

func (n *nullLogger) Errorf(errMsg error, format string, args ...interface{}) {
	n.internal.Errorf(errMsg, format, args...)
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

func (n *nullLogger) SetVerbose(on bool) error {
	return n.internal.SetVerbose(on)
}

func (n *nullLogger) SetOutput(w io.Writer) error {
	return n.internal.SetOutput(w)
}

func (n *nullLogger) AddHooks(hooks ...Hook) error {
	return n.internal.AddHooks(hooks...)
}

func (n *nullLogger) SetLogCaller(on bool) error {
	return n.internal.SetLogCaller(on)
}

func setUpLogger() *standardLogger {
	logger := NewStandardLogger().(*standardLogger)
	logger.internal.SetOutput(&bytes.Buffer{})
	return logger
}

package loggers

import (
	"bytes"

	"github.com/anz-bank/pkg/log"
	"github.com/arr-ai/frozen"
)

type nullLogger struct {
	internal *standardLogger
}

func NewNullLogger() log.Logger {
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

func (n *nullLogger) PutFields(fields frozen.Map) log.Logger {
	n.internal.fields = fields
	return n
}

func (n *nullLogger) Copy() log.Logger {
	return &nullLogger{
		setUpLogger().PutFields(n.internal.fields).(*standardLogger),
	}
}

func setUpLogger() *standardLogger {
	logger := NewStandardLogger().(*standardLogger)
	logger.internal.SetOutput(&bytes.Buffer{})
	return logger
}

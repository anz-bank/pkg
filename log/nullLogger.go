package log

import (
	"io"

	"github.com/arr-ai/frozen"
)

type nullLogger struct{}

// Create a null logger that doesn't log
func NewNullLogger() Logger {
	return &nullLogger{}
}

func (n *nullLogger) Debug(args ...interface{})                               {}
func (n *nullLogger) Debugf(format string, args ...interface{})               {}
func (n *nullLogger) Error(errMsg error, args ...interface{})                 {}
func (n *nullLogger) Errorf(errMsg error, format string, args ...interface{}) {}
func (n *nullLogger) Info(args ...interface{})                                {}
func (n *nullLogger) Infof(format string, args ...interface{})                {}
func (n *nullLogger) Log(entry *LogEntry)                                     {}

func (n *nullLogger) PutFields(fields frozen.Map) Logger {
	return n
}

func (n *nullLogger) Copy() Logger {
	return &nullLogger{}
}

func (n *nullLogger) SetFormatter(formatter Config) error {
	return nil
}

func (n *nullLogger) SetVerbose(on bool) error {
	return nil
}

func (n *nullLogger) SetOutput(w io.Writer) error {
	return nil
}

func (n *nullLogger) AddHooks(hooks ...Hook) error {
	return nil
}

func (n *nullLogger) SetLogCaller(on bool) error {
	return nil
}

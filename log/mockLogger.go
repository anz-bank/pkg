package log

import (
	"io"

	"github.com/arr-ai/frozen"
	"github.com/stretchr/testify/mock"
)

type mockLogger struct {
	Logger
	mock.Mock
}

func newMockLogger() *mockLogger {
	return &mockLogger{}
}

func (m *mockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *mockLogger) Debugf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *mockLogger) Error(errMsg error, args ...interface{}) {
	m.Called(append([]interface{}{errMsg}, args...)...)
}

func (m *mockLogger) Errorf(errMsg error, format string, args ...interface{}) {
	m.Called(append([]interface{}{errMsg, format}, args...)...)
}

func (m *mockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *mockLogger) Infof(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *mockLogger) Log(entry *LogEntry) {
	m.Called(entry)
}

func (m *mockLogger) PutFields(fields frozen.Map) Logger {
	return m.Called(fields).Get(0).(Logger)
}

func (m *mockLogger) Copy() Logger {
	return m.Called().Get(0).(Logger)
}

func (m *mockLogger) SetFormatter(formatter Config) error {
	return m.Called(formatter).Error(0)
}

func (m *mockLogger) SetVerbose(on bool) error {
	return m.Called(on).Error(0)
}

func (m *mockLogger) SetOutput(w io.Writer) error {
	return m.Called(w).Error(0)
}

func (m *mockLogger) AddHooks(hooks ...Hook) error {
	return m.Called(hooks).Error(0)
}

func (m *mockLogger) SetLogCaller(on bool) error {
	return m.Called(on).Error(0)
}

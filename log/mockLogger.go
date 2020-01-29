package log

import (
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

func (m *mockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *mockLogger) Infof(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *mockLogger) PutFields(fields frozen.Map) Logger {
	return m.Called(fields).Get(0).(Logger)
}

func (m *mockLogger) Copy() Logger {
	return m.Called().Get(0).(Logger)
}

func (m *mockLogger) SetFormatter(formatter Config) error {
	res := m.Called(formatter)
	return res.Error(0)
}

package loggers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anz-bank/pkg/log"
	"github.com/arr-ai/frozen"
	"github.com/sirupsen/logrus"
)

const keyFields = "_fields"

type standardLogger struct {
	internal *logrus.Logger
	fields   frozen.Map
}

type standardFormat struct{}

func (sf *standardFormat) Format(entry *logrus.Entry) ([]byte, error) {
	message := strings.Builder{}
	message.WriteString(entry.Time.Format(time.RFC3339Nano))
	message.WriteByte(' ')

	if entry.Data[keyFields] != "" {
		message.WriteString(entry.Data[keyFields].(string))
		message.WriteByte(' ')
	}

	message.WriteString(strings.ToUpper(entry.Level.String()))
	message.WriteByte(' ')

	if entry.Message != "" {
		message.WriteString(entry.Message)
		message.WriteByte(' ')
	}

	// TODO: add codelinker's message here
	message.WriteByte('\n')
	return []byte(message.String()), nil
}

// NewStandardLogger returns a logger with logrus standard logger as the internal logger
func NewStandardLogger() log.Logger {
	return &standardLogger{internal: setupStandardLogger()}
}

func (sl *standardLogger) Debug(args ...interface{}) {
	sl.setInfo().Debug(args...)
}

func (sl *standardLogger) Debugf(format string, args ...interface{}) {
	sl.setInfo().Debugf(format, args...)
}

func (sl *standardLogger) Info(args ...interface{}) {
	sl.setInfo().Info(args...)
}

func (sl *standardLogger) Infof(format string, args ...interface{}) {
	sl.setInfo().Infof(format, args...)
}

func (sl *standardLogger) PutFields(fields frozen.Map) log.Logger {
	sl.fields = fields
	return sl
}

func (sl *standardLogger) Copy() log.Logger {
	return &standardLogger{setupStandardLogger(), sl.fields}
}

func (sl *standardLogger) setInfo() *logrus.Entry {
	// TODO: set linker here
	return sl.internal.WithFields(logrus.Fields{
		keyFields: sl.getFormattedField(),
	})
}

func (sl *standardLogger) getFormattedField() string {
	if sl.fields.Count() == 0 {
		return ""
	}

	fields := strings.Builder{}
	i := sl.fields.Range()
	i.Next()
	fields.WriteString(fmt.Sprintf("%v=%v", i.Key(), i.Value()))
	for i.Next() {
		fields.WriteString(fmt.Sprintf(" %v=%v", i.Key(), i.Value()))
	}
	return fields.String()
}

func setupStandardLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&standardFormat{})

	// makes sure that it always logs every level
	logger.SetLevel(logrus.DebugLevel)

	// explicitly set it to os.Stderr
	logger.SetOutput(os.Stderr)

	return logger
}

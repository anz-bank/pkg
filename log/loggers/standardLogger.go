package loggers

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/marcelocantos/frozen"
	"github.com/sirupsen/logrus"
)

const keyFields = "_fields"

type standardLogger struct {
	fields       frozen.Map
	sortedFields []string
	internal     *logrus.Logger
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
func NewStandardLogger() Logger {
	logger := setupStandardLogger()

	return &standardLogger{
		fields:       frozen.NewMap(),
		sortedFields: []string{},
		internal:     logger,
	}
}

func (sl *standardLogger) Debug(args ...interface{}) {
	sl.setInfo().Debug(args...)
}

func (sl *standardLogger) Debugf(format string, args ...interface{}) {
	sl.setInfo().Debugf(format, args...)
}

func (sl *standardLogger) Error(args ...interface{}) {
	sl.setInfo().Error(args...)
}

func (sl *standardLogger) Errorf(format string, args ...interface{}) {
	sl.setInfo().Errorf(format, args...)
}

func (sl *standardLogger) Exit(code int) {
	sl.internal.Exit(code)
}

func (sl *standardLogger) Fatal(args ...interface{}) {
	sl.setInfo().Fatal(args...)
}

func (sl *standardLogger) Fatalf(format string, args ...interface{}) {
	sl.setInfo().Fatalf(format, args...)
}

func (sl *standardLogger) Panic(args ...interface{}) {
	sl.setInfo().Panic(args...)
}

func (sl *standardLogger) Panicf(format string, args ...interface{}) {
	sl.setInfo().Panicf(format, args...)
}

func (sl *standardLogger) Trace(args ...interface{}) {
	sl.setInfo().Trace(args...)
}

func (sl *standardLogger) Tracef(format string, args ...interface{}) {
	sl.setInfo().Tracef(format, args...)
}

func (sl *standardLogger) Warn(args ...interface{}) {
	sl.setInfo().Warn(args...)
}

func (sl *standardLogger) Warnf(format string, args ...interface{}) {
	sl.setInfo().Warnf(format, args...)
}

func (sl *standardLogger) PutField(key string, val interface{}) Logger {
	if !sl.fields.Has(key) {
		sl.insertFieldsKey(key)
	}
	sl.fields = sl.fields.With(key, val)
	return sl
}

func (sl *standardLogger) PutFields(fields frozen.Map) Logger {
	if fields.Count() == 0 {
		panic("fields can not be empty")
	}

	keys := make([]string, fields.Count())
	index := 0
	for i := fields.Keys().Range(); i.Next(); {
		if !sl.fields.Has(i.Value()) {
			keys[index] = i.Value().(string)
			index++
		}
		sl.fields = sl.fields.With(i.Value(), fields.MustGet(i.Value()))
	}
	if index > 0 {
		sl.insertFieldsKey(keys[:index]...)
	}
	return sl
}

func (sl *standardLogger) insertFieldsKey(fields ...string) {
	newFields := append(sl.sortedFields, fields...)
	sort.Strings(newFields)
	sl.sortedFields = newFields
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
	fields.WriteString(fmt.Sprintf("%s=%v", sl.sortedFields[0], sl.fields.MustGet(sl.sortedFields[0])))
	if sl.fields.Count() > 1 {
		for _, field := range sl.sortedFields[1:] {
			fields.WriteString(fmt.Sprintf(" %s=%v", field, sl.fields.MustGet(field)))
		}
	}
	return fields.String()
}

func (sl *standardLogger) Copy() Logger {
	sortedFields := make([]string, sl.fields.Count())
	copy(sortedFields, sl.sortedFields)

	return &standardLogger{
		fields:       frozen.NewMapFromKeys(sl.fields.Keys(), sl.fields.MustGet),
		internal:     setupStandardLogger(),
		sortedFields: sortedFields,
	}
}

func setupStandardLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&standardFormat{})

	// makes sure that it always logs every level
	logger.SetLevel(logrus.TraceLevel)

	// explicitly set it to os.Stderr
	logger.SetOutput(os.Stderr)

	return logger
}

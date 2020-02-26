package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arr-ai/frozen"
	"github.com/sirupsen/logrus"
)

const keyFields = "_fields"

type standardLogger struct {
	internal *logrus.Logger
	fields   frozen.Map
}

func (sf standardFormat) Format(entry *logrus.Entry) ([]byte, error) {
	sections := append(make([]string, 0, 5), entry.Time.Format(time.RFC3339Nano), strings.ToUpper(entry.Level.String()))

	if entry.Data[keyFields] != nil && entry.Data[keyFields].(frozen.Map).Count() != 0 {
		sections = append(sections, getFormattedField(entry.Data[keyFields].(frozen.Map)))
	}

	if entry.Message != "" {
		sections = append(sections, entry.Message)
	}

	// TODO: add codelinker's message here
	sections = append(sections, "\n")
	return []byte(strings.Join(sections, " ")), nil
}

func (jf jsonFormat) Format(entry *logrus.Entry) ([]byte, error) {
	jsonFile := make(map[string]interface{})
	jsonFile["timestamp"] = entry.Time.Format(time.RFC3339Nano)
	jsonFile["message"] = entry.Message
	jsonFile["level"] = strings.ToUpper(entry.Level.String())
	if entry.Data[keyFields] != nil && entry.Data[keyFields].(frozen.Map).Count() != 0 {
		fields := make(map[string]interface{})
		for i := entry.Data[keyFields].(frozen.Map).Range(); i.Next(); {
			fields[i.Key().(string)] = i.Value()
		}
		jsonFile["fields"] = fields
	}
	data, err := json.Marshal(jsonFile)
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), err
}

// NewStandardLogger returns a logger with logrus standard logger as the internal logger
func NewStandardLogger() Logger {
	logger := logrus.New()
	logger.SetFormatter(&standardFormat{})

	return &standardLogger{internal: logger}
}

func (sl *standardLogger) Debug(args ...interface{}) {
	sl.setInfo().Debug(args...)
}

func (sl *standardLogger) Debugf(format string, args ...interface{}) {
	sl.setInfo().Debugf(format, args...)
}

func (sl *standardLogger) Error(errMsg error, args ...interface{}) {
	if msg, _ := sl.fields.Get(errMsgKey); msg != errMsg.Error() {
		sl.fields = sl.fields.With(errMsgKey, errMsg.Error())
	}
	sl.setInfo().Info(args...)
}

func (sl *standardLogger) Errorf(errMsg error, format string, args ...interface{}) {
	if msg, _ := sl.fields.Get(errMsgKey); msg != errMsg.Error() {
		sl.fields = sl.fields.With(errMsgKey, errMsg.Error())
	}
	sl.setInfo().Infof(format, args...)
}

func (sl *standardLogger) Info(args ...interface{}) {
	sl.setInfo().Info(args...)
}

func (sl *standardLogger) Infof(format string, args ...interface{}) {
	sl.setInfo().Infof(format, args...)
}

func (sl *standardLogger) PutFields(fields frozen.Map) Logger {
	sl.fields = fields
	return sl
}

func (sl *standardLogger) SetFormatter(formatter Config) error {
	logrusFormatter, isLogrusFormatter := formatter.(logrus.Formatter)
	if !isLogrusFormatter {
		return errors.New("formatter is not logrus formatter type")
	}
	sl.internal.SetFormatter(logrusFormatter)
	return nil
}

func (sl *standardLogger) SetVerbose(on bool) error {
	if on {
		sl.internal.SetLevel(logrus.DebugLevel)
	} else {
		sl.internal.SetLevel(logrus.InfoLevel)
	}
	return nil
}

func (sl *standardLogger) Copy() Logger {
	return &standardLogger{sl.getCopiedInternalLogger(), sl.fields}
}

func (sl *standardLogger) setInfo() *logrus.Entry {
	// TODO: set linker here
	return sl.internal.WithFields(logrus.Fields{
		keyFields: sl.fields,
	})
}

func getFormattedField(fields frozen.Map) string {
	if fields.Count() == 0 {
		return ""
	}

	formattedFields := strings.Builder{}
	i := fields.Range()
	i.Next()
	formattedFields.WriteString(fmt.Sprintf("%v=%v", i.Key(), i.Value()))
	for i.Next() {
		formattedFields.WriteString(fmt.Sprintf(" %v=%v", i.Key(), i.Value()))
	}
	return formattedFields.String()
}

func (sl *standardLogger) getCopiedInternalLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(sl.internal.Formatter)
	logger.SetLevel(sl.internal.Level)
	logger.SetOutput(sl.internal.Out)

	return logger
}

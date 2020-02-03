package log

//TODO: decide whether this should just be part of null logger, remember null logger is used for benchmark and merging
// this with that will change benchmark value
import (
	"fmt"
	"time"

	"github.com/arr-ai/frozen"
	"github.com/sirupsen/logrus"
)

const (
	contains int = iota
	equal
)

type TestLogger struct {
	logs []logMessage
	//TODO: mock style expectation or just assert at the end?
	//expectations []expectation
	fields frozen.Map
	format string
}

type logMessage struct {
	Message, Time, Level, Format string
	//TODO: decide whether to output frozen.Map or map[string]interface{}
	Fields frozen.Map
}

//type expectation struct {
//}

// Create a null logger that outputs to a buffer, for benchmarking
func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (tl *TestLogger) Debug(args ...interface{}) {
	tl.logs = append(tl.logs, logMessage{
		Message: fmt.Sprint(args...),
		Time:    logrus.Entry{}.Time.Format(time.RFC3339Nano),
		Level:   logrus.DebugLevel.String(),
		Format:  tl.format,
		Fields:  tl.fields,
	})
}

func (tl *TestLogger) Debugf(Format string, args ...interface{}) {
	tl.logs = append(tl.logs, logMessage{
		Message: fmt.Sprintf(Format, args...),
		Time:    logrus.Entry{}.Time.Format(time.RFC3339Nano),
		Level:   logrus.DebugLevel.String(),
		Format:  tl.format,
		Fields:  tl.fields,
	})
}

func (tl *TestLogger) Info(args ...interface{}) {
	tl.logs = append(tl.logs, logMessage{
		Message: fmt.Sprint(args...),
		Time:    logrus.Entry{}.Time.Format(time.RFC3339Nano),
		Level:   logrus.InfoLevel.String(),
		Format:  tl.format,
		Fields:  tl.fields,
	})
}

func (tl *TestLogger) Infof(Format string, args ...interface{}) {
	tl.logs = append(tl.logs, logMessage{
		Message: fmt.Sprintf(Format, args...),
		Time:    logrus.Entry{}.Time.Format(time.RFC3339Nano),
		Level:   logrus.InfoLevel.String(),
		Format:  tl.format,
		Fields:  tl.fields,
	})
}

func (tl *TestLogger) PutFields(Fields frozen.Map) Logger {
	tl.fields = Fields
	return tl
}

func (tl *TestLogger) Copy() Logger {
	return &TestLogger{tl.logs[:], tl.fields, tl.format}
}

func (tl *TestLogger) SetFormatter(Formatter Config) error {
	//TODO: need to have a talk about this, maybe Config should report its own name since it's a test for custom Formatter
	switch Formatter.(type) {
	case jsonFormat:
		tl.format = "json"
	case standardFormat:
		tl.format = "standard"
	default:
		tl.format = "unrecognizable"
	}
	return nil
}

func (tl *TestLogger) GetLogs() []logMessage {
	return tl.logs
}

func (tl *TestLogger) GetLastLog() logMessage {
	if len(tl.logs) == 0 {
		return logMessage{}
	}
	return tl.logs[len(tl.logs)-1]
}

func (tl *TestLogger) Reset() *TestLogger {
	*tl = *NewTestLogger()
	return tl
}

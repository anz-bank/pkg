package log

import (
	"github.com/arr-ai/frozen"
	"github.com/sirupsen/logrus"
)

const dataFieldKey = "_data"
const dataFieldCaller = "_caller"

// Log the given entry with the logrus logger.
func logWithLogrus(logger *logrus.Logger, e *LogEntry) {
	pkgLogEntryToLogrusEntry(logger, e).Log(verboseToLogrusLevel(e.Verbose), e.Message)
}

// Component to convert a logrus formatter into a pkg formatter.
type logrusFormatterToPkgFormatter struct {
	logger    *logrus.Logger
	formatter logrus.Formatter
}

func (f logrusFormatterToPkgFormatter) Format(entry *LogEntry) (string, error) {
	format, err := f.formatter.Format(pkgLogEntryToLogrusEntry(f.logger, entry))
	if err != nil {
		return "", err
	}
	return string(format), nil
}

// Component to convert a pkg formatter into a logrus formatter.
type pkgFormatterToLogrusFormatter struct {
	formatter Formatter
}

func (f pkgFormatterToLogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	format, err := f.formatter.Format(logrusEntryToPkgLogEntry(entry))
	if err != nil {
		return nil, err
	}
	return []byte(format), nil
}

// Convert the given pkg entry to a logrus entry.
func pkgLogEntryToLogrusEntry(logger *logrus.Logger, entry *LogEntry) *logrus.Entry {
	return &logrus.Entry{
		Logger:  logger,
		Data:    logrus.Fields{dataFieldKey: entry.Data, dataFieldCaller: entry.Caller},
		Time:    entry.Time,
		Level:   verboseToLogrusLevel(entry.Verbose),
		Caller:  nil,
		Message: entry.Message,
		Buffer:  nil,
		Context: nil,
	}
}

// Convert the given pkg entry to a logrus entry.
func logrusEntryToPkgLogEntry(entry *logrus.Entry) *LogEntry {
	return &LogEntry{
		Time:    entry.Time,
		Message: entry.Message,
		Data:    entry.Data[dataFieldKey].(frozen.Map),
		Caller:  entry.Data[dataFieldCaller].(CodeReference),
		Verbose: logrusLevelToVerbose(entry.Level),
	}
}

// Component to convert a pkg hook into a logrus hook.
type pkgHookToLogrusHook struct {
	hook Hook
}

func (h pkgHookToLogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h pkgHookToLogrusHook) Fire(entry *logrus.Entry) error {
	return h.hook.OnLogged(logrusEntryToPkgLogEntry(entry))
}

// Convert the pkg concept of verbosity to a logrus level.
func verboseToLogrusLevel(verbose bool) logrus.Level {
	if verbose {
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}

// Convert a logrus level to the pkg concept of verbosity.
func logrusLevelToVerbose(level logrus.Level) bool {
	return level == logrus.DebugLevel
}

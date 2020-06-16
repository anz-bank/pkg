package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/arr-ai/frozen"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testMessage  = "This is a test message"
	simpleFormat = "%s"
)

var errTest = errors.New("this is an error")

// to test fields output for all log
var testField = generateMultipleFieldsCases()[0].Fields

type recordHook struct {
	entries []*LogEntry
}

func (h *recordHook) OnLogged(entry *LogEntry) error {
	h.entries = append(h.entries, entry)
	return nil
}

func TestCopyStandardLogger(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger().PutFields(
		frozen.NewMap(
			frozen.KV("numberVal", 1),
			frozen.KV("byteVal", 'k'),
			frozen.KV("stringVal", "this is a sentence"),
		),
	).(*standardLogger)
	copiedLogger := logger.Copy().(*standardLogger)
	assert.NotEqual(t, logger.internal, copiedLogger.internal)
	assert.True(t, logger.fields.Equal(copiedLogger.fields))
	assert.True(t, logger != copiedLogger)
}

func TestDebug(t *testing.T) {
	testStandardLogOutput(t, logrus.DebugLevel, frozen.NewMap(), func() {
		getNewStandardLogger().Debug(testMessage)
	})

	testJSONLogOutput(t, logrus.DebugLevel, frozen.NewMap(), func() {
		logger := getNewStandardLogger()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Debug(testMessage)
	})

	testStandardLogOutput(t, logrus.DebugLevel, testField, func() {
		getStandardLoggerWithFields().Debug(testMessage)
	})

	testJSONLogOutput(t, logrus.DebugLevel, testField, func() {
		logger := getStandardLoggerWithFields()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Debug(testMessage)
	})
}

func TestDebugf(t *testing.T) {
	testStandardLogOutput(t, logrus.DebugLevel, frozen.NewMap(), func() {
		getNewStandardLogger().Debugf(simpleFormat, testMessage)
	})

	testJSONLogOutput(t, logrus.DebugLevel, frozen.NewMap(), func() {
		logger := getNewStandardLogger()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Debugf(simpleFormat, testMessage)
	})

	testStandardLogOutput(t, logrus.DebugLevel, testField, func() {
		getStandardLoggerWithFields().Debugf(simpleFormat, testMessage)
	})

	testJSONLogOutput(t, logrus.DebugLevel, testField, func() {
		logger := getStandardLoggerWithFields()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Debugf(simpleFormat, testMessage)
	})
}

func TestError(t *testing.T) {
	testStandardLogOutput(t, logrus.InfoLevel, frozen.NewMap().With(errMsgKey, errTest.Error()), func() {
		NewStandardLogger().Error(errTest, testMessage)
	})

	testJSONLogOutput(t, logrus.InfoLevel, frozen.NewMap().With(errMsgKey, errTest.Error()), func() {
		logger := getNewStandardLogger()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Error(errTest, testMessage)
	})

	testStandardLogOutput(t, logrus.InfoLevel, testField.With(errMsgKey, errTest.Error()), func() {
		getStandardLoggerWithFields().Error(errTest, testMessage)
	})

	testJSONLogOutput(t, logrus.InfoLevel, testField.With(errMsgKey, errTest.Error()), func() {
		logger := getStandardLoggerWithFields()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Error(errTest, testMessage)
	})
}

func TestErrorf(t *testing.T) {
	testStandardLogOutput(t, logrus.InfoLevel, frozen.NewMap().With(errMsgKey, errTest.Error()), func() {
		NewStandardLogger().Errorf(errTest, simpleFormat, testMessage)
	})

	testJSONLogOutput(t, logrus.InfoLevel, frozen.NewMap().With(errMsgKey, errTest.Error()), func() {
		logger := getNewStandardLogger()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Errorf(errTest, simpleFormat, testMessage)
	})

	testStandardLogOutput(t, logrus.InfoLevel, testField.With(errMsgKey, errTest.Error()), func() {
		getStandardLoggerWithFields().Errorf(errTest, simpleFormat, testMessage)
	})

	testJSONLogOutput(t, logrus.InfoLevel, testField.With(errMsgKey, errTest.Error()), func() {
		logger := getStandardLoggerWithFields()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Errorf(errTest, simpleFormat, testMessage)
	})
}

func TestInfo(t *testing.T) {
	testStandardLogOutput(t, logrus.InfoLevel, frozen.NewMap(), func() {
		getNewStandardLogger().Info(testMessage)
	})

	testJSONLogOutput(t, logrus.InfoLevel, frozen.NewMap(), func() {
		logger := getNewStandardLogger()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Info(testMessage)
	})

	testStandardLogOutput(t, logrus.InfoLevel, testField, func() {
		getStandardLoggerWithFields().Info(testMessage)
	})

	testJSONLogOutput(t, logrus.InfoLevel, testField, func() {
		logger := getStandardLoggerWithFields()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Info(testMessage)
	})
}

func TestInfof(t *testing.T) {
	testStandardLogOutput(t, logrus.InfoLevel, frozen.NewMap(), func() {
		getNewStandardLogger().Infof(simpleFormat, testMessage)
	})

	testJSONLogOutput(t, logrus.InfoLevel, frozen.NewMap(), func() {
		logger := getNewStandardLogger()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Infof(simpleFormat, testMessage)
	})

	testStandardLogOutput(t, logrus.InfoLevel, testField, func() {
		getStandardLoggerWithFields().Infof(simpleFormat, testMessage)
	})

	testJSONLogOutput(t, logrus.InfoLevel, testField, func() {
		logger := getStandardLoggerWithFields()
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Infof(simpleFormat, testMessage)
	})
}

func testStandardLogOutput(t *testing.T, level logrus.Level, fields frozen.Map, logFunc func()) {
	expectedOutput := strings.Join([]string{strings.ToUpper(level.String()), testMessage}, " ")
	actualOutput := redirectOutput(t, logFunc)

	// uses Contains to avoid checking timestamps
	assert.Contains(t, actualOutput, expectedOutput)
	for i := fields.Range(); i.Next(); {
		assert.Contains(t, actualOutput, fmt.Sprintf("%s=%v", i.Key(), i.Value()))
	}
}

func testJSONLogOutput(t *testing.T, level logrus.Level, fields frozen.Map, logFunc func()) {
	out := make(map[string]interface{})
	require.NoError(t, json.Unmarshal([]byte(redirectOutput(t, logFunc)), &out))
	assert.Equal(t, out["message"], testMessage)
	assert.Equal(t, out["level"], strings.ToUpper(level.String()))
	if fields.Count() != 0 {
		// type correction because json unmarshall reads numbers as float64
		if fields.Has("byte") && fields.Has("int") {
			fields = fields.With("byte", float64('1')).With("int", float64(123))
		}
		assert.Equal(t,
			convertToGoMap(fields),
			out["fields"].(map[string]interface{}),
		)
	}
}

func TestNewStandardLogger(t *testing.T) {
	t.Parallel()

	logger := NewStandardLogger()

	require.NotNil(t, logger)
	assert.IsType(t, logger, &standardLogger{})
}

func TestGetFormattedFieldEmptyFields(t *testing.T) {
	t.Parallel()

	require.Equal(t, "", getFormattedField(getNewStandardLogger().fields))
}

func TestGetFormattedFieldWithFields(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger().PutFields(
		frozen.NewMap(
			frozen.KV("numberVal", 1),
			frozen.KV("byteVal", 'k'),
			frozen.KV("stringVal", "this is a sentence"),
		),
	).(*standardLogger)
	// fields are in a random order
	expectedFields := []string{"byteVal=107", "numberVal=1", "stringVal=this is a sentence"}
	actualFields := getFormattedField(logger.fields)
	for _, e := range expectedFields {
		assert.Contains(t, actualFields, e)
	}
}

func TestPutFields(t *testing.T) {
	t.Parallel()

	cases := generateMultipleFieldsCases()
	for _, c := range cases {
		c := c
		t.Run(c.Name,
			func(t *testing.T) {
				t.Parallel()

				logger := getNewStandardLogger()
				logger.PutFields(c.Fields)
				assert.True(t, c.Fields.Equal(logger.fields))
			})
	}
}

func TestAddHooks(t *testing.T) {
	hook := recordHook{}
	logger := getNewStandardLogger()
	require.NoError(t, logger.SetLogCaller(true))
	require.NoError(t, logger.AddHooks(&hook))
	logger.Info("info")
	logger.Debug("debug")
	logger.Error(errors.New("error"), "error")
	assert.Equal(t, 3, len(hook.entries))
	assert.Equal(t, "info", hook.entries[0].Message)
	assert.Equal(t, "debug", hook.entries[1].Message)
	assert.Equal(t, "error", hook.entries[2].Message)
}

func TestAddHooksInfoLevel(t *testing.T) {
	hook := recordHook{}
	logger := getNewStandardLogger()
	logger.internal.SetLevel(logrus.InfoLevel)
	require.NoError(t, logger.SetLogCaller(true))
	require.NoError(t, logger.AddHooks(&hook))
	logger.Info("info")
	logger.Debug("debug") // should not be received by hook
	logger.Error(errors.New("error"), "error")
	assert.Equal(t, 2, len(hook.entries))
	assert.Equal(t, "info", hook.entries[0].Message)
	assert.Equal(t, "error", hook.entries[1].Message)
}

func TestLogCaller(t *testing.T) {
	// test standard logger
	actualOutput := redirectOutput(t, func() {
		logger := getNewStandardLogger()
		require.NoError(t, logger.SetLogCaller(true))
		logger.Debug(testMessage)
	})
	assert.Regexp(t, regexp.MustCompile(`.*\[.*standardLogger_test.go:\d+]`), actualOutput)

	// test json logger
	out := make(map[string]interface{})
	require.NoError(t, json.Unmarshal([]byte(redirectOutput(t, func() {
		logger := getNewStandardLogger()
		require.NoError(t, logger.SetLogCaller(true))
		require.NoError(t, logger.SetFormatter(NewJSONFormat()))
		logger.Debug(testMessage)
	})), &out))
	assert.Equal(t, out["message"], testMessage)
	assert.Regexp(t, regexp.MustCompile(`.*standardLogger_test.go:\d+`), out["caller"])
}

func getNewStandardLogger() *standardLogger {
	l := NewStandardLogger().(*standardLogger)
	l.internal.SetLevel(logrus.DebugLevel)
	return l
}

func getStandardLoggerWithFields() *standardLogger {
	logger := getNewStandardLogger().PutFields(testField).(*standardLogger)
	logger.internal.SetLevel(logrus.DebugLevel)
	return logger
}

func TestStandardLogger(t *testing.T) {
	logger := getNewStandardLogger()
	buffer := bytes.Buffer{}
	require.NoError(t, logger.SetOutput(&buffer))
	require.NoError(t, logger.SetVerbose(true))
	logger.Info("info")
	require.True(t, strings.Contains(buffer.String(), "info"))
	require.False(t, strings.Contains(buffer.String(), "standardLogger_test.go")) //don't log caller

	//set log caller
	buffer.Reset()
	require.NoError(t, logger.SetLogCaller(true))
	logger.Info("info")

	require.True(t, strings.Contains(buffer.String(), "info"))
	require.True(t, strings.Contains(buffer.String(), "standardLogger_test.go")) //log caller
}

func TestStandardLoggerWithFields(t *testing.T) {
	logger := getStandardLoggerWithFields()
	buffer := bytes.Buffer{}
	require.NoError(t, logger.SetOutput(&buffer))
	require.NoError(t, logger.SetVerbose(false))

	logger.Info("info")
	require.True(t, strings.Contains(buffer.String(), "info"))
	require.True(t, strings.Contains(buffer.String(), "string=this is an unnecessarily long sentence"))
	require.False(t, strings.Contains(buffer.String(), "standardLogger_test.go")) //don't log caller

	//set log caller
	buffer.Reset()
	require.NoError(t, logger.SetLogCaller(true))
	logger.Info("info")

	require.True(t, strings.Contains(buffer.String(), "info"))
	require.True(t, strings.Contains(buffer.String(), "string=this is an unnecessarily long sentence"))
	require.True(t, strings.Contains(buffer.String(), "standardLogger_test.go")) //log caller
}

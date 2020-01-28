package log

import (
	"context"
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

// to test fields output for all log
var testField = generateMultipleFieldsCases()[0].Fields

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
		NewStandardLogger().Debug(testMessage)
	})

	testStandardLogOutput(t, logrus.DebugLevel, testField, func() {
		getStandardLoggerWithFields().Debug(testMessage)
	})
}

func TestDebugf(t *testing.T) {
	testStandardLogOutput(t, logrus.DebugLevel, frozen.NewMap(), func() {
		NewStandardLogger().Debugf(simpleFormat, testMessage)
	})

	testStandardLogOutput(t, logrus.DebugLevel, testField, func() {
		getStandardLoggerWithFields().Debugf(simpleFormat, testMessage)
	})
}

func TestInfo(t *testing.T) {
	testStandardLogOutput(t, logrus.InfoLevel, frozen.NewMap(), func() {
		NewStandardLogger().Info(testMessage)
	})

	testStandardLogOutput(t, logrus.InfoLevel, testField, func() {
		getStandardLoggerWithFields().Info(testMessage)
	})
}

func TestInfof(t *testing.T) {
	testStandardLogOutput(t, logrus.InfoLevel, frozen.NewMap(), func() {
		NewStandardLogger().Infof(simpleFormat, testMessage)
	})

	testStandardLogOutput(t, logrus.InfoLevel, testField, func() {
		getStandardLoggerWithFields().Infof(simpleFormat, testMessage)
	})
}

func testStandardLogOutput(t *testing.T, level logrus.Level, fields frozen.Map, logFunc func()) {
	expectedOutput := strings.Join([]string{strings.ToUpper(level.String()), testMessage}, " ")
	actualOutput := redirectOutput(t, logFunc)

	// uses Contains to avoid checking timestamps and fields
	assert.Contains(t, actualOutput, expectedOutput)
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

func TestJsonFormat(t *testing.T) {
	WithLogger(NewStandardLogger()).WithConfigs(JSONFormatter{}).From(context.Background()).Info("Hello there")
}

func getNewStandardLogger() *standardLogger {
	return NewStandardLogger().(*standardLogger)
}

func getStandardLoggerWithFields() *standardLogger {
	logger := getNewStandardLogger().PutFields(testField)
	return logger.(*standardLogger)
}

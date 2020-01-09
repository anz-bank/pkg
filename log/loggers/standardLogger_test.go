package loggers

import (
	"sort"
	"strings"
	"testing"

	"github.com/anz-bank/pkg/log/testutil"
	"github.com/marcelocantos/frozen"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testMessage  = "This is a test message"
	simpleFormat = "%s"
)

// to test fields output for all log
var testField = testutil.GenerateMultipleFieldsCases()[0].Fields

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
	assert.Equal(t, logger.sortedFields, copiedLogger.sortedFields)
	assert.True(t, sort.StringsAreSorted(logger.sortedFields))
	assert.True(t, logger != copiedLogger)
}

func TestDebug(t *testing.T) {
	testLogOutput(t, logrus.DebugLevel, frozen.NewMap(), func() {
		NewStandardLogger().Debug(testMessage)
	})

	testLogOutput(t, logrus.DebugLevel, testField, func() {
		getStandardLoggerWithFields().Debug(testMessage)
	})
}

func TestWarn(t *testing.T) {
	testLogOutput(t, logrus.WarnLevel, frozen.NewMap(), func() {
		NewStandardLogger().Warn(testMessage)
	})

	testLogOutput(t, logrus.WarnLevel, testField, func() {
		getStandardLoggerWithFields().Warn(testMessage)
	})
}

func TestTrace(t *testing.T) {
	testLogOutput(t, logrus.TraceLevel, frozen.NewMap(), func() {
		NewStandardLogger().Trace(testMessage)
	})

	testLogOutput(t, logrus.TraceLevel, testField, func() {
		getStandardLoggerWithFields().Trace(testMessage)
	})
}

func TestPanic(t *testing.T) {
	testLogOutput(t, logrus.PanicLevel, frozen.NewMap(), func() {
		require.Panics(t, func() {
			NewStandardLogger().Panic(testMessage)
		})
	})

	testLogOutput(t, logrus.PanicLevel, frozen.NewMap(), func() {
		require.Panics(t, func() {
			getStandardLoggerWithFields().Panic(testMessage)
		})
	})
}

func TestError(t *testing.T) {
	testLogOutput(t, logrus.ErrorLevel, frozen.NewMap(), func() {
		NewStandardLogger().Error(testMessage)
	})

	testLogOutput(t, logrus.ErrorLevel, testField, func() {
		getStandardLoggerWithFields().Error(testMessage)
	})
}

func TestErrorf(t *testing.T) {
	testLogOutput(t, logrus.ErrorLevel, frozen.NewMap(), func() {
		NewStandardLogger().Errorf(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.ErrorLevel, testField, func() {
		getStandardLoggerWithFields().Errorf(simpleFormat, testMessage)
	})
}

func TestDebugf(t *testing.T) {
	testLogOutput(t, logrus.DebugLevel, frozen.NewMap(), func() {
		NewStandardLogger().Debugf(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.DebugLevel, testField, func() {
		getStandardLoggerWithFields().Debugf(simpleFormat, testMessage)
	})
}

func TestWarnf(t *testing.T) {
	testLogOutput(t, logrus.WarnLevel, frozen.NewMap(), func() {
		NewStandardLogger().Warnf(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.WarnLevel, testField, func() {
		getStandardLoggerWithFields().Warnf(simpleFormat, testMessage)
	})
}

func TestTracef(t *testing.T) {
	testLogOutput(t, logrus.TraceLevel, frozen.NewMap(), func() {
		NewStandardLogger().Tracef(simpleFormat, testMessage)
	})

	testLogOutput(t, logrus.TraceLevel, testField, func() {
		getStandardLoggerWithFields().Tracef(simpleFormat, testMessage)
	})
}

func TestPanicf(t *testing.T) {
	testLogOutput(t, logrus.PanicLevel, frozen.NewMap(), func() {
		require.Panics(t, func() {
			NewStandardLogger().Panicf(simpleFormat, testMessage)
		})
	})

	testLogOutput(t, logrus.PanicLevel, testField, func() {
		require.Panics(t, func() {
			getStandardLoggerWithFields().Panicf(simpleFormat, testMessage)
		})
	})
}

func testLogOutput(t *testing.T, level logrus.Level, fields frozen.Map, logFunc func()) {
	outputtedFields := ""
	if fields.Count() != 0 {
		outputtedFields = testutil.OutputFormattedFields(fields)
	}

	expectedOutput := strings.Join([]string{outputtedFields, strings.ToUpper(level.String()), testMessage}, " ")
	actualOutput := testutil.RedirectOutput(t, logFunc)

	// uses Contains to avoid checking timestamps
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

	require.Equal(t, getNewStandardLogger().getFormattedField(), "")
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

	expected := "byteVal=107 numberVal=1 stringVal=this is a sentence"
	assert.Equal(t, expected, logger.getFormattedField())
}

func TestInsertFieldsKeyEmpty(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	logger.insertFieldsKey()
	assert.Equal(t, 0, len(logger.sortedFields))
}

func TestInsertFieldsKey(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	fields := []string{"some", "random", "fields"}
	logger.insertFieldsKey(fields...)

	sort.Strings(fields)
	assert.Equal(t, fields, logger.sortedFields)
}

func TestInsertFieldsKeyAddMoreFields(t *testing.T) {
	t.Parallel()

	logger := getNewStandardLogger()
	fields1 := []string{"some", "random", "fields"}
	fields2 := []string{"even", "more", "stuff"}

	logger.insertFieldsKey(fields1...)
	logger.insertFieldsKey(fields2...)

	combined := append(fields1, fields2...)
	sort.Strings(combined)
	assert.Equal(t, combined, logger.sortedFields)
}

func TestSetInfo(t *testing.T) {
	cases := testutil.GenerateMultipleFieldsCases()
	for _, c := range cases {
		t.Run("TestSetInfo"+c.Name, func(mc testutil.MultipleFields) func(*testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				logger := getNewStandardLogger().PutFields(mc.Fields).(*standardLogger)
				entry := logger.setInfo()
				expected := testutil.OutputFormattedFields(mc.Fields)

				assert.Equal(tt, expected, entry.Data[keyFields])
			}
		}(c))
	}
}

func TestWithFields(t *testing.T) {
	cases := testutil.GenerateMultipleFieldsCases()
	for _, c := range cases {
		t.Run("TestWithFields"+c.Name,
			func(mc testutil.MultipleFields) func(*testing.T) {
				return func(tt *testing.T) {
					tt.Parallel()

					logger := getNewStandardLogger()

					if mc.Fields.Count() == 0 {
						require.Panics(tt, func() {
							logger.PutFields(mc.Fields)
						})
						return
					}

					logger.PutFields(mc.Fields)
					expectedKeys := testutil.GetSortedKeys(mc.Fields)
					assert.Equal(tt, expectedKeys, logger.sortedFields)
					assert.True(tt, mc.Fields.Equal(logger.fields))
				}
			}(c))
	}
}

func TestWithField(t *testing.T) {
	cases := testutil.GenerateSingleFieldCases()
	for _, c := range cases {
		t.Run("TestWithField"+c.Name,
			func(sc testutil.SingleField) func(*testing.T) {
				return func(tt *testing.T) {
					tt.Parallel()

					logger := getNewStandardLogger()
					logger.PutField(sc.Key, sc.Val)
					value, exists := logger.fields.Get(sc.Key)

					require.True(tt, exists)
					assert.Equal(tt, sc.Val, value)
				}
			}(c))
	}
}

func TestWithFieldWithAddingMoreValues(t *testing.T) {
	cases := testutil.GenerateMultipleFieldsCases()
	for _, c := range cases {
		t.Run("TestWithFieldWithAddingMoreValues"+c.Name,
			func(mc testutil.MultipleFields) func(*testing.T) {
				return func(tt *testing.T) {
					tt.Parallel()

					logger := getNewStandardLogger()

					for i := mc.Fields.Range(); i.Next(); {
						logger.PutField(i.Key().(string), i.Value())
					}

					expectedKeys := testutil.GetSortedKeys(mc.Fields)
					assert.Equal(tt, expectedKeys, logger.sortedFields)
					assert.True(tt, mc.Fields.Equal(logger.fields))
				}
			}(c))
	}
}

func TestWithFieldReplaceValues(t *testing.T) {
	t.Parallel()

	key := "random"
	oldVal := 1
	newVal := 2

	logger := getNewStandardLogger()

	logger.PutField(key, oldVal)
	assertFieldExists(t, logger, frozen.NewMap(frozen.KV(key, oldVal)))

	logger.PutField(key, newVal)
	assertFieldExists(t, logger, frozen.NewMap(frozen.KV(key, newVal)))
	assert.Equal(t, []string{key}, logger.sortedFields)
}

func TestWithFieldsReplaceValues(t *testing.T) {
	t.Parallel()
	field := frozen.NewMap(
		frozen.KV("1", 1),
		frozen.KV("2", 2),
		frozen.KV("3", 3),
	)

	logger := getNewStandardLogger().PutFields(field).(*standardLogger)

	assertFieldExists(t, logger, field)

	for i := field.Range(); i.Next(); {
		field = field.With(i.Key(), "replaced")
	}
	logger.PutFields(field)

	assertFieldExists(t, logger, field)
	assert.Equal(t, testutil.GetSortedKeys(field), logger.sortedFields)
}

func assertFieldExists(t *testing.T, logger *standardLogger, expectedFields frozen.Map) {
	for i := expectedFields.Range(); i.Next(); {
		curVal, exists := logger.fields.Get(i.Key())
		require.True(t, exists)
		assert.Equal(t, i.Value(), curVal)
	}
}

func getNewStandardLogger() *standardLogger {
	return NewStandardLogger().(*standardLogger)
}

func getStandardLoggerWithFields() *standardLogger {
	logger := NewStandardLogger().PutFields(testField)
	return logger.(*standardLogger)
}

package log

import (
	"context"
	"strconv"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arr-ai/frozen"
	"github.com/stretchr/testify/mock"
)

type fieldsTest struct {
	name          string
	unresolveds   frozen.Map
	contextFields frozen.Map
	expected      frozen.Map
}

func TestMainDebug(t *testing.T) {
	t.Parallel()

	testLog(t, Debug, "Debug")
}

func TestMainDebugf(t *testing.T) {
	t.Parallel()

	testLogWithFormat(t, Debugf, "Debugf")
}

func TestMainInfo(t *testing.T) {
	t.Parallel()

	testLog(t, Info, "Info")
}

func TestMainInfof(t *testing.T) {
	t.Parallel()

	testLogWithFormat(t, Infof, "Infof")
}

func testLog(t *testing.T, logFunc func(ctx context.Context, args ...interface{}), funcName string) {
	args := []interface{}{"this is a message", 12345, 2.3141, 'k'}
	callLog(t, funcName, args,
		func(m *mockLogger) {
			logFunc(WithLogger(m).Onto(context.Background()), args...)
		},
	)
}

func testLogWithFormat(t *testing.T, logFunc func(ctx context.Context, format string, args ...interface{}), funcName string) {
	args := []interface{}{"this is a format %v %v %v %v", "this is a message", 12345, 2.3141, 'k'}
	callLog(t, funcName, args,
		func(m *mockLogger) {
			logFunc(WithLogger(m).Onto(context.Background()), args[0].(string), args[1:]...)
		},
	)
}

func callLog(t *testing.T, funcName string, args []interface{}, logFunc func(*mockLogger)) {
	logger := newMockLogger()
	setLogMockAssertion(logger, frozen.NewMap())
	logger.On(funcName, args...)
	logFunc(logger)
	logger.AssertExpectations(t)
}

func TestChain(t *testing.T) {
	t.Parallel()

	init := Fields{generateSimpleField(5)}
	fields1 := Fields{frozen.Map{}.With("6", 6).With("7", suppress{})}
	fields2 := Fields{
		frozen.Map{}.
			With("2", 10).
			With("8", func(context.Context) interface{} {
				return 8
			}),
	}
	fields3 := Fields{frozen.Map{}.With("7", 7).With("8", suppress{})}
	expected := generateSimpleField(5).
		With("2", 10).
		With("6", 6).
		With("7", 7).
		With("8", suppress{})

	assert.True(t, expected.Equal(init.Chain(fields1, fields2, fields3).m))
}

func TestWithConfigsSameConfigType(t *testing.T) {
	t.Parallel()

	expectedConfig := frozen.Map{}.
		With(standardFormat{}.TypeKey(), standardFormat{})
	f := WithConfigs(NewJSONFormat(), NewStandardFormat())
	assert.True(t, expectedConfig.Equal(f.m))
}

func TestFrom(t *testing.T) {
	for _, c := range getUnresolvedFieldsCases() {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			logger := newMockLogger()
			setLogMockAssertion(logger, c.expected)
			ctx := context.WithValue(context.Background(), fieldsContextKey{}, c.unresolveds.With(loggerKey{}, logger))
			for i := c.contextFields.Range(); i.Next(); {
				ctx = context.WithValue(ctx, i.Key(), i.Value())
			}
			From(ctx)
			logger.AssertExpectations(t)
		})
	}
}

func TestOnto(t *testing.T) {
	cases := generateMultipleFieldsCases()
	for _, c := range cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			logger := newMockLogger()
			setMockCopyAssertion(logger)

			fields := WithLogger(logger)
			for i := c.Fields.Range(); i.Next(); {
				fields = fields.With(i.Key().(string), i.Value())
			}
			ctx := fields.Onto(context.Background())

			logger = getLoggerFromContext(t, ctx)

			logger.AssertExpectations(t)
		})
	}
}

func TestSuppress(t *testing.T) {
	runFieldsMethod(
		t,
		func(t *testing.T) {
			t.Parallel()

			assert.True(t,
				frozen.NewMapFromKeys(
					generateSimpleField(3).Keys(),
					func(_ interface{}) interface{} {
						return suppress{}
					},
				).
					Equal(Suppress("0", "1", "2").m),
			)
		},
		func(t *testing.T) {
			t.Parallel()

			initial := generateSimpleField(5)

			expected := generateSimpleField(5).
				With("2", suppress{}).
				With("4", suppress{}).
				With("5", suppress{})

			assert.True(t, expected.Equal(Fields{initial}.Suppress("2", "4", "5").m))
		},
	)
}

func TestWith(t *testing.T) {
	runFieldsMethod(t,
		func(t *testing.T) {
			t.Parallel()

			assert.True(t,
				generateSimpleField(5).
					Equal(
						With("0", 0).
							With("1", 1).
							With("2", 2).
							With("3", 3).
							With("4", 4).m),
			)
		},
		func(t *testing.T) {
			t.Parallel()

			simpleField := generateSimpleField(5)
			expected := simpleField.With("1", 4).With("2", 5)
			assert.True(t, expected.Equal(Fields{simpleField}.With("1", 4).With("2", 5).m))
		},
	)
}

func TestWithCtxRef(t *testing.T) {
	t.Parallel()

	f := WithCtxRef("key1", key1{}).WithCtxRef("key2", key2{}).WithCtxRef("key3", key3{})

	for i := f.m.Range(); i.Next(); {
		assert.IsType(t, ctxRef{}, i.Value())
	}
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	logger := newMockLogger()
	setMockCopyAssertion(logger)
	WithLogger(logger).Onto(context.Background())

	logger.AssertExpectations(t)
}

func setLogMockAssertion(logger *mockLogger, fields frozen.Map) {
	setMockCopyAssertion(logger)
	setPutFieldsAssertion(logger, fields)
}

func setPutFieldsAssertion(logger *mockLogger, fields frozen.Map) {
	logger.On(
		"PutFields",
		mock.MatchedBy(
			func(arg frozen.Map) bool {
				return fields.Equal(arg)
			},
		),
	).Return(logger)
}

func getLoggerFromContext(t *testing.T, ctx context.Context) *mockLogger {
	m, exists := ctx.Value(fieldsContextKey{}).(frozen.Map)
	if !exists {
		t.Fatal("Fields not set yet")
	}
	return m.MustGet(loggerKey{}).(*mockLogger)
}

func setMockCopyAssertion(logger *mockLogger) *mock.Call {
	// set to return the same logger for testing purposes, in real case it will return
	// a copied logger. Tests that use these usually are not checked for their return value
	// as the return value is mocked
	return logger.On("Copy").Return(logger)
}

func runFieldsMethod(t *testing.T, empty, nonEmpty func(*testing.T)) {
	t.Run("empty fields", empty)
	t.Run("non empty fields", nonEmpty)
}

func generateSimpleField(limit int) frozen.Map {
	keys := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	return frozen.NewMapFromKeys(
		frozen.NewSetFromStrings(keys...),
		func(a interface{}) interface{} {
			num, err := strconv.Atoi(a.(string))
			if err != nil {
				panic("not a number")
			}
			return num
		},
	)
}

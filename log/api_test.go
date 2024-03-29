package log

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arr-ai/frozen"
	"github.com/stretchr/testify/mock"
)

var (
	formattedArgs = []interface{}{"this is a format %v %v %v %v", "this is a message", 12345, 2.3141, 'k'}
	regularArgs   = []interface{}{"this is a message", 12345, 2.3141, 'k'}
	errMsg        = errors.New("this is an error")
)

type fieldsTest struct {
	name          string
	unresolveds   frozen.Map[any, any]
	contextFields frozen.Map[any, any]
	expected      frozen.Map[any, any]
}

type testHook struct{}

func (h *testHook) OnLogged(entry *LogEntry) error {
	return nil
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

func TestMainError(t *testing.T) {
	t.Parallel()

	callLog(
		t,
		"Error",
		append([]interface{}{errMsg}, regularArgs...),
		frozen.Map[any, any]{}.With(errMsgKey, errMsg.Error()),
		func(m *mockLogger) {
			Error(WithLogger(m).Onto(context.Background()), errMsg, regularArgs...)
		},
	)
}

func TestMainErrorf(t *testing.T) {
	t.Parallel()

	callLog(
		t,
		"Errorf",
		append([]interface{}{errMsg}, formattedArgs...),
		frozen.Map[any, any]{}.With(errMsgKey, errMsg.Error()),
		func(m *mockLogger) {
			Errorf(WithLogger(m).Onto(context.Background()), errMsg, formattedArgs[0].(string), formattedArgs[1:]...)
		},
	)
}

func testLog(
	t *testing.T,
	logFunc func(ctx context.Context, args ...interface{}),
	funcName string,
) {
	callLog(t, funcName, regularArgs, frozen.NewMap[any, any](),
		func(m *mockLogger) {
			logFunc(WithLogger(m).Onto(context.Background()), regularArgs...)
		},
	)
}

func testLogWithFormat(
	t *testing.T,
	logFunc func(ctx context.Context, format string, args ...interface{}),
	funcName string,
) {
	callLog(t, funcName, formattedArgs, frozen.NewMap[any, any](),
		func(m *mockLogger) {
			logFunc(WithLogger(m).Onto(context.Background()), formattedArgs[0].(string), formattedArgs[1:]...)
		},
	)
}

func callLog(
	t *testing.T,
	funcName string,
	args []interface{},
	fields frozen.Map[any, any],
	logFunc func(*mockLogger),
) {
	logger := newMockLogger()
	setLogMockAssertion(logger, fields)
	logger.On(funcName, args...)
	logFunc(logger)
	logger.AssertExpectations(t)
}

func TestChain(t *testing.T) {
	t.Parallel()

	init := Fields{generateSimpleField(5)}
	fields1 := Fields{frozen.Map[any, any]{}.With("6", 6).With("7", suppress{})}
	fields2 := Fields{
		frozen.Map[any, any]{}.
			With("2", 10).
			With("8", func(context.Context) interface{} {
				return 8
			}),
	}
	fields3 := Fields{frozen.Map[any, any]{}.With("7", 7).With("8", suppress{})}
	expected := generateSimpleField(5).
		With("2", 10).
		With("6", 6).
		With("7", 7).
		With("8", suppress{})

	assert.True(t, expected.Equal(init.Chain(fields1, fields2, fields3).m))
}

func TestWithConfigsSameConfigType(t *testing.T) {
	t.Parallel()

	expectedConfig := frozen.Map[any, any]{}.
		With(standardFormat{}.TypeKey(), standardFormat{}).
		With(verboseMode{}.TypeKey(), verboseMode{true})

	f := WithConfigs(
		NewJSONFormat(),
		NewStandardFormat(),
		SetVerboseMode(true),
	)
	assert.True(t, expectedConfig.Equal(f.m))
}

func TestWithConfigLevel(t *testing.T) {
	t.Parallel()

	logger := newMockLogger()
	setLogMockAssertion(logger, frozen.NewMap[any, any]())
	logger.On("SetVerbose", true).Return(nil)
	WithConfigs(SetVerboseMode(false), SetVerboseMode(true)).WithLogger(logger).From(context.Background())
	logger.AssertExpectations(t)
}

func TestWithConfigFormat(t *testing.T) {
	t.Parallel()

	logger := newMockLogger()
	setLogMockAssertion(logger, frozen.NewMap[any, any]())
	logger.On("SetFormatter", jsonFormat{}).Return(nil)
	WithConfigs(NewStandardFormat(), NewJSONFormat()).WithLogger(logger).From(context.Background())
	logger.AssertExpectations(t)
}

func TestWithConfigOutput(t *testing.T) {
	t.Parallel()

	logger := newMockLogger()
	setLogMockAssertion(logger, frozen.NewMap[any, any]())
	logger.On("SetOutput", &bytes.Buffer{}).Return(nil)
	WithConfigs(SetOutput(&bytes.Buffer{})).WithLogger(logger).From(context.Background())
	logger.AssertExpectations(t)
}

func TestWithHooks(t *testing.T) {
	t.Parallel()

	logger := newMockLogger()
	setLogMockAssertion(logger, frozen.NewMap[any, any]())
	logger.On("AddHooks", mock.Anything).Return(nil)
	WithConfigs(AddHooks(&testHook{})).WithLogger(logger).From(context.Background())
	logger.AssertExpectations(t)
}

func TestWithConfigLogCaller(t *testing.T) {
	t.Parallel()

	logger := newMockLogger()
	setLogMockAssertion(logger, frozen.NewMap[any, any]())
	logger.On("SetLogCaller", true).Return(nil)
	WithConfigs(SetLogCaller(false), SetLogCaller(true)).WithLogger(logger).From(context.Background())
	logger.AssertExpectations(t)
}

func TestNewForwardingHook(t *testing.T) {
	t.Parallel()

	targetBuffer := bytes.Buffer{}
	targetLogger := WithConfigs(SetOutput(&targetBuffer)).
		WithLogger(NewStandardLogger()).
		From(context.Background())
	sourceBuffer := bytes.Buffer{}
	sourceLogger := WithConfigs(SetOutput(&sourceBuffer), AddHooks(NewForwardingHook(targetLogger))).
		WithLogger(NewStandardLogger()).
		From(context.Background())

	sourceLogger.Info("message")
	assert.True(t, sourceBuffer.Len() > 0)
	assert.Equal(t, sourceBuffer, targetBuffer)
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

			logger = getLoggerFromContext(ctx, t)

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

func TestWithContextRef(t *testing.T) {
	t.Parallel()

	f := WithContextKey("key1", key1{}).WithContextKey("key2", key2{}).WithContextKey("key3", key3{})

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

func TestFieldsFrom(t *testing.T) {
	t.Parallel()

	fields := Fields{}.With("key", "value")
	ctx := fields.Onto(context.Background())
	retrieved := FieldsFrom(ctx)
	assert.Equal(t, fields, retrieved)
}

func setLogMockAssertion(logger *mockLogger, fields frozen.Map[any, any]) {
	setMockCopyAssertion(logger)
	setPutFieldsAssertion(logger, fields)
}

func setPutFieldsAssertion(logger *mockLogger, fields frozen.Map[any, any]) {
	logger.On(
		"PutFields",
		mock.MatchedBy(fields.Equal),
	).Return(logger)
}

func getLoggerFromContext(ctx context.Context, t *testing.T) *mockLogger {
	m, exists := ctx.Value(fieldsContextKey{}).(frozen.Map[any, any])
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

func generateSimpleField(limit int) frozen.Map[any, any] {
	keys := make([]any, 0, limit)
	for i := 0; i < limit; i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	return frozen.NewMapFromKeys(
		frozen.NewSet[any](keys...),
		func(a interface{}) interface{} {
			num, err := strconv.Atoi(a.(string))
			if err != nil {
				panic("not a number")
			}
			return num
		},
	)
}

package log

import (
	"context"
	"testing"

	"github.com/arr-ai/frozen"
	"github.com/stretchr/testify/assert"
)

type key1 struct{}
type key2 struct{}
type key3 struct{}

func getUnresolvedFieldsCases() []fieldsTest {
	return []fieldsTest{
		{
			name: "regular unresolved fields",
			unresolveds: frozen.NewMap(
				frozen.KV("a", 1),
				frozen.KV("b", 2),
				frozen.KV("c", ctxRef{key1{}}),
				frozen.KV("d", suppress{}),
				frozen.KV("e", func(context.Context) interface{} { return "f" }),
			),
			contextFields: frozen.NewMap(
				frozen.KV(key1{}, "g"),
			),
			expected: frozen.NewMap(
				frozen.KV("a", 1),
				frozen.KV("b", 2),
				frozen.KV("c", "g"),
				frozen.KV("e", "f"),
			),
		},
		{
			name: "key does not exist in context",
			unresolveds: frozen.NewMap(
				frozen.KV("a", 1),
				frozen.KV("b", 2),
				frozen.KV("c", ctxRef{key1{}}),
				frozen.KV("d", suppress{}),
				frozen.KV("e", func(context.Context) interface{} { return "f" }),
			),
			expected: frozen.NewMap(
				frozen.KV("a", 1),
				frozen.KV("b", 2),
				frozen.KV("e", "f"),
			),
		},
		{
			name: "nothing to resolve",
			unresolveds: frozen.NewMap(
				frozen.KV("a", 1),
				frozen.KV("b", 2),
				frozen.KV("c", 3),
			),
			expected: frozen.NewMap(
				frozen.KV("a", 1),
				frozen.KV("b", 2),
				frozen.KV("c", 3),
			),
		},
	}
}

func TestConfigureLogger(t *testing.T) {
	for _, c := range getUnresolvedFieldsCases() {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			logger := newMockLogger()
			setPutFieldsAssertion(logger, c.expected)
			ctx := context.Background()
			for i := c.contextFields.Range(); i.Next(); {
				ctx = context.WithValue(ctx, i.Key(), i.Value())
			}
			logger = configureLogger(Logger(logger).(fieldSetter), c.expected, []Config{}).(*mockLogger)
			logger.AssertExpectations(t)
		})
	}
}

func TestConfigureLoggerWithConfigs(t *testing.T) {
	t.Parallel()

	//TODO: add more configs
	testCase := getUnresolvedFieldsCases()[0]
	unresolved := Fields{testCase.unresolveds}.WithConfigs(NewJSONFormat())
	expected := testCase.expected

	logger := newMockLogger()
	setPutFieldsAssertion(logger, expected)
	logger.On("SetFormatter", NewJSONFormat()).Return(nil)

	ctx := context.Background()
	for i := testCase.contextFields.Range(); i.Next(); {
		ctx = context.WithValue(ctx, i.Key(), i.Value())
	}
	resolved, configs := unresolved.getResolvedFields(ctx)
	configureLogger(Logger(logger).(fieldSetter), resolved, configs)
	logger.AssertExpectations(t)
}

func TestRegisterListener(t *testing.T) {
	t.Parallel()

	mockListener1 := OnLog

	ctx := RegisterListener(context.Background(), mockListener1)
	cbs, isCallbackList := ctx.Value(listenerKey{}).([]func(context.Context, Fields))
	assert.True(t, isCallbackList)
	assert.Equal(t, 1, len(cbs))

	mockListener2 := func(context.Context, Fields){}
	ctx = RegisterListener(ctx, mockListener2)
	cbs, isCallbackList = ctx.Value(listenerKey{}).([]func(context.Context, Fields))
	assert.True(t, isCallbackList)
	assert.Equal(t, 2, len(cbs))
}

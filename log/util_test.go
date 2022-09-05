package log

import (
	"context"
	"testing"

	"github.com/arr-ai/frozen"
)

type key1 struct{}
type key2 struct{}
type key3 struct{}

func getUnresolvedFieldsCases() []fieldsTest {
	return []fieldsTest{
		{
			name: "regular unresolved fields",
			unresolveds: frozen.NewMap(
				frozen.KV[any, any]("a", 1),
				frozen.KV[any, any]("b", 2),
				frozen.KV[any, any]("c", ctxRef{key1{}}),
				frozen.KV[any, any]("d", suppress{}),
				frozen.KV[any, any]("e", func(context.Context) interface{} { return "f" }),
			),
			contextFields: frozen.NewMap(
				frozen.KV[any, any](key1{}, "g"),
			),
			expected: frozen.NewMap(
				frozen.KV[any, any]("a", 1),
				frozen.KV[any, any]("b", 2),
				frozen.KV[any, any]("c", "g"),
				frozen.KV[any, any]("e", "f"),
			),
		},
		{
			name: "key does not exist in context",
			unresolveds: frozen.NewMap(
				frozen.KV[any, any]("a", 1),
				frozen.KV[any, any]("b", 2),
				frozen.KV[any, any]("c", ctxRef{key1{}}),
				frozen.KV[any, any]("d", suppress{}),
				frozen.KV[any, any]("e", func(context.Context) interface{} { return "f" }),
			),
			expected: frozen.NewMap(
				frozen.KV[any, any]("a", 1),
				frozen.KV[any, any]("b", 2),
				frozen.KV[any, any]("e", "f"),
			),
		},
		{
			name: "nothing to resolve",
			unresolveds: frozen.NewMap(
				frozen.KV[any, any]("a", 1),
				frozen.KV[any, any]("b", 2),
				frozen.KV[any, any]("c", 3),
			),
			expected: frozen.NewMap(
				frozen.KV[any, any]("a", 1),
				frozen.KV[any, any]("b", 2),
				frozen.KV[any, any]("c", 3),
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
			Fields{c.unresolveds}.configureLogger(ctx, Logger(logger).(fieldSetter))
			logger.AssertExpectations(t)
		})
	}
}

func TestConfigureLoggerWithConfigs(t *testing.T) {
	t.Parallel()

	//TODO: add more configs
	testCase := getUnresolvedFieldsCases()[0]
	unresolveds := Fields{testCase.unresolveds}.WithConfigs(NewJSONFormat())
	expected := testCase.expected

	logger := newMockLogger()
	setPutFieldsAssertion(logger, expected)
	logger.On("SetFormatter", NewJSONFormat()).Return(nil)

	ctx := context.Background()
	for i := testCase.contextFields.Range(); i.Next(); {
		ctx = context.WithValue(ctx, i.Key(), i.Value())
	}

	unresolveds.configureLogger(ctx, Logger(logger).(fieldSetter))
	logger.AssertExpectations(t)
}

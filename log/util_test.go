package log

import (
	"context"
	"testing"

	"github.com/alecthomas/assert"
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

func TestResolveFields(t *testing.T) {
	for _, c := range getUnresolvedFieldsCases() {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			for i := c.contextFields.Range(); i.Next(); {
				ctx = context.WithValue(ctx, i.Key(), i.Value())
			}

			assert.True(t, c.expected.Equal(Fields{c.unresolveds}.resolveFields(ctx)))
		})
	}
}

func TestSetConfigToLogger(t *testing.T) {

}

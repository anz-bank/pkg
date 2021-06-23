package log

import (
	"context"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arr-ai/frozen"
)

func TestCanonicalLog(t *testing.T) {
	t.Parallel()

	ctx := WithCanonicalLog(context.Background())
	logger := newMockLogger()

	mb := getCanonicalFields(ctx)

	addValues(mb, logger)
	setLogMockAssertion(logger, mb.Finish().Without(frozen.NewSet(loggerKey{})))
	logger.On("Info")

	// add the same values again because Finish would remove the values from the map builder.
	addValues(mb, logger)
	CanonicalLog(ctx)

	logger.AssertExpectations(t)
}

func TestOnLogWithoutCanonicalFieldsSetup(t *testing.T) {
	t.Parallel()

	logger := newMockLogger()
	logger.On("Info", "Canonical fields have not been set up properly, fields will not be collected")
	setLogMockAssertion(logger, frozen.NewMap())
	OnLog(WithLogger(logger).Onto(context.Background()), Fields{frozen.NewMap()})
	logger.AssertExpectations(t)
}

func TestOnLog(t *testing.T) {
	t.Parallel()

	ctx := WithCanonicalLog(context.Background())
	mb := getCanonicalFields(ctx)
	addValues(mb, newMockLogger())

	testField := frozen.NewMap().With("additional", "value").With("1", 1).With("byte", 'l')
	OnLog(ctx, Fields{testField})

	m := mb.Finish()
	for i := testField.Range(); i.Next(); {
		val, exists := m.Get(i.Key())
		assert.True(t, exists && val == i.Value())
	}
}

func addValues(mb *frozen.MapBuilder, logger Logger) {
	mb.Put("test", 1)
	mb.Put("test2", 'k')
	mb.Put("test3", "string value")
	mb.Put(loggerKey{}, logger)
}
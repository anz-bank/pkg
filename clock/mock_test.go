package clock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := Now(ctx)
	m := NewMock()
	m.AddTicks(
		now.Add(0*time.Hour),
		now.Add(1*time.Hour),
		now.Add(2*time.Hour),
	)
	ctx = Onto(ctx, m)
	assert.Equal(t, now.Add(0*time.Hour), Now(ctx))
	assert.Equal(t, now.Add(1*time.Hour), Now(ctx))
	assert.Equal(t, now.Add(2*time.Hour), Now(ctx))
	assert.Panics(t, func() { Now(ctx) })
}

func TestMockNotImplemented(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = Onto(ctx, NewMock())

	assert.Panics(t, func() { After(ctx, time.Millisecond) })
	assert.Panics(t, func() { AfterFunc(ctx, time.Millisecond, func() {}) })
	assert.Panics(t, func() { NewTicker(ctx, time.Millisecond) })
	assert.Panics(t, func() { NewTimer(ctx, time.Millisecond) })
}

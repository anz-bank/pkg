package clock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeTravel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	t0 := Now(ctx)
	tt := NewTimeTravel(time.Hour)
	ctx = Onto(ctx, tt)
	assert.True(t, Now(ctx).Sub(t0) >= time.Hour)
	tt.Jump(time.Hour)
	assert.True(t, Now(ctx).Sub(t0) >= 2*time.Hour)
	tt.JumpTo(time.Now().Add(3 * time.Hour))
	assert.True(t, Now(ctx).Sub(t0) >= 3*time.Hour)
}

func TestTimeTravelWaiting(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = Onto(ctx, NewTimeTravel(time.Hour))

	realT0 := time.Now()
	t0 := Now(ctx)

	// Travelling forward 1h should take no time at all in the real world.
	After(ctx, time.Hour)

	realT1 := time.Now()
	t1 := Now(ctx)

	assert.True(t, realT1.Sub(realT0) < time.Second)
	assert.True(t, t1.Sub(t0) >= time.Hour)

	// Trying to travelling backwards 1h shouldn't do anything.
	After(ctx, -time.Hour)

	realT2 := time.Now()
	t2 := Now(ctx)

	assert.True(t, realT2.Sub(realT0) < time.Second)
	assert.True(t, t2.Sub(t0) >= time.Hour)
}

func TestTimeTravelNotImplemented(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = Onto(ctx, NewTimeTravel(0))

	assert.Panics(t, func() { AfterFunc(ctx, time.Millisecond, func() {}) })
	assert.Panics(t, func() { NewTicker(ctx, time.Millisecond) })
	assert.Panics(t, func() { NewTimer(ctx, time.Millisecond) })
}

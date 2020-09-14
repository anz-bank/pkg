package clock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClock(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	assert.NotPanics(t, func() { Sleep(ctx, 0) })
	assert.NotPanics(t, func() { Now(ctx) })
}

func TestClockWaiting(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	assert.NotPanics(t, func() { <-After(ctx, time.Millisecond) })

	ch := make(chan int)
	assert.NotPanics(t, func() {
		AfterFunc(ctx, 0, func() { ch <- 42 })
	})
	assert.Equal(t, 42, <-ch)
}

func TestClockRelative(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	loc, err := time.LoadLocation("")
	require.NoError(t, err)

	epoch := time.Date(1970, time.January, 1, 0, 0, 0, 0, loc)
	assert.True(t, 0 < Since(ctx, epoch))

	nextMillenium := time.Date(3000, time.January, 1, 0, 0, 0, 0, loc)
	assert.True(t, 0 < Until(ctx, nextMillenium))
}

func TestClockTicker(t *testing.T) {
	// Time-sensitive. Do not parallelise.

	ctx := context.Background()

	t0 := Now(ctx)
	ticker := NewTicker(ctx, 100*time.Millisecond)
	t1 := <-ticker.C()
	t2 := <-ticker.C()
	ticker.Stop()
	assert.InEpsilon(t, float64(t1.Sub(t0)), float64(100*time.Millisecond), float64(10*time.Millisecond))
	assert.InEpsilon(t, float64(t2.Sub(t1)), float64(100*time.Millisecond), float64(10*time.Millisecond))
}

func TestClockTimer(t *testing.T) {
	// Time-sensitive. Do not parallelise.

	ctx := context.Background()

	t0 := Now(ctx)
	timer := NewTimer(ctx, 100*time.Millisecond)
	t1 := <-timer.C()
	assert.False(t, timer.Stop())
	assert.InEpsilon(t, float64(t1.Sub(t0)), float64(100*time.Millisecond), float64(10*time.Millisecond))
}

func TestClockTimerReadAfterFiring(t *testing.T) {
	// Time-sensitive. Do not parallelise.

	ctx := context.Background()

	t0 := Now(ctx)
	timer := NewTimer(ctx, 100*time.Millisecond)
	Sleep(ctx, 200*time.Millisecond)
	if assert.False(t, timer.Stop()) {
		t1 := <-timer.C()
		assert.InEpsilon(t, float64(t1.Sub(t0)), float64(100*time.Millisecond), float64(10*time.Millisecond))
	}
}

func TestClockTimerStop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	timer := NewTimer(ctx, 100*time.Millisecond)
	assert.True(t, timer.Stop())
}

func TestClockTimerReset(t *testing.T) {
	// Time-sensitive. Do not parallelise.

	ctx := context.Background()

	t0 := Now(ctx)
	timer := NewTimer(ctx, 200*time.Millisecond)
	Sleep(ctx, 50*time.Millisecond)
	assert.True(t, timer.Reset(400*time.Millisecond))
	t1 := <-timer.C()
	assert.False(t, timer.Stop())
	assert.InEpsilon(t, float64(t1.Sub(t0)), float64(400*time.Millisecond), float64(20*time.Millisecond))
}

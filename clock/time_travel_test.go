package clock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertWithinOneSecond(t *testing.T, expected, actual time.Time, msgAndArgs ...interface{}) bool {
	d := actual.Sub(expected)
	if assert.True(t, -time.Second <= d && d <= time.Second, "expected = %v\nactual   = %v", expected, actual) {
		return true
	}
	if len(msgAndArgs) > 0 {
		t.Logf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return false
}

func TestTimeTravel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	t0 := Now(ctx)

	tt := NewTimeTravel(time.Hour)
	defer tt.Close()
	ctx = Onto(ctx, tt)
	assert.True(t, assertWithinOneSecond(t, t0.Add(time.Hour), Now(ctx)))
	<-tt.After(time.Hour)
	assert.True(t, assertWithinOneSecond(t, t0.Add(2*time.Hour), Now(ctx)))
	<-tt.After(time.Hour)
	assert.True(t, assertWithinOneSecond(t, t0.Add(3*time.Hour), Now(ctx)))
}

func TestTimeTravelStartingAt(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	t0 := time.Unix(0, 0)
	tt := NewTimeTravelStartingAt(t0.Add(time.Hour))
	defer tt.Close()
	ctx = Onto(ctx, tt)
	assert.True(t, assertWithinOneSecond(t, t0.Add(time.Hour), Now(ctx)))
	<-tt.After(time.Hour)
	assert.True(t, assertWithinOneSecond(t, t0.Add(2*time.Hour), Now(ctx)))
	<-tt.After(time.Hour)
	assert.True(t, assertWithinOneSecond(t, t0.Add(3*time.Hour), Now(ctx)))
}

func TestTimeTravelWaiting(t *testing.T) {
	t.Parallel()

	tt := NewTimeTravel(time.Hour)
	defer tt.Close()
	ctx := Onto(context.Background(), tt)

	realT0 := time.Now()
	t0 := Now(ctx)

	// Travelling forward 1h should take no time at all in the real world.
	<-After(ctx, time.Hour)

	realT1 := time.Now()
	t1 := Now(ctx)

	assertWithinOneSecond(t, realT0, realT1)
	assertWithinOneSecond(t, t0.Add(time.Hour), t1)

	// Trying to travel backwards 1h should have no effect.
	<-After(ctx, -time.Hour)

	realT2 := time.Now()
	t2 := Now(ctx)

	assertWithinOneSecond(t, realT0, realT2)
	assertWithinOneSecond(t, t0.Add(time.Hour), t2)
}

func TestTimeTravelMultipleAfters(t *testing.T) {
	t.Parallel()

	type wait struct {
		wait     time.Duration
		expected time.Duration
	}
	cases := [][]wait{
		{{time.Hour, time.Hour}, {2 * time.Hour, 2 * time.Hour}},
		{{2 * time.Hour, 2 * time.Hour}, {time.Hour, 2 * time.Hour}},
	}
	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()
			ctx := Onto(context.Background(), NewTimeTravel(0))
			t0 := Now(ctx)
			chs := make([]<-chan time.Time, 0, len(c))
			for _, w := range c {
				chs = append(chs, After(ctx, w.wait))
			}
			for j, w := range c {
				now := <-chs[j]
				assertWithinOneSecond(t, t0.Add(w.expected), now, "wait %d", j)
			}
		})
	}
}

func TestTimeTravelNotImplemented(t *testing.T) {
	t.Parallel()

	ctx := Onto(context.Background(), NewTimeTravel(0))

	assert.Panics(t, func() { AfterFunc(ctx, time.Millisecond, func() {}) })
	assert.Panics(t, func() { NewTicker(ctx, time.Millisecond) })
	assert.Panics(t, func() { NewTimer(ctx, time.Millisecond) })
}

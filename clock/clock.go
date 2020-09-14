// Package clock is a context-driven wrapper for the time library. It allows
// substitution of mock clocks via context.Context for testing and other
// purposes.
package clock

import (
	"context"
	"time"
)

// Clock models the system clock. The functions mirror their counterparts in the
// time package.
type Clock interface {
	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) Timer
	NewTicker(d time.Duration) Ticker
	NewTimer(d time.Duration) Timer
	Now() time.Time
}

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type Timer interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

type clockKey struct{}

// NewClock sets the context clock. Pass nil to revert to the default system clock.
func Onto(ctx context.Context, clock Clock) context.Context {
	return context.WithValue(ctx, clockKey{}, clock)
}

func From(ctx context.Context) Clock {
	if clock := ctx.Value(clockKey{}); clock != nil {
		return clock.(Clock)
	}
	return defaultClock{}
}

func After(ctx context.Context, d time.Duration) <-chan time.Time {
	return From(ctx).After(d)
}

func AfterFunc(ctx context.Context, d time.Duration, f func()) Timer {
	return From(ctx).AfterFunc(d, f)
}

func NewTicker(ctx context.Context, d time.Duration) Ticker {
	return From(ctx).NewTicker(d)
}

func NewTimer(ctx context.Context, d time.Duration) Timer {
	return From(ctx).NewTimer(d)
}

func Now(ctx context.Context) time.Time {
	return From(ctx).Now()
}

func Since(ctx context.Context, t time.Time) time.Duration {
	return From(ctx).Now().Sub(t)
}

func Sleep(ctx context.Context, d time.Duration) {
	<-From(ctx).After(d)
}

func Until(ctx context.Context, t time.Time) time.Duration {
	return t.Sub(From(ctx).Now())
}

package clock

import "time"

type defaultClock struct{}

var _ Clock = defaultClock{}

func (defaultClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (defaultClock) AfterFunc(d time.Duration, f func()) Timer {
	return defaultTimer{timer: time.AfterFunc(d, f)}
}

func (defaultClock) NewTicker(d time.Duration) Ticker {
	return defaultTicker{ticker: time.NewTicker(d)}
}

func (defaultClock) NewTimer(d time.Duration) Timer {
	return defaultTimer{timer: time.NewTimer(d)}
}

func (defaultClock) Now() time.Time {
	return time.Now()
}

type defaultTicker struct {
	ticker *time.Ticker
}

func (t defaultTicker) C() <-chan time.Time {
	return t.ticker.C
}

func (t defaultTicker) Stop() {
	t.ticker.Stop()
}

type defaultTimer struct {
	timer *time.Timer
}

func (t defaultTimer) C() <-chan time.Time {
	return t.timer.C
}

func (t defaultTimer) Reset(d time.Duration) bool {
	return t.timer.Reset(d)
}

func (t defaultTimer) Stop() bool {
	return t.timer.Stop()
}

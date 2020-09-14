package clock

import "time"

// TimeTravel is a mock Clock that offsets time by a given offset. Also, when
// a call is made that involves a time delay, it travels to that time and
// returns instantly.
type TimeTravel struct {
	offset time.Duration
}

var _ Clock = &TimeTravel{}

func NewTimeTravel(initialOffset time.Duration) *TimeTravel {
	return &TimeTravel{offset: initialOffset}
}

func (t *TimeTravel) Jump(d time.Duration) {
	t.offset += d
}

func (t *TimeTravel) JumpTo(target time.Time) {
	t.offset = time.Until(target)
}

func (t *TimeTravel) After(d time.Duration) <-chan time.Time {
	// Waiting never travels backwards in time.
	if d < 0 {
		d = 0
	}
	t.offset += d
	ch := make(chan time.Time, 1)
	ch <- t.Now()
	return ch
}

func (*TimeTravel) AfterFunc(d time.Duration, f func()) Timer {
	panic("not implemented")
}

func (*TimeTravel) NewTicker(d time.Duration) Ticker {
	panic("not implemented")
}

func (*TimeTravel) NewTimer(d time.Duration) Timer {
	panic("not implemented")
}

func (t *TimeTravel) Now() time.Time {
	return time.Now().Add(t.offset)
}

package clock

import "time"

// TimeTravel is a mock Clock that offsets time by a given offset. Also, when
// a call is made that involves a time delay, it travels to that time and
// returns instantly.
type TimeTravel struct {
	offset time.Duration
}

var _ Clock = &TimeTravel{}

// NewTimeTravel creates a TimeTravel Clock.
func NewTimeTravel(initialOffset time.Duration) *TimeTravel {
	return &TimeTravel{offset: initialOffset}
}

// Travel through time by the duration. Backwards time travel is allowed.
func (t *TimeTravel) Jump(d time.Duration) {
	t.offset += d
}

// Travel to a point in time. Backwards time travel is allowed.
func (t *TimeTravel) JumpTo(target time.Time) {
	t.offset = time.Until(target)
}

// Wait returns a channel that time-travels by the duration and immediately
// sends the new time on the returned channel immediately. This function won't
// time-travel backwards. Negative durations will be treated as 0.
func (t *TimeTravel) After(d time.Duration) <-chan time.Time {
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

// Now returns the current time adjusted by the time-travel offset.
func (t *TimeTravel) Now() time.Time {
	return time.Now().Add(t.offset)
}

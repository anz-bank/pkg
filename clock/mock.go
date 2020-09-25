package clock

import "time"

// Mock implements a mock Clock that reports a fixed sequence of times each time
// it is queried.
type Mock struct {
	pending []time.Time
}

var _ Clock = &Mock{}

// NewMock creates a Mock Clock.
func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) AddTicks(t ...time.Time) {
	m.pending = append(m.pending, t...)
}

func (*Mock) After(d time.Duration) <-chan time.Time {
	panic("not implemented")
}

func (*Mock) AfterFunc(d time.Duration, f func()) Timer {
	panic("not implemented")
}

func (*Mock) NewTicker(d time.Duration) Ticker {
	panic("not implemented")
}

func (*Mock) NewTimer(d time.Duration) Timer {
	panic("not implemented")
}

func (m *Mock) Now() time.Time {
	if len(m.pending) == 0 {
		panic("clock.Mock.Now: no ticks available")
	}
	t := m.pending[0]
	m.pending = m.pending[1:]
	return t
}

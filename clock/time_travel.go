package clock

import (
	"reflect"
	"time"
)

// TimeTravel is a mock Clock that offsets time by a given offset. Also, when
// a call is made that involves a time delay, it travels to that time and
// returns instantly.
type TimeTravel struct {
	stop  chan struct{}
	after chan afterRequest
	now   chan (chan<- time.Time)
}

var _ Clock = &TimeTravel{}

// NewTimeTravel creates a TimeTravel Clock.
func NewTimeTravel(initialOffset time.Duration) *TimeTravel {
	t := &TimeTravel{
		stop:  make(chan struct{}),
		after: make(chan afterRequest),
		now:   make(chan (chan<- time.Time)),
	}
	go t.run(initialOffset)
	return t
}

func NewTimeTravelStartingAt(when time.Time) *TimeTravel {
	return NewTimeTravel(time.Until(when))
}

func (t *TimeTravel) run(offset time.Duration) {
	cases := []reflect.SelectCase{
		{Dir: reflect.SelectSend, Chan: reflect.ValueOf(t.stop), Send: reflect.ValueOf(struct{}{})},
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.after)},
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.now)},
	}

	nowPlus := func(d time.Duration) time.Time {
		return time.Now().Add(offset + d)
	}

	for {
		chosen, recv, _ := reflect.Select(cases)
		switch chosen {
		case 0:
			return
		case 1:
			req := recv.Interface().(afterRequest)
			if req.d < 0 {
				req.d = 0
			}
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(req.ch),
				Send: reflect.ValueOf(nowPlus(req.d)),
			})
		case 2:
			ch := recv.Interface().(chan<- time.Time)
			ch <- nowPlus(0)
		default:
			const tail = 3
			c := cases[chosen]
			cases = append(cases[:chosen], cases[chosen+1:]...)

			sent := c.Send.Interface().(time.Time)
			offset = time.Until(sent)

			waiters := cases[tail:]
			for i, w := range waiters {
				if w.Send.Interface().(time.Time).Before(sent) {
					w.Send = reflect.ValueOf(nowPlus(0))
					waiters[i] = w
				}
			}
		}
	}
}

func (t *TimeTravel) Close() {
	<-t.stop
}

// Wait returns a channel that time-travels by the duration and immediately
// sends the new time on the returned channel immediately. To conform to the
// real-world time API and avoid paradoxes, this function won't time-travel
// backwards. Negative durations will be treated as 0.
func (t *TimeTravel) After(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time)
	t.after <- afterRequest{d: d, ch: ch}
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
	ch := make(chan time.Time)
	t.now <- ch
	return <-ch
}

type afterRequest struct {
	d  time.Duration
	ch chan<- time.Time
}

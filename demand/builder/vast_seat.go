package builder

import (
	"time"
)

type vastSeat struct {
	seat
	linear    bool
	duration  int
	skipAfter int
}

// Linear is only usable in vast subsystem!
func (s *vastSeat) Linear() bool {
	return s.linear
}

func (s *vastSeat) Duration() time.Duration {
	return time.Second * time.Duration(s.duration)
}

func (s *vastSeat) SkipAfter() time.Duration {
	return time.Second * time.Duration(s.skipAfter)
}

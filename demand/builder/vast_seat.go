package builder

import (
	"fmt"
	"net/url"
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

func (s *vastSeat) ImpressionURL() *url.URL {
	if s.imp != nil {
		return s.imp
	}
	if s.winnerAd == nil {
		panic("no winner")
	}

	s.imp = s.makeURL(
		"pixel",
		map[string]string{"rh": s.ReservedHash(), "size": fmt.Sprint(s.size), "type": s.Type(), "subtype": s.SubType()},
		s.cpm,
		showExpire.Duration(),
	)
	return s.imp
}

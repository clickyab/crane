package builder

import (
	"fmt"
	"net/url"
	"time"

	"clickyab.com/crane/demand/entity"
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

// @override
func (s *vastSeat) ImpressionURL() *url.URL {
	if s.imp != nil {
		return s.imp
	}
	if s.winnerAd == nil {
		panic("no winner")
	}

	s.imp = s.makeURL(
		"pixel",
		map[string]string{"rh": s.ReservedHash(), "size": fmt.Sprint(s.size), "type": s.Type().String(), "subtype": s.RequestType().String(), "pt": s.context.publisher.Type().String()},
		s.cpm,
		showExpire.Duration(),
	)
	return s.imp
}

// @override
func (s *vastSeat) Acceptable(advertise entity.Creative) bool {
	if !s.genericTests(advertise) {
		return false
	}

	// TODO : the following line is correct. but, since we use an invalid form of ads in our system, we should comment it
	//return in.Campaign().Target() == entity.TargetVast
	if advertise.Target() != entity.TargetVast {
		// TODO : remove it when the new console is awaken!
		if advertise.Size() != 9 && advertise.Target() != entity.TargetWeb { // there is a fucking decision to show web size 9 in vast network.
			return false
		}
	}
	return true
}

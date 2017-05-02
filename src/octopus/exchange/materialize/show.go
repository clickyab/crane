package materialize

import (
	"encoding/json"
	"octopus/exchange"
	"services/broker"
	"time"
)

type show struct {
	data map[string]interface{}
	time time.Duration
	key  string

	src []byte
}

// Encode encode
func (s *show) Encode() ([]byte, error) {
	if s.src == nil {
		s.data["time"] = s.time
		s.src, _ = json.Marshal(s.data)
	}

	return s.src, nil
}

// Length return length
func (s *show) Length() int {
	x, _ := s.Encode()
	return len(x)
}

// Topic return topic
func (*show) Topic() string {
	return "show"
}

// Key return key
func (s *show) Key() string {
	return s.key
}

// Report report
func (*show) Report() func(error) {
	return func(error) {}
}

// ShowJob return a broker job
func ShowJob(imp exchange.Impression, dmn exchange.Demand, ad exchange.Advertise, t time.Duration, slotID string) broker.Job {
	return &show{
		data: winnerToMap(imp, dmn, ad, slotID),
		time: t,
		key:  imp.IP().String(),
	}
}

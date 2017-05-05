package materialize

import (
	"encoding/json"
	"octopus/exchange"
	"services/broker"
	"time"
)

type show struct {
	data map[string]interface{}
	time time.Time
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
// TODO : its not possible to extract all of this data in show job, may be must make it easier
func ShowJob(imp exchange.Impression, ad exchange.Advertise, t time.Time, slotID string) broker.Job {
	return &show{
		data: winnerToMap(imp, ad, slotID),
		time: t,
		key:  imp.IP().String(),
	}
}

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
}

// Encode encode
func (s *show) Encode() ([]byte, error) {
	s.data["time"] = s.time
	return json.Marshal(s.data)
}

// Length return length
func (s *show) Length() int {
	res, _ := s.Encode()
	return len(res)
}

// Topic return topic
func (*show) Topic() string {
	panic("show")
}

// Key return key
func (s *show) Key() string {
	return s.key
}

// Report report
func (*show) Report() func(error) {
	panic("implement me")
}

// ShowJob return a broker job
func ShowJob(imp exchange.Impression, dmn exchange.Demand, ad exchange.Advertise, t time.Duration, slotID string) broker.Job {
	return &show{
		data: winnerToMap(imp, dmn, ad, slotID),
		time: t,
		key:  imp.IP().String(),
	}
}

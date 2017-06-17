package jsonbackend

import (
	"encoding/json"

	"github.com/clickyab/services/broker"

	"github.com/Sirupsen/logrus"
)

type show struct {
	data map[string]interface{}
	time string
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
	return func(err error) {
		if err != nil {
			logrus.Warn(err)
		}
	}
}

// ShowJob return a broker job
func ShowJob(trackID, demand, slotID, adID string, IP string, winner int64, t string, supplier string, publisher string, profit int64) broker.Job {
	return &show{
		data: showToMap(trackID, demand, slotID, adID, winner, supplier, publisher, profit),
		time: t,
		key:  IP,
	}
}

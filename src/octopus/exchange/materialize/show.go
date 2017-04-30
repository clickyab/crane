package materialize

import (
	"octopus/exchange"
	"services/broker"
	"time"
)

type show struct {
}

func (a *show) Encode() ([]byte, error) {
	panic("implement me")
}

func (s *show) Length() int {
	panic("implement me")
}

func (s *show) Topic() string {
	panic("implement me")
}

func (s *show) Key() string {
	panic("implement me")
}

func (s *show) Report() func(error) {
	panic("implement me")
}

func ShowJob(imp exchange.Impression, dmn exchange.Demand, ad exchange.Advertise, t time.Duration, price int64, slotID string) broker.Job {
	return nil
}

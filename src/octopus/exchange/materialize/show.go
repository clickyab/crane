package materialize

import (
	"octopus/exchange"
	"services/broker"
	"time"
)

type show struct {
}

func (*show) Encode() ([]byte, error) {
	panic("implement me")
}

func (*show) Length() int {
	panic("implement me")
}

func (*show) Topic() string {
	panic("implement me")
}

func (*show) Key() string {
	panic("implement me")
}

func (*show) Report() func(error) {
	panic("implement me")
}

func ShowJob(imp exchange.Impression, dmn exchange.Demand, ad exchange.Advertise, t time.Duration, price int64, slotID string) broker.Job {
	return nil
}

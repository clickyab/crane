package materialize

import (
	"octopus/exchange"
	"services/broker"
)

type winner struct {
}

func (w *winner) Encode() ([]byte, error) {
	panic("implement me")
}

func (w *winner) Length() int {
	panic("implement me")
}

func (w *winner) Topic() string {
	panic("implement me")
}

func (w *winner) Key() string {
	panic("implement me")
}

func (w *winner) Report() func(error) {
	panic("implement me")
}

func WinnerJob(imp exchange.Impression, dmn exchange.Demand, cpm int64, ad exchange.Advertise, slotID string) broker.Job {
	return nil
}

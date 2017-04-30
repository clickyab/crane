package materialize

import (
	"octopus/exchange"
	"services/broker"
)

type demand struct {
}

func (d *demand) Encode() ([]byte, error) {
	panic("implement me")
}

func (d *demand) Length() int {
	panic("implement me")
}

func (d *demand) Topic() string {
	panic("implement me")
}

func (d *demand) Key() string {
	panic("implement me")
}

func (d *demand) Report() func(error) {
	panic("implement me")
}

func DemandJob(imp exchange.Impression, dmn exchange.Demand, ads []exchange.Advertise) broker.Job {
	return nil
}

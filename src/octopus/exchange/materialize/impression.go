package materialize

import (
	"octopus/exchange"
	"services/broker"
)

type impression struct {
}

func (i *impression) Encode() ([]byte, error) {
	panic("implement me")
}

func (i *impression) Length() int {
	panic("implement me")
}

func (i *impression) Topic() string {
	panic("implement me")
}

func (i *impression) Key() string {
	panic("implement me")
}

func (i *impression) Report() func(error) {
	panic("implement me")
}

func ImpressionJob(imp exchange.Impression) broker.Job {
	return nil
}

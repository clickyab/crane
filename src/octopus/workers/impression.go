package workers

import (
	"time"

	"services/assert"
	"services/broker"
	"services/safe"
)

type impressionModel struct {
	Time   time.Time `json:"time"`
	Source struct {
		Name     string `json:"name"`
		Supplier struct {
			Name string `json:"name"`
		} `json:"supplier"`
	} `json:"source"`

	Slots []struct{} `json:"slots"`
}

type impressionConsumer struct {
}

func (*impressionConsumer) Topic() string {
	return "show"
}

func (*impressionConsumer) Queue() string {
	panic("implement me")
}

func (s *impressionConsumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery, 0)
	safe.GoRoutine(func() {
		for {
			select {
			case del := <-chn:
				obj := impressionModel{}
				err := del.Decode(&obj)
				assert.Nil(err)
				dataChannel <- tableModel{
					Request:      1,
					Impression:   len(obj.Slots),
					Source:       obj.Source.Name,
					Supplier:     obj.Source.Supplier.Name,
					Time:         factTableID(obj.Time),
					Acknowledger: &del.(Acknowledger),
				}
			}
		}
	})

	return chn
}

func init() {
	broker.RegisterConsumer(impressionConsumer{})
}

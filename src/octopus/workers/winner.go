package workers

import (
	"services/assert"
	"services/broker"
	"services/safe"
	"time"
)

type winnerModel struct {
	Impression struct {
		Source struct {
			Name     string `json:"name"`
			Supplier struct {
				Name string `json:"name"`
			} `json:"supplier"`
		}
		Time time.Time `json:"time"`
	} `json:"impression"`
	Advertise struct {
		ImpressionBid int64 `json:"winner_cpm"`
		Demand        struct {
			Name string `json:"name"`
		} `json:"demand"`
	}
}

// winnerConsumer
type winnerConsumer struct {
}

func (*winnerConsumer) Topic() string {
	return "winner"
}

func (*winnerConsumer) Queue() string {
	return "winner_que"
}

func (w *winnerConsumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery)
	safe.GoRoutine(func() {
		for {
			select {
			case del := <-chn:
				obj := winnerModel{}
				err := del.Decode(&obj)
				assert.Nil(err)
				dataChannel <- tableModel{
					Supplier:      obj.Impression.Source.Supplier.Name,
					Source:        obj.Impression.Source.Name,
					Demand:        obj.Advertise.Demand.Name,
					ImpressionBid: obj.Advertise.ImpressionBid,
					Time:          factTableID(obj.Impression.Time),
					Acknowledger:  &del,
				}
			}
		}
	})

	return chn
}

func init() {
	broker.RegisterConsumer(&winnerConsumer{})
}

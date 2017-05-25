package winner

import (
	"context"
	"octopus/workers/internal/manager"
	"services/assert"
	"services/broker"
	"services/initializer"
	"services/safe"
	"time"
)

type model struct {
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

// consumer
type consumer struct {
	ctx context.Context
}

func (s *consumer) Initialize(ctx context.Context) {
	s.ctx = ctx
	broker.RegisterConsumer(s)
}

func (consumer) Topic() string {
	return "winner"
}

func (consumer) Queue() string {
	return "winner_aggregate"
}

func (s *consumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery)
	done := s.ctx.Done()
	safe.GoRoutine(func() {
		for {
			select {
			case del := <-chn:
				obj := model{}
				err := del.Decode(&obj)
				assert.Nil(err)
				manager.DataChannel <- manager.TableModel{
					Supplier:      obj.Impression.Source.Supplier.Name,
					Source:        obj.Impression.Source.Name,
					Demand:        obj.Advertise.Demand.Name,
					ImpressionBid: obj.Advertise.ImpressionBid,
					Time:          manager.FactTableID(obj.Impression.Time),
					Acknowledger:  &del,
				}
			case <-done:
				return
			}
		}
	})

	return chn
}

func init() {
	initializer.Register(&consumer{}, 10000)
}

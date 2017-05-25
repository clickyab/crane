package impression

import (
	"context"
	"time"

	"octopus/workers/internal/manager"
	"services/assert"
	"services/broker"
	"services/initializer"
	"services/safe"
)

type model struct {
	Time   time.Time `json:"time"`
	Source struct {
		Name     string `json:"name"`
		Supplier struct {
			Name string `json:"name"`
		} `json:"supplier"`
	} `json:"source"`

	Slots []struct{} `json:"slots"`
}

type consumer struct {
	ctx context.Context
}

func (s *consumer) Initialize(ctx context.Context) {
	s.ctx = ctx
	broker.RegisterConsumer(s)
}

func (consumer) Topic() string {
	return "impression"
}

func (consumer) Queue() string {
	return "impression_aggregate"
}

func (s *consumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery, 0)
	done := s.ctx.Done()
	safe.GoRoutine(func() {
		for {
			select {
			case del := <-chn:
				obj := model{}
				err := del.Decode(&obj)
				assert.Nil(err)
				manager.DataChannel <- manager.TableModel{
					Request:      1,
					Impression:   int64(len(obj.Slots)),
					Source:       obj.Source.Name,
					Supplier:     obj.Source.Supplier.Name,
					Time:         manager.FactTableID(obj.Time),
					Acknowledger: &del,
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

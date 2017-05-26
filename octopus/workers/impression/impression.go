package impression

import (
	"context"
	"time"

	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/initializer"
	"clickyab.com/exchange/services/safe"
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
				datamodels.ActiveAggregator().Channel() <- datamodels.TableModel{
					Request:      1,
					Impression:   int64(len(obj.Slots)),
					Source:       obj.Source.Name,
					Supplier:     obj.Source.Supplier.Name,
					Time:         datamodels.FactTableID(obj.Time),
					Acknowledger: del,
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

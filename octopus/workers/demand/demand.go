package demand

import (
	"context"
	"time"

	"clickyab.com/exchange/octopus/models"
	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/config"
	"clickyab.com/exchange/services/initializer"
	"clickyab.com/exchange/services/random"
	"clickyab.com/exchange/services/safe"
)

type model struct {
	Impression struct {
		Source struct {
			Name     string `json:"name"`
			Supplier struct {
				Name string `json:"name"`
			} `json:"supplier"`
		} `json:"source"`
		Time  time.Time `json:"time"`
		Slots []struct {
			Ad struct {
				MaxCPM int64 `json:"max_cpm,omitempty"`
			} `json:"ad,omitempty"`
		} `json:"slots"`
	} `json:"impression"`
	Demand struct {
		Name string `json:"name"`
	} `json:"demand"`
}

var extraCount = config.RegisterInt("octopus.workers.extra.count", 10, "the consumer count for a worker")

type consumer struct {
	ctx      context.Context
	workerID string
}

func (s *consumer) Initialize(ctx context.Context) {
	s.ctx = ctx
	broker.RegisterConsumer(s)

	for i := 1; i < *extraCount; i++ {
		broker.RegisterConsumer(
			&consumer{
				ctx:      ctx,
				workerID: <-random.ID,
			},
		)
	}
}

func (consumer) Topic() string {
	return "demand"
}

func (consumer) Queue() string {
	return "demand_aggregate"
}

func (s *consumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery, 0)
	done := s.ctx.Done()
	safe.ContinuesGoRoutine(func(cnl context.CancelFunc) {
		var del broker.Delivery
		defer func() {
			if del != nil {
				del.Reject(false)
			}
		}()
		for {
			select {
			case del = <-chn:
				obj := model{}
				err := del.Decode(&obj)
				assert.Nil(err)
				var win int64
				for i := range obj.Impression.Slots {
					if cpm := obj.Impression.Slots[i].Ad.MaxCPM; cpm > 0 {
						win++
					}
				}
				datamodels.ActiveAggregator().Channel() <- datamodels.TableModel{
					Supplier:           obj.Impression.Source.Supplier.Name,
					Source:             obj.Impression.Source.Name,
					Demand:             obj.Demand.Name,
					Time:               models.FactTableID(obj.Impression.Time),
					RequestOutCount:    1,
					ImpressionOutCount: int64(len(obj.Impression.Slots)),
					AdInCount:          win,
					Acknowledger:       del,
					WorkerID:           s.workerID,
				}
			case <-done:
				cnl()
				del = nil
				return
			}
		}
	}, time.Second)
	return chn
}

func init() {
	initializer.Register(&consumer{workerID: <-random.ID}, 1000)
}

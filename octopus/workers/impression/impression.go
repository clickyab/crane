package impression

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
	Time   time.Time `json:"time"`
	Source struct {
		Name     string `json:"name"`
		Supplier struct {
			Name string `json:"name"`
		} `json:"supplier"`
	} `json:"source"`

	Slots []struct{} `json:"slots"`
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
	return "impression"
}

func (consumer) Queue() string {
	return "impression_aggregate"
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
				datamodels.ActiveAggregator().Channel() <- datamodels.TableModel{
					Source:            obj.Source.Name,
					Supplier:          obj.Source.Supplier.Name,
					Time:              models.FactTableID(obj.Time),
					RequestInCount:    1,
					ImpressionInCount: int64(len(obj.Slots)),
					Acknowledger:      del,
					WorkerID:          s.workerID,
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
	initializer.Register(&consumer{workerID: <-random.ID}, 10000)
}

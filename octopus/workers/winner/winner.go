package winner

import (
	"context"
	"time"

	"clickyab.com/exchange/octopus/models"
	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/safe"
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
		WinnerCpm int64 `json:"winner_cpm"`
		Demand    struct {
			Name string `json:"name"`
		} `json:"demand"`
	}
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
	return "winner"
}

func (consumer) Queue() string {
	return "winner_aggregate"
}

func (s *consumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery)
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
					Supplier:     obj.Impression.Source.Supplier.Name,
					Source:       obj.Impression.Source.Name,
					Demand:       obj.Advertise.Demand.Name,
					Time:         models.FactTableID(obj.Impression.Time),
					AdOutCount:   1,
					AdOutBid:     obj.Advertise.WinnerCpm,
					Acknowledger: del,
					WorkerID:     s.workerID,
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

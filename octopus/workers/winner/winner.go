package winner

import (
	"context"
	"time"

	"clickyab.com/exchange/octopus/models"
	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/initializer"
	"clickyab.com/exchange/services/safe"
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
					WinCount:     1,
					WinBid:       obj.Advertise.WinnerCpm,
					Acknowledger: del,
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
	initializer.Register(&consumer{}, 10000)
}

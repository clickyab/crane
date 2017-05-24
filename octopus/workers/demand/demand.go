package demand

import (
	"context"
	"time"

	"clickyab.com/exchange/octopus/workers/internal/manager"
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
		} `json:"source"`
		Time time.Time `json:"time"`
	} `json:"impression"`
	Demand struct {
		Name string `json:"name"`
	} `json:"dem"`
	Ads map[string]struct {
		WinnerCPM int64 `json:"winner_cpm"`
	} `json:"ads"`
}

type consumer struct {
	ctx context.Context
}

func (s *consumer) Initialize(ctx context.Context) {
	s.ctx = ctx
	broker.RegisterConsumer(s)
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
	safe.GoRoutine(func() {
		for {
			select {
			case del := <-chn:
				obj := model{}
				err := del.Decode(&obj)
				assert.Nil(err)
				manager.DataChannel <- manager.TableModel{
					Supplier:     obj.Impression.Source.Supplier.Name,
					Source:       obj.Impression.Source.Name,
					Time:         manager.FactTableID(obj.Impression.Time),
					Demand:       obj.Demand.Name,
					Win:          int64(len(obj.Ads)),
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
	initializer.Register(&consumer{}, 1000)
}

package demand

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
				var winOut int64
				for i := range obj.Impression.Slots {
					winOut++
					if cpm := obj.Impression.Slots[i].Ad.MaxCPM; cpm > 0 {
						win++
					}
				}
				datamodels.ActiveAggregator().Channel() <- datamodels.TableModel{
					Supplier:           obj.Impression.Source.Supplier.Name,
					Source:             obj.Impression.Source.Name,
					Demand:             obj.Demand.Name,
					Time:               models.FactTableID(obj.Impression.Time),
					ImpressionInCount:  win,
					RequestOutCount:    1,
					ImpressionOutCount: winOut,
					Acknowledger:       del,
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
	initializer.Register(&consumer{}, 1000)
}

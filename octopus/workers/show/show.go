package show

import (
	"context"
	"strconv"
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

// TODO : is this model is correct? remove this tod if it is.
type model struct {
	TrackID    string `json:"track_id"`
	DemandName string `json:"demand_name"`
	Price      int64  `json:"price"`
	SlotID     string `json:"slot_id"`
	AdID       string `json:"ad_id"`
	Supplier   string `json:"supplier"`
	Publisher  string `json:"publisher"`
	Time       string `json:"time"`
	Profit     int    `json:"profit"`
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
	return "show"
}

func (consumer) Queue() string {
	return "show_aggregate"
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
					Supplier:     obj.Supplier,
					Source:       obj.Publisher,
					Demand:       obj.DemandName,
					DeliverBid:   obj.Price,
					DeliverCount: 1,
					Profit:       int64(obj.Profit),
					// TODO : why this is different with other?? make it same.
					Time:         models.FactTableID(timestampToTime(obj.Time)),
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

func timestampToTime(s string) time.Time {

	i, err := strconv.ParseInt(s, 10, 0)
	assert.Nil(err)
	return time.Unix(i, 0)

}

func init() {
	initializer.Register(&consumer{workerID: <-random.ID}, 10000)
}

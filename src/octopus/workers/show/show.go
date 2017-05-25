package show

import (
	"context"
	"octopus/workers/internal/manager"
	"services/assert"
	"services/broker"
	"services/initializer"
	"services/safe"
	"strconv"
	"time"
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
}

type consumer struct {
	ctx context.Context
}

func (s *consumer) Initialize(ctx context.Context) {
	s.ctx = ctx
	broker.RegisterConsumer(s)
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
	safe.GoRoutine(func() {
		for {
			select {
			case del := <-chn:
				obj := model{}
				err := del.Decode(&obj)
				assert.Nil(err)
				manager.DataChannel <- manager.TableModel{
					Supplier: obj.Supplier,
					Source:   obj.Publisher,
					Demand:   obj.DemandName,
					ShowBid:  obj.Price,
					Show:     1,
					// TODO : why this is different with other?? make it same.
					Time:         manager.FactTableID(timestampToTime(obj.Time)),
					Acknowledger: &del,
				}
			case <-done:
				return
			}
		}
	})
	return chn
}

func timestampToTime(s string) time.Time {

	i, err := strconv.ParseInt(s, 10, 0)
	assert.Nil(err)
	return time.Unix(i, 0)

}

func init() {
	initializer.Register(&consumer{}, 10000)
}

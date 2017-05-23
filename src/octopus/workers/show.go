package workers

import (
	"services/assert"
	"services/broker"
	"services/safe"
)

type showModel struct {
	TrackID    string `json:"track_id"`
	DemandName string `json:"demand_name"`
	Price      int64  `json:"price"`
	SlotID     string `json:"slot_id"`
	AdID       string `json:"ad_id"`
	Supplier   string `json:"supplier"`
	Publisher  string `json:"publisher"`
	Time       string `json:"time"`
}

type showConsumer struct {
}

func (*showConsumer) Topic() string {
	return "show"
}

func (*showConsumer) Queue() string {
	panic("implement me")
}

func (s *showConsumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery, 0)
	safe.GoRoutine(func() {
		for {
			select {
			case del := <-chn:
				obj := showModel{}
				err := del.Decode(&obj)
				assert.Nil(err)
				dataChannel <- tableModel{
					Supplier:     obj.Supplier,
					Source:       obj.Publisher,
					Demand:       obj.DemandName,
					ShowBid:      obj.Price,
					Show:         1,
					Time:         factTableID(timestampToTime(obj.Time)),
					Acknowledger: &del.(Acknowledger),
				}
			}
		}
	})
	return chn
}
func init() {
	broker.RegisterConsumer(showConsumer{})
}

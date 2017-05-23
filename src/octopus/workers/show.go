package workers

import (
	"services/assert"
	"services/broker"
	"services/safe"
)

type showConsumer struct {
	Data struct {
		TrackID    string `json:"track_id"`
		DemandName string `json:"demand_name"`
		Price      int64  `json:"price"`
		SlotID     string `json:"slot_id"`
		AdID       string `json:"ad_id"`
		Supplier   string `json:"supplier"`
		Publisher  string `json:"publisher"`
	} `json:"data"`
	Time string `json:"time"`
	IP   string `json:"ip"`
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
				obj := showConsumer{}
				err := del.Decode(&obj)
				assert.Nil(err)
				dataChannel <- tableModel{
					ShowBid: obj.Data.Price,
					Show:    1,
					Time:    factID(timestampToTime(obj.Time)),
				}
			}
		}
	})
	return chn

}

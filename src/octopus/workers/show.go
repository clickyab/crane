package workers

import (
	"time"

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
	safe.GoRoutine(func() { s.fillChan(chn) })
	return chn

}
func (s *showConsumer) fillChan(chn chan broker.Delivery) {
	obj := showConsumer{}
	h := time.After(2 * time.Minute)
	var ins = make(map[string]map[string]int64)
	var counter = 0
	var ack Acknowledger
	for {
		select {
		case del := <-chn:
			ack = del
			err := del.Decode(&obj)
			assert.Nil(err)
			//que = append(que, obj)

			key := genID(timestampToTime(obj.Time), obj.Data.Supplier, obj.Data.Publisher)
			val, ok := ins[key]
			if !ok {
				val = make(map[string]int64)
			}

			// increment
			val["request"] += 1
			val["impression"] += 1
			val["show"] += 1
			val["show_bid"] += obj.Data.Price

			ins[key] = val

			counter++
			if counter > limit {
				err := s.flush(ins)
				if err == nil {
					ack.Ack(true)
				} else {
					ack.Nack(true, true)

				}
				ins = make(map[string]map[string]int64)
				counter = 0
			}

		case <-h:
			s.flush(ins)
			ins = make(map[string]map[string]int64)
			counter = 0
			h = time.After(2 * time.Minute)
		}
	}
}

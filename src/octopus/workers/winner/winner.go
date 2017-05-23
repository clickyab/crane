package winner

import (
	"services/assert"
	"services/broker"
	"services/safe"
	"time"
)

const limit = 1000

// winnerConsumer asd
type winnerConsumer struct {
	Data struct {
		TrackID    string `json:"track_id"`
		DemandName string `json:"demand_name"`
		Price      int64  `json:"price"`
		SlotID     string `json:"slot_id"`
		AdID       string `json:"ad_id"`
	} `json:"data"`

	key string `json:"key"`
}

func (*winnerConsumer) Topic() string {
	return "winner"
}

func (*winnerConsumer) Queue() string {
	return "winner_que"
}

func (w *winnerConsumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery, 0)
	safe.GoRoutine(func() { w.consume(chn) })

	return chn
}

func (w *winnerConsumer) consume(chn chan broker.Delivery) {
	obj := winnerConsumer{}
	res := make(map[string]map[string]int64)
	dead := time.After(2 * time.Minute)

	for {
		select {
		case del := <-chn:
			err := del.Decode(&obj)
			assert.Nil(err)

			key := obj.Data.DemandName
			val, ok := res[key]
			if !ok {
				val = make(map[string]int64)
			}

			val["imp_bid"] += obj.Data.Price
			val["win"] += 1

			res[key] = val

			if len(res) > limit {
				err := obj.flush()
				if err == nil {
					del.Ack(true)
				} else {
					del.Nack(true, true)

				}
				res = make(map[string]map[string]int64)
			}
		case <-dead:
			obj.flush()
			res = make(map[string]map[string]int64)
			dead = time.After(2 * time.Minute)
		}
	}
}

func (w *winnerConsumer) flush() error {
	return nil
}

func (w *winnerConsumer) Initialize() {
	broker.RegisterConsumer(w)
}

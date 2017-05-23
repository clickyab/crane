package workers

import (
	"fmt"
	"strconv"
	"time"

	"services/assert"
	"services/broker"
	"services/safe"
)
func impressionMapper(m impResponse) sourceSupplier {
	return sourceSupplier{
		Time:     m.Data.Time,
		Supplier: m.Data.Source.Supplier.Name,
		Request:  1,
		Source:   m.Data.Source.Name,
	}
}

type sourceSupplier struct {
	Time     time.Time
	Request  int64
	Slot     int
	Bid      int64
	ShowBid  int64
	Supplier string
	Source   string
}

var impressionChan = make(chan Impression)

type impResponse struct {
	Data struct {
		Time   time.Time `json:"time"`
		Source struct {
			Name     string `json:"name"`
			Supplier struct {
				Name string
			} `json:"supplier"`
		} `json:"source"`
	} `json:"data"`
}

type impressionConsumer struct {
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

func (*impressionConsumer) Topic() string {
	return "show"
}

func (*impressionConsumer) Queue() string {
	panic("implement me")
}

func (s *impressionConsumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery, 0)
	safe.GoRoutine(func() { s.fillChan(chn) })
	return chn

}

func (s *impressionConsumer) fillChan(chn chan broker.Delivery) {
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

			key := genID() fmt.Sprintf("%d-%s-%s", hours, obj.Data.Supplier, obj.Data.Publisher)
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
func (s *impressionConsumer) flush(data map[string]map[string]int64) error {
	//q := "INSERT  ON DUPLICATE UPDATE "
	return nil
}


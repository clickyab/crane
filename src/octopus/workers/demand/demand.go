package demand

import (
	"fmt"
	"services/assert"
	"services/broker"
	"services/safe"
	"strconv"
	"time"
)

const limit = 1000

type Acknowledger interface {
	Ack(multiple bool) error
	// Nack negatively acknowledge the delivery of message(s) identified by the delivery tag from either the client or server.
	Nack(multiple, requeue bool) error
	// Reject delegates a negatively acknowledgement through the Acknowledger interface.
	Reject(requeue bool) error
}

type Ad struct {
	Winner_cpm int64 `json:"winner_cpm"`
}

// winnerConsumer asd
type demandConsumer struct {
	Data struct {
		Impression struct {
			Source struct {
				Name     string `json:"name"`
				Supplier struct {
					Name string `json:"name"`
				} `json:"supplier"`
			} `json:"source"`
		} `json:"impression"`
		Demand struct {
			Name string `json:"name"`
		} `json:"dem"`
		Ads  map[string]Ad `json:"ads"`
		Time string        `json:"time"`
	} `json:"data"`
}

func (*demandConsumer) Topic() string {
	return "winner"
}

func (*demandConsumer) Queue() string {
	return "winner_que"
}

func (s *demandConsumer) Consume() chan<- broker.Delivery {
	chn := make(chan broker.Delivery, 0)
	safe.GoRoutine(func() { s.fillChan(chn) })
	return chn

}
func (s *demandConsumer) fillChan(chn chan broker.Delivery) {
	obj := demandConsumer{}
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
			i, err := strconv.ParseInt(obj.Data.Time, 10, 0)
			assert.Nil(err)
			t1 := time.Unix(i, 0)
			layout := "2006-01-02T15:04:05.000Z"
			str := "2017-03-21T00:00:00.000Z"
			t, err := time.Parse(layout, str)

			if err != nil {
				fmt.Println(err)
			}
			assert.Nil(err)
			hours := t1.Sub(t).Hours() + 1
			key := fmt.Sprintf("%s-%d-%s-%s-%s", "supplier_demand_source", hours, obj.Data.Impression.Source.Supplier.Name, obj.Data.Demand.Name, obj.Data.Impression.Source.Name)
			val, ok := ins[key]
			if !ok {
				val = make(map[string]int64)
			}

			// increment
			val["request"] += 1
			val["impression"] += int64(len(obj.Data.Ads))
			val["show"] += 1
			for a := range obj.Data.Ads {
				val["show_bid"] += obj.Data.Ads[a].Winner_cpm
			}

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

func (s *demandConsumer) flush(data map[string]map[string]int64) error {
	q := "INSERT  ON DUPLICATE UPDATE "
}

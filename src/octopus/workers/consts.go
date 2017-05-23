package workers

import (
	"fmt"
	"strconv"
	"time"

	"services/assert"
	"services/safe"
)

// TODO get this from config
const limit = 1000

type Acknowledger interface {
	Ack(multiple bool) error
	// Nack negatively acknowledge the delivery of message(s) identified by the delivery tag from either the client or server.
	Nack(multiple, requeue bool) error
	// Reject delegates a negatively acknowledgement through the Acknowledger interface.
	Reject(requeue bool) error
}

func factID(tm time.Time) int {

	layout := "2006-01-02T15:04:05.000Z"
	str := "2017-03-21T00:00:00.000Z"
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(err)
	return int(tm.Sub(t).Hours()) + 1

}

var dataChannel = make(chan tableModel)

func manager() {
	supDemSrcTable := make(map[string]tableModel)
	supSrcTable := make(map[string]tableModel)
	h := time.After(2 * time.Minute)
	var counter = 0
	var ack Acknowledger
	for {
		select {
		case p := <-dataChannel:
			supDemSrcKey := fmt.Sprint(p.Time, p.Supplier, p.Source, p.Demand)
			m, mk := supDemSrcTable[supDemSrcKey]
			if !mk {
				supDemSrcTable[supDemSrcKey] = p

			} else {
				m = aggregator(m, p)
				supDemSrcTable[supDemSrcKey] = m
			}

			supSrcTableKey := fmt.Sprint(p.Time, p.Supplier, p.Source)
			u, uk := supSrcTable[supSrcTableKey]
			if !uk {

				supSrcTable[supSrcTableKey] = p

			} else {
				u = aggregator(u, p)
				supSrcTable[supSrcTableKey] = u
			}
			counter++
			if counter > limit {
				err := flush(supDemSrcTable, supSrcTable)
				if err == nil {
					ack.Ack(true)
				} else {
					ack.Nack(true, true)

				}
				counter = 0
				supDemSrcTable = make(map[string]tableModel)
				supSrcTable = make(map[string]tableModel)
			}
		case <-h:
			err := flush(supDemSrcTable, supSrcTable)
			if err == nil {
				ack.Ack(true)
			} else {
				ack.Nack(true, true)

			}
			counter = 0
			supDemSrcTable = make(map[string]tableModel)
			supSrcTable = make(map[string]tableModel)
			h = time.After(2 * time.Minute)

		}
	}
}

func init() {
	safe.GoRoutine(func(){
		manager()
	})
}

func timestampToTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 0)
	assert.Nil(err)
	return time.Unix(i, 0)
}

func aggregator(a tableModel, b tableModel) tableModel {
	res := tableModel{}
	res.ShowBid = a.ShowBid + b.ShowBid
	res.Show = a.Show + b.Show
	res.Request = a.Request + b.Request
	res.Impression = a.Impression + b.Impression
	res.ImpressionBid = a.ImpressionBid + b.ImpressionBid
	res.Win = a.Win + b.Win
	if a.Time != 0 {
		res.Time = a.Time
	} else {
		res.Time = b.Time
	}
	return res
}

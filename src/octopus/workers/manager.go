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
}

func factTableID(tm time.Time) int {

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

func worker() {
	supDemSrcTable := make(map[string]*tableModel)
	supSrcTable := make(map[string]*tableModel)

	//TODO Get this from timeout
	h := time.After(2 * time.Minute)
	var counter = 0
	var ack Acknowledger

	defer func() {
		if ack != nil {
			ack.Nack(true, true)
		}
	}()

	flushAndClean := func() {
		err := flush(supDemSrcTable, supSrcTable)
		if err == nil {
			ack.Ack(true)
		} else {
			ack.Nack(true, true)

		}
		counter = 0
		supDemSrcTable = make(map[string]*tableModel)
		supSrcTable = make(map[string]*tableModel)
	}
	for {
		select {
		case p := <-dataChannel:

			if p.Time == 0 {
				assert.NotNil(nil, "Time should not be equal 0")
			}
			if p.Source == "" || p.Supplier == "" {
				assert.NotNil(nil, "Source and supplier can not be empty")
			}
			ack = p.Acknowledger
			supDemSrcKey := fmt.Sprint(p.Time, p.Supplier, p.Source, p.Demand)
			supDemSrcTable[supDemSrcKey] = aggregate(supDemSrcTable[supDemSrcKey], p)

			if p.Demand != "" {
				supSrcTableKey := fmt.Sprint(p.Time, p.Supplier, p.Source)
				supSrcTable[supSrcTableKey] = aggregate(supSrcTable[supSrcTableKey], p)
			}

			counter++

			if counter > limit {
				flushAndClean()
			}

		case <-h:

			flushAndClean()
		}
	}
}

func init() {

	safe.GoRoutine(func() {
		worker()
	})

}

func timestampToTime(s string) time.Time {

	i, err := strconv.ParseInt(s, 10, 0)
	assert.Nil(err)
	return time.Unix(i, 0)

}

func aggregate(a *tableModel, b tableModel) *tableModel {
	if a == nil {
		a = &tableModel{}
	}
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
	return &res
}

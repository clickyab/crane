package manager

import (
	"fmt"
	"time"

	"services/assert"
	"services/broker"
	"services/config"
	"services/safe"
)

var (
	limit   = config.RegisterInt("octopus.worker.manager.limit", 1000, "the limit for points in manager")
	timeout = config.RegisterDuration("octopus.worker.manager.timeout", time.Minute, "the timeout to flush data")
	epoch   time.Time
)

// FactTableID is a helper function to get the fact table id from time
func FactTableID(tm time.Time) int64 {
	return int64(tm.Sub(epoch).Hours()) + 1
}

// DataChannel is a channel to handle data entry for workers without lock
var DataChannel = make(chan TableModel)

func worker() {
	supDemSrcTable := make(map[string]*TableModel)
	supSrcTable := make(map[string]*TableModel)

	h := time.After(*timeout)
	var counter = 0
	var ack broker.Delivery

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
		supDemSrcTable = make(map[string]*TableModel)
		supSrcTable = make(map[string]*TableModel)
	}
	for {
		select {
		case p := <-DataChannel:

			if p.Time == 0 {
				assert.NotNil(nil, "Time should not be equal 0")
			}
			if p.Source == "" || p.Supplier == "" {
				assert.NotNil(nil, "Source and supplier can not be empty")
			}
			ack = *p.Acknowledger
			supDemSrcKey := fmt.Sprint(p.Time, p.Supplier, p.Source, p.Demand)
			supDemSrcTable[supDemSrcKey] = aggregate(supDemSrcTable[supDemSrcKey], p)

			if p.Demand != "" {
				supSrcTableKey := fmt.Sprint(p.Time, p.Supplier, p.Source)
				supSrcTable[supSrcTableKey] = aggregate(supSrcTable[supSrcTableKey], p)
			}

			counter++

			if counter > *limit {
				flushAndClean()
			}

		case <-h:
			flushAndClean()
		}
	}
}

func aggregate(a *TableModel, b TableModel) *TableModel {
	if a == nil {
		a = &TableModel{}
	}
	res := TableModel{}
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

func init() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2017-03-21T00:00:00.000Z"
	var err error
	epoch, err = time.Parse(layout, str)
	assert.Nil(err)

	safe.GoRoutine(func() {
		worker()
	})
}

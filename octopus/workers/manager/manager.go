package manager

import (
	"fmt"
	"time"

	"context"

	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/config"
	"clickyab.com/exchange/services/mysql"
	"clickyab.com/exchange/services/safe"
)

var (
	limit   = config.RegisterInt("octopus.worker.manager.limit", 1000, "the limit for points in manager")
	timeout = config.RegisterDuration("octopus.worker.manager.timeout", time.Minute, "the timeout to flush data")
)

type starter struct {
	channel chan datamodels.TableModel
}

func (s *starter) Initialize() {
	datamodels.RegisterAggregator(s)
	safe.ContinuesGoRoutine(func(context.CancelFunc) {
		s.worker()
	}, time.Second)
}

func (s *starter) Channel() chan<- datamodels.TableModel {
	return s.channel
}

func (s *starter) worker() {
	supDemSrcTable := make(map[string]*datamodels.TableModel)
	supSrcTable := make(map[string]*datamodels.TableModel)

	t := *timeout
	if t < 10*time.Second {
		t = 10 * time.Second
	}
	var counter = 0
	var ack datamodels.Acknowledger

	defer func() {
		if ack != nil {
			// Make sure the packet is rejected to prevent another requeue of an invalid
			// job
			assert.Nil(ack.Reject(false))
		}
	}()

	flushAndClean := func() {
		err := flush(supDemSrcTable, supSrcTable)
		if ack != nil {
			if err == nil {
				assert.Nil(ack.Ack(true))
			} else {
				assert.Nil(ack.Nack(true, true))
			}
		}
		ack = nil
		counter = 0
		supDemSrcTable = make(map[string]*datamodels.TableModel)
		supSrcTable = make(map[string]*datamodels.TableModel)
	}
	ticker := time.NewTicker(t)

bigLoop:
	for {
		select {
		case p := <-s.channel:

			ack = p.Acknowledger
			if p.Time == 0 {
				//assert.NotNil(nil, "Time should not be equal 0")
				assert.Nil(ack.Reject(false))
				continue bigLoop
			}
			if p.Source == "" || p.Supplier == "" {
				//assert.NotNil(nil, "Source and supplier can not be empty")
				assert.Nil(ack.Reject(false))
				continue bigLoop
			}

			supSrcTableKey := fmt.Sprint(p.Time, p.Supplier, p.Source)
			supSrcTable[supSrcTableKey] = aggregate(supSrcTable[supSrcTableKey], p)

			if p.Demand != "" {
				supDemSrcKey := fmt.Sprint(p.Time, p.Supplier, p.Source, p.Demand)
				supDemSrcTable[supDemSrcKey] = aggregate(supDemSrcTable[supDemSrcKey], p)
			}

			counter++

			if counter > *limit {
				flushAndClean()
			}

		case <-ticker.C:
			flushAndClean()
		}
	}
}

func aggregate(a *datamodels.TableModel, b datamodels.TableModel) *datamodels.TableModel {
	if a == nil {
		return &b
	}

	assert.True(a.Time == b.Time, "[BUG] times are not same")

	a.RequestInCount += b.RequestInCount
	a.RequestOutCount += b.RequestOutCount
	a.ImpressionInCount += b.ImpressionInCount
	a.ImpressionOutCount += b.ImpressionOutCount
	a.WinCount += b.WinCount
	a.DeliverCount += b.DeliverCount
	a.WinBid += b.WinBid
	a.DeliverBid += b.DeliverBid
	a.Profit += b.Profit

	return a
}

func init() {
	//make sure worker start after mysql
	mysql.Register(&starter{channel: make(chan datamodels.TableModel)})
}

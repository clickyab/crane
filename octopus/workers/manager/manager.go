package manager

import (
	"fmt"
	"time"

	"context"

	"clickyab.com/exchange/octopus/workers/internal/datamodels"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/safe"
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
	var allAck = make(map[string]datamodels.Acknowledger)

	defer func() {
		if e := recover(); e != nil {
			for i := range allAck {
				// Make sure the packet is rejected to prevent another requeue of an invalid
				// job
				assert.Nil(allAck[i].Reject(false))

			}
		}
	}()

	flushAndClean := func() {
		err := flush(supDemSrcTable, supSrcTable)
		for i := range allAck {
			if err == nil {
				assert.Nil(allAck[i].Ack(true))
			} else {
				assert.Nil(allAck[i].Nack(true, true))
			}
		}
		allAck = make(map[string]datamodels.Acknowledger)
		counter = 0
		supDemSrcTable = make(map[string]*datamodels.TableModel)
		supSrcTable = make(map[string]*datamodels.TableModel)
	}
	ticker := time.NewTicker(t)

bigLoop:
	for {
		select {
		case p := <-s.channel:

			allAck[p.WorkerID] = p.Acknowledger
			if p.Time == 0 {
				//assert.NotNil(nil, "Time should not be equal 0")
				assert.Nil(p.Acknowledger.Reject(false))
				continue bigLoop
			}
			if p.Source == "" || p.Supplier == "" {
				//assert.NotNil(nil, "Source and supplier can not be empty")
				assert.Nil(p.Acknowledger.Reject(false))
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
	a.AdOutCount += b.AdOutCount
	a.DeliverCount += b.DeliverCount
	a.AdOutBid += b.AdOutBid
	a.DeliverBid += b.DeliverBid
	a.Profit += b.Profit
	a.AdInCount += b.AdInCount

	return a
}

func init() {
	//make sure worker start after mysql
	mysql.Register(&starter{channel: make(chan datamodels.TableModel)})
}

package impression

import (
	"context"
	"time"

	"clickyab.com/crane/models/ads"
	m "clickyab.com/crane/workers/models"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/xlog"
)

type consumer struct {
}

func (c *consumer) Topic() string {
	return "#.impression"
}

func (c *consumer) Queue() string {
	return "core-impression"
}

func (c *consumer) Consume(ctx context.Context) chan<- broker.Delivery {
	ch := make(chan broker.Delivery)
	safe.GoRoutine(ctx, func() {
	bigLoop:
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-ch:
				data := job{}
				err := d.Decode(&data)
				if err != nil {
					xlog.GetWithError(ctx, err).Error("impression reject")
					assert.Nil(d.Reject(false))
					continue bigLoop
				}
				if err := data.process(ctx); err != nil {
					xlog.GetWithError(ctx, err).Error("impression nack")
					assert.Nil(d.Nack(false, false))
					continue bigLoop
				}
				assert.Nil(d.Ack(false))
			}
		}
	})
	return ch
}

// NewConsumer return a new consumer
func NewConsumer() broker.Consumer {
	return &consumer{}
}

func init() {
	go bulkInsert()
}

func bulkInsert() {

	imps := make([]m.Impression, 0)
	t := time.After(bulkTime.Duration())
	for {
		select {
		case <-t:
			if len(imps) > 0 {
				err := ads.AddMultiImpression(imps...)
				if err != nil {
					xlog.GetWithError(context.Background(), err)
				}
				imps = make([]m.Impression, 0)
			}
			t = time.After(bulkTime.Duration())

		case s := <-impressions:
			imps = append(imps, s)
			if len(imps) > bulkCount.Int() {
				err := ads.AddMultiImpression(imps...)
				if err != nil {
					xlog.GetWithError(context.Background(), err)
				}
				imps = make([]m.Impression, 0)
			}
			t = time.After(bulkTime.Duration())
		}
	}
}

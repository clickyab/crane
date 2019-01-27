package show

import (
	"context"

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
					xlog.GetWithError(ctx, err).Error("show reject")
					assert.Nil(d.Reject(false))
					continue bigLoop
				}
				if err := data.process(ctx); err != nil {
					xlog.GetWithError(ctx, err).Error("show nack")
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

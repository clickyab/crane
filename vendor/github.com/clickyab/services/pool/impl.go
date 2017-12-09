package pool

import (
	"context"
	"time"

	"sync/atomic"

	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
)

type pl struct {
	loader Loader
	exp    time.Duration
	driver Driver
	retry  int

	started int64
	fail    int

	notify chan time.Time
}

func (p *pl) Notify() <-chan time.Time {
	return p.notify
}

func (p *pl) All() map[string]kv.Serializable {
	return p.driver.All()
}

func (p *pl) Start(ctx context.Context) context.Context {
	if !atomic.CompareAndSwapInt64(&p.started, 0, 1) {
		logrus.Panic("already started")
	}
	return safe.ContinuesGoRoutine(ctx, func(cnl context.CancelFunc) {
		data, err := p.loader(ctx)
		if err != nil {
			p.fail++
			if p.fail > p.retry {
				cnl()
				atomic.SwapInt64(&p.started, 0)
			}
			return
		}
		// There is no need to lock here. implementation decide if it require a lock or not
		err = p.driver.Store(data, time.Duration(p.retry)*p.exp)
		if err != nil {
			p.fail++
		} else {
			p.fail = 0
			select {
			case p.notify <- time.Now():
			default:
			}
		}
	}, p.exp)
}

func (p *pl) Get(s string, data kv.Serializable) (kv.Serializable, error) {
	return p.driver.Fetch(s, data)
}

// NewPool return a new pool object, must start it and watch for ending context
func NewPool(loader Loader, driver Driver, exp time.Duration, retry int) Interface {
	return &pl{
		loader: loader,
		driver: driver,
		exp:    exp,
		retry:  retry,
		notify: make(chan time.Time, 10),
	}
}

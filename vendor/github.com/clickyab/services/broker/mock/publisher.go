package mock

import (
	"sync"

	"github.com/clickyab/services/broker"
)

var (
	b = &chnBroker{}
)

// chnBroker is a mock for publisher interface with ability to
// handle publishers queue
type chnBroker struct {
	lock sync.RWMutex
	out  []chan broker.Job
}

// Publish is here for satisfy the interface
func (p *chnBroker) Publish(j broker.Job) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for i := range p.out {
		// Not blocking
		select {
		case p.out[i] <- j:
		default:
		}
	}
}

// getChannel return a channel for getting messages on
func (p *chnBroker) getChannel(size int) <-chan broker.Job {
	p.lock.Lock()
	defer p.lock.Unlock()

	c := make(chan broker.Job, size)
	p.out = append(p.out, c)
	return c
}

func (p *chnBroker) RegisterConsumer(consumer broker.Consumer) error {
	return nil
}

// GetChannel return a channel for getting messages on
func GetChannel(size int) <-chan broker.Job {
	return b.getChannel(size)
}

// GetChannelBroker return the current active broker
func GetChannelBroker() broker.Publisher {
	return b
}

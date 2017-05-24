package mock

import (
	"sync"
	"time"
)

type channelStore struct {
	lock sync.Mutex
	c    map[string]chan string
}

func (c *channelStore) makeChan(key string) chan string {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ch, ok := c.c[key]; ok {
		return ch
	}

	c.c[key] = make(chan string, 1)
	return c.c[key]
}

func (c *channelStore) closeChan(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.c[key]; ok {
		delete(c.c, key)
	}
}

// Push data in the store
func (c *channelStore) Push(key string, value string, t time.Duration) {
	ch := c.makeChan(key)
	ch <- value
	go func() {
		// after the time is reached or the channel is open for write again
		defer c.closeChan(key)
		select {
		case <-time.After(t):
		case ch <- "":
		}
	}()
}

// Pop and remove data from store, its blocking pop
func (c *channelStore) Pop(key string, t time.Duration) (string, bool) {
	defer c.closeChan(key)
	select {
	case <-time.After(t):
		return "", false
	case data := <-c.makeChan(key):
		return data, true
	}
}

// NewMockChannelStore return a new mock store for testing
func NewMockChannelStore() interface {
	// Push data in the store
	Push(string, string, time.Duration)
	// Pop and remove data from store, its blocking pop
	Pop(string, time.Duration) (string, bool)
} {
	return &channelStore{
		c: make(map[string]chan string),
	}
}

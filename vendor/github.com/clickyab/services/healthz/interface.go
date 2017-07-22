package healthz

import (
	"context"
	"sync"
)

// Interface is the checker interface
type Interface interface {
	// Check must return an error if the health is not ok
	Healthy(context.Context) error
}

var (
	all  []Interface
	lock sync.RWMutex
)

// Register add a new health check service to system
func Register(checker ...Interface) {
	lock.Lock()
	defer lock.Unlock()

	all = append(all, checker...)
}

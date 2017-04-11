package initializer

import (
	"context"
	"time"
)

// Interface is the type to call early on system initialize call
type Interface interface {
	Initialize(context.Context)
}

type group []Interface

var (
	gr = make(group, 0)
)

// Register a module in initializer
func Register(initializer Interface) {
	gr = append(gr, initializer)
}

// Initialize all modules and return the finalizer function
func Initialize() func() {
	ctx, cnl := context.WithCancel(context.Background())
	for i := range gr {
		gr[i].Initialize(ctx)
	}

	return func() {
		cnl()
		<-time.After(time.Second)
	}
}

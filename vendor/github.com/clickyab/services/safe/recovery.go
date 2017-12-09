package safe

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// RecoverHook is the interface to handle recovery
type RecoverHook interface {
	Recover(error, []byte, ...interface{})
}

var (
	recoverHooks []RecoverHook
	lock         = &sync.RWMutex{}
)

// Register is a way to register a hook to trigger after panic
func Register(hook RecoverHook) {
	lock.Lock()
	defer lock.Unlock()

	recoverHooks = append(recoverHooks, hook)
}

func call(err error, stack []byte, extra ...interface{}) {
	newExtra := make([]interface{}, 0)
	// if there is an function the call it here as call back.
	// not a cool idea but leave it here for now.
	for i := range extra {
		if fn, ok := extra[i].(func()); ok {
			fn()
		} else {
			newExtra = append(newExtra, extra[i])
		}
	}
	go func() {
		lock.RLock()
		defer lock.RUnlock()
		defer func() {
			if e := recover(); e != nil {
				logrus.Error("What? the recover function is panicked!")
				logrus.Error(e)
			}
		}()

		for i := range recoverHooks {
			recoverHooks[i].Recover(err, stack, newExtra...)
		}
	}()
}

package safe

import (
	"context"
	"fmt"
	"runtime/debug"

	"clickyab.com/exchange/services/version"

	"time"

	"github.com/Sirupsen/logrus"
)

func mkTitle(err interface{}, title error, commits int64, short string) error {
	switch err.(type) {
	case string:
		title = fmt.Errorf("[%s, %d] %s", short, commits, err.(string))
	case error:
		title = fmt.Errorf("[%s, %d] %s", short, commits, err.(error).Error())
	case *logrus.Entry:
		title = fmt.Errorf("[%s, %d] %s", short, commits, err.(*logrus.Entry).Message)
	}

	return title
}

// GoRoutine is a safe go routine system with recovery and a way to inform finish of the routine
func GoRoutine(f func(), extra ...interface{}) context.Context {
	ctx, cl := context.WithCancel(context.Background())
	go func() {
		defer cl()
		defer func() {
			if e := recover(); e != nil {
				v := version.GetVersion()
				title := fmt.Errorf("[%s, %d] cannot extract title, the type is %T, value is %v", v.Short, v.Count, e, e)
				err := mkTitle(e, title, v.Count, v.Short)
				stack := debug.Stack()
				call(err, stack, extra...)
			}
		}()

		f()
	}()

	return ctx
}

// ContinuesGoRoutine is a safe go routine system with recovery, its continue after recovery
func ContinuesGoRoutine(f func(context.CancelFunc), delay time.Duration, extra ...interface{}) {
	parent, cnl := context.WithCancel(context.Background())
	go func() {
		for i := 1; ; i++ {
			ctx := GoRoutine(func() { f(cnl) }, extra...)
			select {
			case <-ctx.Done():
				time.Sleep(delay)
				logrus.Debugf("restart the routine for %d time", i)
			case <-parent.Done():
				logrus.Debugf("finalize function and exit")
				return
			}
		}
	}()
}

// Routine is a safe routine system with recovery
func Routine(f func(), extra ...interface{}) {
	defer func() {
		if e := recover(); e != nil {
			v := version.GetVersion()
			title := fmt.Errorf("[%s, %d] cannot extract title, the type is %T, value is %v", v.Short, v.Count, e, e)
			err := mkTitle(e, title, v.Count, v.Short)
			stack := debug.Stack()
			call(err, stack, extra...)
		}
	}()

	f()
}

package safe

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/clickyab/services/version"

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

func actual(a func() error, extra ...interface{}) (res error) {
	defer func() {
		if e := recover(); e != nil {
			v := version.GetVersion()
			title := fmt.Errorf("[%s, %d] cannot extract title, the type is %T, value is %v", v.Short, v.Count, e, e)
			err := mkTitle(e, title, v.Count, v.Short)
			stack := debug.Stack()
			call(err, stack, extra...)
			res = err
		}
	}()
	res = a()
	return
}

// Try retry by fibonacci way the given function
func Try(a func() error, max time.Duration, extra ...interface{}) {
	x, y := 0, 1
	for {
		err := actual(a, extra...)
		if err == nil {
			return
		}
		logrus.Error(err)
		t := time.Duration(x) * time.Second
		if t < max {
			x, y = y, x+y
		}
		time.Sleep(t)
	}

}

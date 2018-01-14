package safe

import (
	"context"
	"testing"
	"time"
)

func TestContinuesGoRoutineNormal(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	tm := time.After(time.Second)
	counter := 0
	ch := make(chan int)
	res := make(chan int)
	go ContinuesGoRoutine(ctx, func(x context.CancelFunc) time.Duration {
		counter += 1
		if counter > 3 {
			t.Logf("DymanicContinuesGoRoutine should run 3 times but ran %d", counter)
			t.Fail()
		}
		ch <- 0
		if c := <-res; c == 1 {
			x()
			return 0
		}
		return time.Millisecond
	})
	for {
		select {
		case <-ch:
			if counter == 3 {
				res <- 1
				return
			}
			res <- 0
		case <-tm:
			t.Log("something wrong with DymanicContinuesGoRoutine")
			t.Fail()
		}
	}
}

func TestContinuesGoRoutineWithCancel(t *testing.T) {
	ctx, cl := context.WithCancel(context.Background())
	tm := time.After(time.Second)
	counter := 0
	ch := make(chan int)
	go ContinuesGoRoutine(ctx, func(x context.CancelFunc) time.Duration {
		counter += 1
		if counter > 3 {
			t.Logf("DymanicContinuesGoRoutine should run 3 times but ran %d", counter)
			t.Fail()
		}
		ch <- 0

		return time.Millisecond
	})
	for {
		select {
		case <-ch:
			if counter == 3 {
				cl()
				return
			}
		case <-tm:
			t.Log("something wrong with DymanicContinuesGoRoutine")
			t.Fail()
		}
	}
}

// utils.go
package main

import (
	"runtime"

	"github.com/gobuild/log"
)

const (
	cYELLOW = "33"
	cGREEN  = "32"
	cPURPLE = "35"
)

func cPrintf(ansiColor string, format string, args ...interface{}) {
	if runtime.GOOS != windows {
		format = "\033[" + ansiColor + "m" + format + "\033[0m"
	}
	log.Printf(format, args...)
}

//func groupKill(cmd *exec.Cmd, signal string) (err error) {
//	log.Println("\033[33mprogram terminated\033[0m")
//	var pid, pgid int
//	if cmd.Process != nil {
//		pid = cmd.Process.Pid
//		sess := sh.NewSession()
//		if *verbose {
//			sess.ShowCMD = true
//		}
//		c := sess.Command("/bin/ps", "-o", "pgid", "-p", strconv.Itoa(pid)).Command("sed", "-n", "2,$p")
//		var out []byte
//		out, err = c.Output()
//		if err != nil {
//			return
//		}
//		_, err = fmt.Sscanf(string(out), "%d", &pgid)
//		if err != nil {
//			return
//		}
//		err = sess.Command("pkill", "-"+signal, "--pgroup", strconv.Itoa(pgid)).Run()
//	}
//	return
//}

//func GetFunctionName(i interface{}) string {
//	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
//}

func chanGo(f func() error) chan error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return ch
}

//// block event after specified duration
//func delayEvent(event chan *fsnotify.FileEvent, notifyDelay time.Duration) {
//	for {
//		select {
//		case <-event: //filterEvent:
//			continue
//		case <-time.After(notifyDelay):
//			return
//		}
//	}
//}

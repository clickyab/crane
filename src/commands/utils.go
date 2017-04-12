package commands

import (
	"os"
	"os/signal"
	"syscall"
)

const (
	// AppName the application name
	AppName string = "exchange"
	// Organization the organization name
	Organization = "clickyab"
	// Prefix the prefix for config loader from env
	Prefix = "EXC"
)

// WaitExitSignal get os signal
func WaitExitSignal() os.Signal {
	quit := make(chan os.Signal, 6)
	signal.Notify(quit, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT)
	return <-quit
}

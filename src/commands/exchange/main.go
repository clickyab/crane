package main

import (
	"commands"
	"config"
	"services/initializer"

	"github.com/Sirupsen/logrus"
)

func main() {
	config.Initialize()
	defer initializer.Initialize()()

	go runStatusServer()

	sig := commands.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

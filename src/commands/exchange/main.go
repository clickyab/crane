package main

import (
	"commands"
	"services/config"
	"services/initializer"

	"github.com/Sirupsen/logrus"
)

func main() {
	config.Initialize()
	defer initializer.Initialize()()

	sig := commands.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

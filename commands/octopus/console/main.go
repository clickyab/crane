package main

import (
	"clickyab.com/exchange/commands"
	_ "clickyab.com/exchange/octopus/console/routes"
	"clickyab.com/exchange/services/config"

	"clickyab.com/exchange/services/initializer"

	_ "clickyab.com/exchange/services/eav/redis"
	"clickyab.com/exchange/services/shell"
	"github.com/Sirupsen/logrus"
)

func main() {
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultConfig())
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

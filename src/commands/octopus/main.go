package main

import (
	"commands"
	_ "octopus/demands"
	_ "octopus/router/restful"
	_ "services/broker/selector"
	"services/config"
	_ "services/eav/redis"
	"services/initializer"
	_ "services/statistic/redis"
	_ "services/store/redis"

	"github.com/Sirupsen/logrus"
)

func main() {
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultConfig())
	defer initializer.Initialize()()

	sig := commands.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

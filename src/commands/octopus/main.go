package main

import (
	"commands"
	_ "octopus/demands"
	_ "octopus/router"
	_ "services/broker/selector"
	"services/config"
	_ "services/dset/redis"
	_ "services/eav/redis"
	"services/initializer"
	_ "services/statistic/redis"
	_ "services/store/redis"

	"services/dlock"
	"services/dlock/mock"

	"github.com/Sirupsen/logrus"
)

func main() {

	// TODO : after implementing dlock backend remove the next line
	dlock.Register(mock.NewMockDistributedLocker)

	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultConfig())
	defer initializer.Initialize()()

	sig := commands.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

package main

import (
	"clickyab.com/exchange/commands"
	_ "clickyab.com/exchange/octopus/demands"
	_ "clickyab.com/exchange/octopus/router"
	_ "clickyab.com/exchange/services/broker/selector"
	"clickyab.com/exchange/services/config"
	_ "clickyab.com/exchange/services/dset/redis"
	_ "clickyab.com/exchange/services/eav/redis"
	"clickyab.com/exchange/services/initializer"
	_ "clickyab.com/exchange/services/statistic/redis"
	_ "clickyab.com/exchange/services/store/redis"

	"clickyab.com/exchange/services/dlock"
	"clickyab.com/exchange/services/dlock/mock"

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
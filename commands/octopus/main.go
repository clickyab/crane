package main

import (
	"clickyab.com/exchange/commands"
	_ "clickyab.com/exchange/octopus/demands"
	_ "clickyab.com/exchange/octopus/router"
	_ "github.com/clickyab/services/broker/selector"
	"github.com/clickyab/services/config"
	_ "github.com/clickyab/services/dset/redis"
	_ "github.com/clickyab/services/eav/redis"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	_ "github.com/clickyab/services/statistic/redis"
	_ "github.com/clickyab/services/store/redis"

	"github.com/clickyab/services/dlock"
	"github.com/clickyab/services/dlock/mock"

	_ "github.com/clickyab/services/mysql/connection/mysql"
	"github.com/Sirupsen/logrus"
)

func main() {

	// TODO : after implementing dlock backend remove the next line
	dlock.Register(mock.NewMockDistributedLocker)

	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultConfig())
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

package main

import (
	"os"

	"clickyab.com/crane/commands"
	_ "clickyab.com/crane/demand/layers/input/common"
	_ "clickyab.com/crane/demand/layers/input/ortb"
	_ "github.com/clickyab/services/broker/selector"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	_ "github.com/clickyab/services/kv/redis"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	"github.com/clickyab/services/shell"
	_ "github.com/clickyab/services/slack"
	"github.com/sirupsen/logrus"
)

func main() {
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultConfig())
	config.DumpConfig(os.Stdout)

	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

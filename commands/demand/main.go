package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"clickyab.com/crane/commands"
	_ "clickyab.com/crane/demand/layers/input/allads"
	_ "clickyab.com/crane/demand/layers/input/ortb"
	_ "clickyab.com/crane/demand/layers/input/ortb/stream"
	_ "clickyab.com/crane/demand/layers/output/statics"
	_ "clickyab.com/crane/supplier/layers/exam"
	_ "clickyab.com/crane/supplier/middleware/user"
	_ "github.com/clickyab/services/broker/selector"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	_ "github.com/clickyab/services/kv/redis"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	"github.com/clickyab/services/shell"
)

func main() {
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultConfig())
	config.DumpConfig(os.Stdout)
	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

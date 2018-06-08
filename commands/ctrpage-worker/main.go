package main

import (
	"os"

	"clickyab.com/crane/commands"
	_ "clickyab.com/crane/models/ads/statistics/locationctr"
	_ "clickyab.com/crane/models/pages"
	_ "clickyab.com/crane/models/seats"
	"clickyab.com/crane/workers/ctrpage"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	_ "github.com/clickyab/services/broker/selector"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	_ "github.com/clickyab/services/kv/redis"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
)

func main() {
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultWorkerConfig())
	config.DumpConfig(os.Stdout)

	defer initializer.Initialize()()

	assert.Nil(broker.RegisterConsumer(ctrpage.NewConsumer()))

	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}
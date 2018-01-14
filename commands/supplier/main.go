package main

import (
	"os"

	"clickyab.com/crane/commands"
	_ "clickyab.com/crane/supplier/layer/app"
	_ "clickyab.com/crane/supplier/layer/video"
	_ "clickyab.com/crane/supplier/layer/web"
	// CORS is required for supplier
	_ "clickyab.com/crane/supplier/middleware/cors"
	"github.com/clickyab/services/assert"
	_ "github.com/clickyab/services/broker/selector"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	_ "github.com/clickyab/services/kv/redis"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
)

func main() {
	// its important for supplier to have no mount point. since it need to handle some BC routes
	d := commands.DefaultConfig()
	d.Add("", "services.framework.controller.mount_point", "")
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, d)
	config.DumpConfig(os.Stdout)

	assert.True(config.GetString("services.framework.controller.mount_point") == "", "do not set end point for this app")

	defer initializer.Initialize()()

	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

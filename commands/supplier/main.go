package main

import (
	"os"

	"clickyab.com/crane/commands"
	_ "clickyab.com/crane/supplier/layers/app"
	_ "clickyab.com/crane/supplier/layers/asset"
	_ "clickyab.com/crane/supplier/layers/native"
	_ "clickyab.com/crane/supplier/layers/output/statics"
	_ "clickyab.com/crane/supplier/layers/video"
	_ "clickyab.com/crane/supplier/layers/web"

	_ "clickyab.com/crane/models/staticseat"
	// CORS is required for supplier
	_ "clickyab.com/crane/supplier/middleware/cors"
	_ "clickyab.com/crane/supplier/middleware/user"
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

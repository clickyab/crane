package main

import (
	"clickyab.com/crane/commands"

	// CORS is required for supplier
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
)

func main() {

	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, config.NewDescriptiveLayer())

	assert.True(config.GetString("services.framework.controller.mount_point") == "", "do not set end point for this app")

	defer initializer.Initialize()()
	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

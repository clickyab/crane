package main

import (
	"os"

	"clickyab.com/crane/commands"

	// CORS is required for supplier
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
)

func main() {

	if _, err := os.Stat(os.Getenv("CRN_CLICKYAB_CERT")); os.IsNotExist(err) {
		logrus.Fatal(`certificate not found!
export CRN_CLICKYAB_CERT=/absolute/path/to/cert/file`)

	}

	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, config.NewDescriptiveLayer())

	assert.True(config.GetString("services.framework.controller.mount_point") == "", "do not set end point for this app")

	defer initializer.Initialize()()
	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

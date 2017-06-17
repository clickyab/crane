package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/ip2location"

	"clickyab.com/exchange/commands"
	"github.com/clickyab/services/shell"
	"github.com/Sirupsen/logrus"
)

var (
	listenAddress = config.RegisterString("exchange.ip2location.listen", ":8190", "exchange ip2location port")
)

func main() {
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix)
	defer initializer.Initialize()()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Trim(r.URL.String(), "/ ")
		tmp := ip2location.GetAll(ip)
		dec := json.NewEncoder(w)
		assert.Nil(dec.Encode(tmp))
	})
	go func() {
		http.ListenAndServe(*listenAddress, nil)
	}()

	sig := shell.WaitExitSignal()
	logrus.Debugf("%s received, exiting...", sig.String())
}

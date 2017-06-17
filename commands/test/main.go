package main

import (
	"fmt"
	"net/http"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	"github.com/clickyab/services/shell"
)

var (
	port = config.RegisterString("test.config", ":9898", "test config port")
)

func main() {
	config.Initialize("test", "test", "test")
	http.HandleFunc("/", demandDemo)
	fmt.Println(port.String())
	err := http.ListenAndServe(port.String(), nil)
	assert.Nil(err)
	shell.WaitExitSignal()
}

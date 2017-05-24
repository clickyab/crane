package main

import (
	"fmt"
	"net/http"

	"clickyab.com/exchange/commands"
	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/config"
)

var (
	port = config.RegisterString("test.config", ":9898", "test config port")
)

func main() {
	config.Initialize("test", "test", "test")
	http.HandleFunc("/getad", getAdd)
	http.HandleFunc("/", demandDemo)
	fmt.Println(*port)
	err := http.ListenAndServe(*port, nil)
	assert.Nil(err)
	commands.WaitExitSignal()
}

package main

import (
	"fmt"
	"net/http"

	"clickyab.com/exchange/commands"
	"clickyab.com/exchange/services/config"
)

var (
	port = config.RegisterString("test.config", ":3500", "desc")
)

func main() {
	config.Initialize("test", "test", "test")
	http.HandleFunc("/start", getSupplierDemo)
	http.HandleFunc("/send", postSupplierDemo)
	fmt.Println(*port)
	http.ListenAndServe(*port, nil)
	commands.WaitExitSignal()
}

package main

import (
	"fmt"
	"net/http"

	"clickyab.com/exchange/services/config"
	"clickyab.com/exchange/services/shell"
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
	shell.WaitExitSignal()
}

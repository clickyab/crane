package main

import (
	"fmt"
	"net/http"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/shell"
)

var (
	port = config.RegisterString("test.config", ":3500", "desc")
)

func main() {
	config.Initialize("test", "test", "test")
	http.HandleFunc("/start", getSupplierDemo)
	http.HandleFunc("/send", postSupplierDemo)
	fmt.Println(port.String())
	http.ListenAndServe(port.String(), nil)
	shell.WaitExitSignal()
}

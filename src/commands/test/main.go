package main

import (
	"commands"
	"net/http"
	"services/config"

	"fmt"
	"services/assert"
)

var (
	port = config.RegisterString("test.config", ":3412")
)

func main() {
	config.Initialize("test", "test", "test")
	http.HandleFunc("/getad", getAdd)
	http.HandleFunc("/", demandDemo)
	fmt.Println(*port)
	err := http.ListenAndServe(*port, nil)
	fmt.Println("asd", err)
	assert.Nil(err)
	commands.WaitExitSignal()
}

package main

import (
	"context"
	"core"
	"net/http"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

func runStatusServer() {
	mux := xmux.New()
	core.Mount(mux)

	http.ListenAndServe(":8080", xhandler.New(context.Background(), mux))
}

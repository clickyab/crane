package framework

import (
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

// Handler is a type for all controllers
type Handler xhandler.HandlerFuncC

// Middleware is a middleware generator
type Middleware func(Handler) Handler

// GlobalMiddleware is the middleware that must be on all routes
type GlobalMiddleware interface {
	Handler(Handler) Handler

	PreRoute() bool
}

// Mux is the simple router interface
type Mux interface {

	// GET is a shortcut for mux.Handle("GET", path, handler)
	GET(string, string, Handler)

	// HEAD is a shortcut for mux.Handle("HEAD", path, handler)
	HEAD(string, string, Handler)

	// OPTIONS is a shortcut for mux.Handle("OPTIONS", path, handler)
	OPTIONS(string, string, Handler)

	// POST is a shortcut for mux.Handle("POST", path, handler)
	POST(string, string, Handler)

	// PUT is a shortcut for mux.Handle("PUT", path, handler)
	PUT(string, string, Handler)

	// PATCH is a shortcut for mux.Handle("PATCH", path, handler)
	PATCH(string, string, Handler)

	// DELETE is a shortcut for mux.Handle("DELETE", path, handler)
	DELETE(string, string, Handler)

	// NewGroup creates a new routes group with the provided path prefix.
	// All routes added to the returned group will have the path prepended.
	NewGroup(string) Mux

	// RootMux return the root mux without any prefix for routes like health
	// currently just for health check route and things like that
	RootMux() *xmux.Mux
}

// Routes the base rote structure
type Routes interface {
	// Routes is for adding new controller
	Routes(Mux)
}

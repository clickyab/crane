package framework

import "github.com/rs/xhandler"

// Handler is a type for all controllers
type Handler xhandler.HandlerFuncC

// Middleware is a middleware generator
type Middleware func(Handler) Handler

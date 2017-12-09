package router

import (
	"path/filepath"

	"github.com/clickyab/services/framework"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

type xmuxer struct {
	path       string
	root       *xmux.Mux
	engine     *xmux.Mux
	group      *xmux.Group
	middleware func(next framework.Handler) framework.Handler
}

func (x *xmuxer) getFunc(handler framework.Handler) xhandler.HandlerFuncC {
	return xhandler.HandlerFuncC(x.middleware(handler))
}

func (x *xmuxer) NewGroup(path string) framework.Mux {
	xm := &xmuxer{
		path:       filepath.Join(x.path, path),
		root:       x.root,
		engine:     nil,
		middleware: x.middleware,
	}
	if x.engine != nil {
		xm.group = x.engine.NewGroup(path)
	} else {
		xm.group = x.group.NewGroup(path)
	}

	return xm
}

func (x *xmuxer) RootMux() *xmux.Mux {
	return x.root
}

func (x *xmuxer) GET(name string, path string, handler framework.Handler) {
	AddRoute(name, filepath.Join(x.path, path))
	if x.engine != nil {
		x.engine.GET(path, x.getFunc(handler))
		return
	}
	x.group.GET(path, x.getFunc(handler))
}

func (x *xmuxer) HEAD(name string, path string, handler framework.Handler) {
	AddRoute(name, filepath.Join(x.path, path))
	if x.engine != nil {
		x.engine.HEAD(path, x.getFunc(handler))
		return
	}
	x.group.HEAD(path, x.getFunc(handler))
}

func (x *xmuxer) OPTIONS(name string, path string, handler framework.Handler) {
	AddRoute(name, filepath.Join(x.path, path))
	if x.engine != nil {
		x.engine.HEAD(path, x.getFunc(handler))
		return
	}
	x.group.HEAD(path, x.getFunc(handler))
}

func (x *xmuxer) POST(name string, path string, handler framework.Handler) {
	AddRoute(name, filepath.Join(x.path, path))
	if x.engine != nil {
		x.engine.POST(path, x.getFunc(handler))
		return
	}
	x.group.POST(path, x.getFunc(handler))
}

func (x *xmuxer) PUT(name string, path string, handler framework.Handler) {
	AddRoute(name, filepath.Join(x.path, path))
	if x.engine != nil {
		x.engine.PUT(path, x.getFunc(handler))
		return
	}
	x.group.PUT(path, x.getFunc(handler))
}

func (x *xmuxer) PATCH(name string, path string, handler framework.Handler) {
	AddRoute(name, filepath.Join(x.path, path))
	if x.engine != nil {
		x.engine.PATCH(path, x.getFunc(handler))
		return
	}
	x.group.PATCH(path, x.getFunc(handler))
}

func (x *xmuxer) DELETE(name string, path string, handler framework.Handler) {
	AddRoute(name, filepath.Join(x.path, path))
	if x.engine != nil {
		x.engine.DELETE(path, x.getFunc(handler))
		return
	}
	x.group.DELETE(path, x.getFunc(handler))
}

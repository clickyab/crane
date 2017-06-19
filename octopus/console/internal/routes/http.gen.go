// Code generated build with router DO NOT EDIT.

package routes

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/middleware"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/initializer"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

// Routes return the route registered with this
func (c *Controller) Routes(r *xmux.Mux, mountPoint string) {

	groupMiddleware := []framework.Middleware{}

	group := r.NewGroup(mountPoint + "/user")

	/* Route {
		"Route": "/login",
		"Method": "POST",
		"Function": "Controller.login",
		"RoutePkg": "routes",
		"RouteMiddleware": null,
		"RouteFuncMiddleware": "",
		"RecType": "Controller",
		"RecName": "c",
		"Payload": "loginPayload",
		"Resource": "",
		"Scope": ""
	} with key 0 */
	m0 := append(groupMiddleware, []framework.Middleware{}...)

	// Make sure payload is the last middleware
	m0 = append(m0, middleware.PayloadUnMarshallerGenerator(loginPayload{}))
	group.POST("/login", xhandler.HandlerFuncC(framework.Mix(c.login, m0...)))
	// End route with key 0

	initializer.DoInitialize(c)
}

func init() {
	router.Register(&Controller{})
}

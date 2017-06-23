// Code generated build with router DO NOT EDIT.

package routes

import (
	"clickyab.com/exchange/octopus/console/user/routes"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/initializer"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
)

// Routes return the route registered with this
func (c *Controller) Routes(r *xmux.Mux, mountPoint string) {

	groupMiddleware := []framework.Middleware{}

	group := r.NewGroup(mountPoint + "/report")

	/* Route {
		"Route": "/demand/:from/:to",
		"Method": "GET",
		"Function": "Controller.demand",
		"RoutePkg": "routes",
		"RouteMiddleware": [
			"routes.Authenticate"
		],
		"RouteFuncMiddleware": "",
		"RecType": "Controller",
		"RecName": "c",
		"Payload": "",
		"Resource": "",
		"Scope": ""
	} with key 0 */
	m0 := append(groupMiddleware, []framework.Middleware{
		routes.Authenticate,
	}...)

	group.GET("/demand/:from/:to", xhandler.HandlerFuncC(framework.Mix(c.demand, m0...)))
	// End route with key 0

	/* Route {
		"Route": "/exchange/:from/:to",
		"Method": "GET",
		"Function": "Controller.exchange",
		"RoutePkg": "routes",
		"RouteMiddleware": [
			"routes.Authenticate"
		],
		"RouteFuncMiddleware": "",
		"RecType": "Controller",
		"RecName": "c",
		"Payload": "",
		"Resource": "",
		"Scope": ""
	} with key 1 */
	m1 := append(groupMiddleware, []framework.Middleware{
		routes.Authenticate,
	}...)

	group.GET("/exchange/:from/:to", xhandler.HandlerFuncC(framework.Mix(c.exchange, m1...)))
	// End route with key 1

	/* Route {
		"Route": "/supplier/:from/:to",
		"Method": "GET",
		"Function": "Controller.supplier",
		"RoutePkg": "routes",
		"RouteMiddleware": [
			"routes.Authenticate"
		],
		"RouteFuncMiddleware": "",
		"RecType": "Controller",
		"RecName": "c",
		"Payload": "",
		"Resource": "",
		"Scope": ""
	} with key 2 */
	m2 := append(groupMiddleware, []framework.Middleware{
		routes.Authenticate,
	}...)

	group.GET("/supplier/:from/:to", xhandler.HandlerFuncC(framework.Mix(c.supplier, m2...)))
	// End route with key 2

	initializer.DoInitialize(c)
}

func init() {
	router.Register(&Controller{})
}

// Code generated build with router DO NOT EDIT.

package user

import (
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/middleware"
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
		"RoutePkg": "user",
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

	/* Route {
		"Route": "/logout",
		"Method": "GET",
		"Function": "Controller.logout",
		"RoutePkg": "user",
		"RouteMiddleware": [
			"Authenticate"
		],
		"RouteFuncMiddleware": "",
		"RecType": "Controller",
		"RecName": "c",
		"Payload": "",
		"Resource": "",
		"Scope": ""
	} with key 1 */
	m1 := append(groupMiddleware, []framework.Middleware{
		Authenticate,
	}...)

	group.GET("/logout", xhandler.HandlerFuncC(framework.Mix(c.logout, m1...)))
	// End route with key 1

	initializer.DoInitialize(c)
}

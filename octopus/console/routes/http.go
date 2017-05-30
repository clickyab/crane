package routes

import (
	"clickyab.com/exchange/services/httplib/controller"
	"gopkg.in/labstack/echo.v3"
)

type initConsole struct{}

func (initConsole) Routes(e *echo.Echo, mountPoint string) {
	e.GET(mountPoint+"/login", loginGet, recovery)
	e.POST(mountPoint+"/login", loginPost, recovery)
	e.GET(mountPoint+"/logout", logout, recovery, auth)
}

func init() {
	controller.Register(&initConsole{})
}

package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	echo "gopkg.in/labstack/echo.v3"

	"github.com/Sirupsen/logrus"
)

// Recovery is the middleware to prevent the panic to crash the app
func Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				ctx.JSON(
					http.StatusInternalServerError,
					struct {
						Error string `json:"error"`
					}{
						Error: http.StatusText(http.StatusInternalServerError),
					},
				)
				stack := debug.Stack()
				dump := "TODO: create dump" //httputil.DumpRequest(ctx.Request(), true)
				data := fmt.Sprintf("Request : \n %s \n\nStack : \n %s", dump, stack)
				logrus.WithField("error", err).Warn(err, data)
				// TODO : use safe package
			}
		}()

		return next(ctx)
	}
}

package routes

import (
	"net/http"

	"clickyab.com/exchange/octopus/console/internal/manager"
	"clickyab.com/exchange/services/eav"
	"clickyab.com/exchange/services/safe"
	"gopkg.in/labstack/echo.v3"
)

const userData = "__user_data__"
const tokenData = "__token__"

func auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("token")
		if token != "" {
			val := eav.NewEavStore(token).SubKey("token")
			if val != "" {
				user, err := manager.NewManager().GetUserByToken(val)
				if err == nil {
					c.Set(userData, user)
					c.Set(tokenData, token)
					return next(c)
				}
			}
		}
		return c.JSON(http.StatusUnauthorized, struct {
			error string
		}{
			error: http.StatusText(http.StatusUnauthorized),
		})
	}
}

// recovery is the middleware to prevent the panic to crash the app
func recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		var failed bool
		safe.Routine(func() {
			err = next(ctx)
		}, func() {
			failed = true
		})
		if failed {
			return ctx.JSON(http.StatusInternalServerError, struct {
				error string
			}{
				error: http.StatusText(http.StatusInternalServerError),
			})
		}
		return err
	}
}

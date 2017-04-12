package middlewares

import (
	"context"
	"net/http"

	"encoding/json"

	"github.com/fzerorubigd/xhandler"
)

// Recovery is the middleware to prevent the panic to crash the app
func Recovery(next xhandler.HandlerFuncC) xhandler.HandlerFuncC {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				resp := struct {
					Error string `json:"error"`
				}{
					Error: http.StatusText(http.StatusInternalServerError),
				}
				enc := json.NewEncoder(w)
				enc.Encode(resp)

			}
		}()

		next(c, w, r)
	}
}

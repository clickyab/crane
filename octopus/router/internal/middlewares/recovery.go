package middlewares

import (
	"context"
	"net/http"

	"encoding/json"

	"clickyab.com/exchange/services/safe"

	"github.com/fzerorubigd/xhandler"
)

// Recovery is the middleware to prevent the panic to crash the app
func Recovery(next xhandler.HandlerFuncC) xhandler.HandlerFuncC {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {

		errResp := func() {
			resp := struct {
				Error string `json:"error"`
			}{
				Error: http.StatusText(http.StatusInternalServerError),
			}
			enc := json.NewEncoder(w)
			enc.Encode(resp)
		}

		safe.Routine(func() { next(c, w, r) }, errResp, r)
	}
}

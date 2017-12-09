package middleware

import (
	"net/http"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/safe"
)

// Recovery is the middleware to prevent the panic to crash the app
func Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var failed bool
		safe.Routine(
			func() {
				next(w, r)
			},
			func() {
				failed = true
			},
			r,
		)
		if failed {
			framework.JSON(
				w,
				http.StatusInternalServerError,
				struct {
					Error string `json:"error"`
				}{
					Error: http.StatusText(http.StatusInternalServerError),
				},
			)
		}
	}
}

package cors

import (
	"context"
	"net/http"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/cors"
)

type middleware struct {
}

func (middleware) PreRoute() bool {
	return true
}

func (middleware) Handler(next framework.Handler) framework.Handler {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"token", "content-type"},
		AllowOriginFunc: func(origin string) bool {
			return true // TODO : write the real code here
		},
	})

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		c.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
			next(ctx, w, r)
		})
	}
}

func init() {
	router.RegisterGlobalMiddleware(&middleware{})
}

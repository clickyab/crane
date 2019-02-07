package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/version"
)

type middleware struct {
	data string
}

func (middleware) PreRoute() bool {
	return true
}

func (m *middleware) Handler(next framework.Handler) framework.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cy-Version", m.data)
		next(ctx, w, r)
	}
}

func init() {
	v := version.GetVersion()
	data := fmt.Sprintf("%s:%d", v.Short, v.Count)
	router.RegisterGlobalMiddleware(&middleware{data: data})
}

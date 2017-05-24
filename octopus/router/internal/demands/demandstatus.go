package demands

import (
	"context"
	"net/http"

	"clickyab.com/exchange/octopus/core"

	"github.com/fzerorubigd/xmux"
)

// Status is the demand status route
func Status(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	name := xmux.Param(ctx, "name")
	demand, err := core.GetDemand(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	demand.Status(ctx, w, r)
}

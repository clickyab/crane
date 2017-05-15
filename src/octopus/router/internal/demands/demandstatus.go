package demands

import (
	"context"
	"net/http"
	"octopus/core"

	"github.com/fzerorubigd/xmux"
)

// DemandStatus is the demand status route
func DemandStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	name := xmux.Param(ctx, "name")
	demand, err := core.GetDemand(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	demand.Status(ctx, w, r)
}

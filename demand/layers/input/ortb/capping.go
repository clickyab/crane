package ortb

import (
	"context"
	"net/http"
	"strconv"

	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
	"github.com/rs/xmux"
)

const cappingURL = "/capping/:adID/:userID/:cappMode"

func applyCapp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	adID := xmux.Param(ctx, "adID")
	userID := xmux.Param(ctx, "userID")
	cappMode := xmux.Param(ctx, "cappMode")
	cappModeInt, err := strconv.Atoi(cappMode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cp := entity.CappingMode(cappModeInt)

	if !cp.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if cp != entity.CappingNone {
		adIDInt, err := strconv.ParseInt(adID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		capping.StoreCapping(cp, userID, adIDInt)
	}

	w.WriteHeader(http.StatusOK)
}

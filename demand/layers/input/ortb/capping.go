package ortb

import (
	"context"
	"net/http"
	"strconv"

	"clickyab.com/crane/demand/capping"
	"clickyab.com/crane/demand/entity"
	"github.com/rs/xmux"
)

const cappingURL = "/capping/:ad_id/:user_id/:capp_mode"

func applyCapp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	adID := xmux.Param(ctx, "ad_id")
	userID := xmux.Param(ctx, "user_id")
	cappMode := xmux.Param(ctx, "capp_mode")
	cappModeInt, err := strconv.Atoi(cappMode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cp := entity.CappingMode(cappModeInt)

	if cp.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("validation failed"))
		return
	}

	if cp != entity.CappingNone {
		adIDInt, err := strconv.ParseInt(adID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid id"))
			return
		}

		capping.StoreCapping(cp, userID, adIDInt)
	}

	w.WriteHeader(http.StatusOK)
}

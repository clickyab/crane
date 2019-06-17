package ortb

import (
	"context"
	"net/http"
	"time"

	"github.com/clickyab/services/kv"
	"github.com/rs/xmux"
)

const balancePath = "/balance/:page"

func balancer(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	k := kv.NewAEAVStore("BLC_"+xmux.Param(ctx, "page"), time.Hour*72)
	k.IncSubKey("c", 1)
}

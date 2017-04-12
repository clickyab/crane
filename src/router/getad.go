package router

import (
	"context"
	"core"
	"encoding/json"
	"net/http"
	"rtb"
	"supliers"

	"entity"

	"github.com/fzerorubigd/xmux"
)

func getAd(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	dec := json.NewEncoder(w)
	key := xmux.Param(ctx, "key")
	imp, err := supliers.GetImpression(key, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		dec.Encode(struct {
			Error error
		}{
			Error: err,
		})
		return
	}
	nCtx, cnl := context.WithCancel(ctx)
	defer cnl()
	ads := core.Call(nCtx, imp)
	var has bool
	for i := range ads {
		if len(ads[i]) > 0 {
			has = true
			break
		}
	}
	if !has {
		w.WriteHeader(http.StatusNoContent)
		dec.Encode(struct {
			Error string
		}{
			Error: "no data",
		})
		return
	}
	res := rtb.SelectCPM(imp, ads)
	has = false
	data := make(map[string]entity.DumbAd)
	for i := range res {
		if res[i] != nil {
			has = true
			data[i] = res[i].Morph("winpixel")
		}
	}

	if !has {
		w.WriteHeader(http.StatusNoContent)
		dec.Encode(struct {
			Error string
		}{
			Error: "no data",
		})
		return
	}

	dec.Encode(data)
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/bsm/openrtb"
	"github.com/clickyab/services/xlog"
)

// Call an openrtp end point
func Call(ctx context.Context, method, url string, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	d, err := json.Marshal(pl)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("failed")
		return nil, err
	}
	buf := bytes.NewReader(d)
	r, err := http.NewRequest(method, url, buf)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("failed")
		return nil, err
	}
	nCtx, _ := context.WithTimeout(ctx, timeout.Duration())
	resp, err := httpClient.Do(r.WithContext(nCtx))
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("failed")
		return nil, err
	}

	bid := openrtb.BidResponse{}
	err = json.NewDecoder(resp.Body).Decode(&bid)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("failed")
		return nil, err
	}
	defer resp.Body.Close()

	return &bid, nil
}

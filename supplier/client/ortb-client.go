package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"clickyab.com/crane/openrtb"
	"github.com/clickyab/services/xlog"
)

// Call an openrtp end point
func Call(ctx context.Context, method, url string, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	bid := &openrtb.BidResponse{}
	if len(pl.Imp) == 0 {
		return bid, nil
	}
	d, err := json.Marshal(pl)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("marshal failed")
		return nil, err
	}
	buf := bytes.NewReader(d)
	r, err := http.NewRequest(method, url, buf)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("request create failed")
		return nil, err
	}
	nCtx, _ := context.WithTimeout(ctx, timeout.Duration())
	resp, err := httpClient.Do(r.WithContext(nCtx))
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("request do failed")
		return nil, err
	}
	// in any case, on error and non-error situation we use the body, so defer close here is a good idea :)
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("invalid status %d, message was %s", resp.StatusCode, string(data))
		xlog.GetWithError(ctx, err).Debug("request do status failed")
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(bid)
	if err != nil {
		xlog.GetWithError(ctx, err).Debug("decode failed")
		return nil, err
	}
	return bid, nil
}

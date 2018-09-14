package client

import (
	"context"
	"fmt"

	"clickyab.com/crane/openrtb/v2.5"
)

// Call an openrtp end point
func Call(ctx context.Context, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	ct, cl := context.WithTimeout(ctx, timeout.Duration())
	res := make(chan *openrtb.BidResponse)
	RequestChannel <- &StreamRequest{
		BidRequest: pl,
		Response:   res,
		Context:    ct,
	}
	select {
	case rs := <-res:
		return rs, nil
	case <-ct.Done():
		cl()
		return nil, fmt.Errorf("timeout excced for: %s", pl.Id)
	}

}

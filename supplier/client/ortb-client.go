package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"

	openrtb "clickyab.com/crane/openrtb/v2.5"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/xlog"
)

var demand = config.RegisterString("crane.supplier.client.mode", "managed", "")

// Call an openrtb end point
func Call(ctx context.Context, url string, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	switch demand.String() {
	case "managed":
		return managed(ctx, pl)
	case "unary":
		return UnaryCall(ctx, pl)
	case "stream":
		return StreamCall(ctx, pl)
	case "rest":
		return RestCall(ctx, url, pl)
	default:
		return RestCall(ctx, url, pl)
	}
}

// RestCall an openrtb end point
func RestCall(ctx context.Context, url string, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
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
	r, err := http.NewRequest("POST", url, buf)
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

func managed(ctx context.Context, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	res := make(chan *openrtb.BidResponse)
	RequestChannel <- &StreamRequest{
		Response:   res,
		BidRequest: pl,
		Context:    ctx,
	}
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context done")
	case rs := <-res:
		if rs == nil {
			return nil, fmt.Errorf("no ad")
		}
		return rs, nil
	}
}

var ucl openrtb.OrtbServiceClient
var recon chan int
var uclock = sync.RWMutex{}

func unaryInit(ctx context.Context) {
	for {
		uclock.Lock()
		recon = make(chan int)
	RC:
		conn, err := grpc.Dial(insecureSever.String(), grpc.WithInsecure())
		if err != nil {
			fmt.Println(fmt.Sprintf("filed to connect: %s", err))
			time.Sleep(time.Second * 2)
			goto RC
		}
		ucl = openrtb.NewOrtbServiceClient(conn)
		uclock.Unlock()
		select {
		case <-recon:
			_ = conn.Close()
		case <-ctx.Done():
			break
		}
	}
}

// UnaryCall an openrtp end point
func UnaryCall(ctx context.Context, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	uclock.RLock()
	defer uclock.RUnlock()
	pl.Token = token.String()
	res, err := ucl.Ortb(ctx, pl)
	if err != nil {
		close(recon)
	}
	return res, err
}

// UnaryCall an openrtp end point
// func UnaryCall(ctx context.Context, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
// 	conn, err := grpc.Dial(insecureSever.String(), grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		_ = conn.Close()
// 	}()
// 	client := openrtb.NewOrtbServiceClient(conn)
// 	pl.Token = token.String()
// 	return client.Ortb(ctx, pl)
// }

// StreamCall an openrtp end point
func StreamCall(ctx context.Context, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	ct, cl := context.WithTimeout(ctx, timeout.Duration())
	defer cl()
	res := make(chan *openrtb.BidResponse)
	pl.Token = token.String()

	RequestChannel <- &StreamRequest{
		BidRequest: pl,
		Response:   res,
		Context:    ct,
	}
	defer func() {
		lock.Lock()
		delete(inprogress, pl.Id)
		lock.Unlock()
	}()
	select {
	case rs := <-res:
		return rs, nil
	case <-ct.Done():
		return nil, fmt.Errorf("timeout excced for: %s", pl.Id)
	}

}

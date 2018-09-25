package main

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/crane/commands/streaming-supplier/client"
	"clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/safe"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
)

type controller struct {
}

// Routes is for registering routes
func (controller) Routes(r framework.Mux) {
	r.POST("ortb", "/ortb", openRTBInput)
}

func init() {
	router.Register(&controller{})
}

var token = config.RegisterString("clickyab.token", "", "")
var timeout = config.RegisterDuration("clickyab.timeout", time.Millisecond*150, "maximum timeout")

// openRTBInput is the route for rtb input layer
//func openRTBInputStream(ct context.Context, w http.ResponseWriter, r *http.Request) {
//	ctx, cl := context.WithTimeout(ct, timeout.Duration())
//	defer cl()
//	//tk := time.Now()
//	payload := &openrtb.BidRequest{}
//	err := jsonpb.Unmarshal(r.Body, payload)
//	if err != nil {
//		logrus.Warn(err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	payload.Token = token.String()
//	w.Header().Set("content-type", "application/json")
//	m := jsonpb.Marshaler{}
//
//	res := &openrtb.BidResponse{
//		Id: payload.Id,
//	}
//	rc := make(chan *openrtb.BidResponse)
//	client.RequestChannel <- &client.StreamRequest{
//		BidRequest: payload,
//		Context:    ctx,
//		Response:   rc,
//	}
//
//	select {
//	case <-ctx.Done():
//		logrus.Infof("timeout exceed for request id: ", payload.Id)
//		assert.Nil(m.Marshal(w, res))
//	case rs := <-rc:
//		assert.Nil(m.Marshal(w, rs))
//
//	}
//
//}

// openRTBInput is the route for rtb input layer
func openRTBInput(ct context.Context, w http.ResponseWriter, r *http.Request) {
	ctx, cl := context.WithTimeout(ct, timeout.Duration())
	defer cl()
	//tk := time.Now()
	payload := &openrtb.BidRequest{}
	err := jsonpb.Unmarshal(r.Body, payload)
	if err != nil {
		logrus.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload.Token = token.String()
	w.Header().Set("content-type", "application/json")
	m := jsonpb.Marshaler{}

	rc := make(chan *openrtb.BidResponse)

	safe.GoRoutine(ctx, func() {
		res, err := client.UnaryCall(ctx, payload)
		if err != nil {
			rc <- nil
			return
		}
		rc <- res
	})
	res := &openrtb.BidResponse{
		Id: payload.Id,
	}
	select {
	case <-ctx.Done():
		logrus.Infof("timeout exceed for request id: ", payload.Id)
	case rs := <-rc:
		if rs != nil {
			res = rs
		}
	}
	assert.Nil(m.Marshal(w, res))

}

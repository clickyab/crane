package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type initRestClient struct {
}

var (
	httpClient *http.Client
	maxIdle    = config.RegisterInt("crane.supplier.max_idle_connection", 30, "maximum idle connection count")
)

func (*initRestClient) Initialize(context.Context) {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdle.Int(),
			MaxIdleConns:        maxIdle.Int() + 1,
		},
	}
}

func init() {
	initializer.Register(&initRestClient{}, 100)
}

var (
	connection    = config.RegisterInt("crane.supplier.stream.connection", 40, "")
	insecureSever = config.RegisterString("crane.supplier.stream.address", "crane-stream:9001", "")
	token         = config.RegisterString("crane.supplier.demand.token", "forbidden", "")
	timeout       = config.RegisterDuration("crane.supplier.timeout", time.Millisecond*150, "maximum timeout")
	// RequestChannel for stream
	RequestChannel = make(chan *StreamRequest, 100000)
	lock           = sync.RWMutex{}
)

// StreamRequest for stream
type StreamRequest struct {
	BidRequest *openrtb.BidRequest
	Context    context.Context
	Response   chan<- *openrtb.BidResponse
}

type initClient struct{}

var inprogress = make(map[string]*StreamRequest)

func (*initClient) Initialize(ctx context.Context) {
	//safe.Try(func() error {
	//	for {
	//		if conn == nil && client == nil {
	//			time.Sleep(time.Second)
	//			continue
	//		}
	//		break
	//	}
	//	if st := conn.GetState(); st != connectivity.Idle &&
	//		st != connectivity.Connecting &&
	//		st != connectivity.Ready {
	//		connlock.Lock()
	//		defer connlock.Unlock()
	//		var err error
	//		conn, err = grpc.Dial(insecureSever.String(), grpc.WithInsecure())
	//		if err != nil {
	//			return err
	//		}
	//		client = openrtb.NewOrtbServiceClient(conn)
	//	}
	//	return fmt.Errorf("next round")
	//
	//}, time.Millisecond*100)
	//safe.GoRoutine(ctx, func() {
	//	connlock.Lock()
	//	defer connlock.Unlock()
	//	var err error
	//	conn, err = grpc.Dial(insecureSever.String(), grpc.WithInsecure())
	//	assert.Nil(err)
	//	client = openrtb.NewOrtbServiceClient(conn)
	//	<-ctx.Done()
	//})

	for i := 0; i < connection.Int(); i++ {
		safe.Try(func() error {
			var cerr = make(chan error)
			conn, err := grpc.Dial(insecureSever.String(), grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer func() { _ = conn.Close() }()
			client := openrtb.NewOrtbServiceClient(conn)
			cl, err := client.OrtbStream(ctx)
			if err != nil {
				return err
			}
			defer func() { _ = cl.CloseSend() }()
			safe.GoRoutine(ctx, func() {
				for {
					rq := <-RequestChannel
					lock.Lock()
					inprogress[rq.BidRequest.Id] = rq
					lock.Unlock()

					err := cl.Send(rq.BidRequest)
					if err != nil {
						cerr <- fmt.Errorf("stream send: %v", err)
					}
				}
			})
			safe.GoRoutine(ctx, func() {
				for {
					rs, err := cl.Recv()
					if err == io.EOF {
						continue
					}
					if err != nil {
						cerr <- fmt.Errorf("stream recv: %v", err)
					}
					lock.RLock()
					if or, ok := inprogress[rs.Id]; ok {
						or.Response <- rs
					}
					lock.RUnlock()
				}
			})
			select {
			case err := <-cerr:
				logrus.Error(err)
				return err
			case <-ctx.Done():
				return nil
			}
		}, time.Nanosecond)
	}

}

func init() {
	initializer.Register(&initClient{}, 100)
}

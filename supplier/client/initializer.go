package client

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	connection    = config.RegisterInt("crane.supplier.stream.connection", 40, "")
	insecureSever = config.RegisterString("crane.supplier.stream.address", "crane-stream:9001", "")
	token         = config.RegisterString("crane.supplier.demand.token", "forbidden", "")
	timeout       = config.RegisterDuration("crane.supplier.timeout", time.Millisecond*150, "maximum timeout")
	// RequestChannel for stream
	RequestChannel = make(chan *StreamRequest, 1000)
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
	for i := 0; i < connection.Int(); i++ {
		safe.Try(func() error {
			var cerr = make(chan error)
			conn, err := grpc.Dial(insecureSever.String(), grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer func() { assert.Nil(conn.Close()) }()
			client := openrtb.NewOrtbServiceClient(conn)
			cl, err := client.OrtbStream(ctx)
			if err != nil {
				return err
			}
			defer func() { assert.Nil(cl.CloseSend()) }()
			safe.GoRoutine(ctx, func() {
				for {
					rq := <-RequestChannel
					go func() {
						<-rq.Context.Done()
						lock.Lock()
						delete(inprogress, rq.BidRequest.Id)
						lock.Unlock()
					}()
					rq.BidRequest.Token = token.String()
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
				logrus.Debug(err)
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

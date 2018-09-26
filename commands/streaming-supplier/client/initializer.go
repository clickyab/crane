package client

import (
	"context"

	"clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	//connection  = config.RegisterInt("clickyab.stream.connection", 40, "")
	secureSever = config.RegisterString("clickyab.stream.address", "stream.clickyab.com:30100", "")
	cert        = config.RegisterString("clickyab.cert", "", "")
	// RequestChannel for stream
	//RequestChannel = make(chan *StreamRequest, 10000)
	//lock           = sync.RWMutex{}
	//inprogress     = make(map[string]*StreamRequest)
)

// UnaryCall an openrtp end point
func UnaryCall(ctx context.Context, pl *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	creads, err := credentials.NewClientTLSFromFile(cert.String(), "")
	if err != nil {
		logrus.Fatal("certificate is not valid")
	}
	defer creads.Clone()
	conn, err := grpc.Dial(secureSever.String(), grpc.WithTransportCredentials(creads))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()
	client := openrtb.NewOrtbServiceClient(conn)
	return client.Ortb(ctx, pl)
}

func init() {
	//initializer.Register(&handler{}, 1000)
}

//
//// StreamRequest for streaming
//type StreamRequest struct {
//	BidRequest *openrtb.BidRequest
//	Context    context.Context
//	Response   chan *openrtb.BidResponse
//}
//
//type handler struct {
//}
//
//func (*handler) Initialize(ct context.Context) {
//	for i := 0; i < connection.Int(); i++ {
//
//		safe.Try(func() error {
//			ctx := ocontext.Context(ct)
//
//			var cerr = make(chan error)
//			creads, err := credentials.NewClientTLSFromFile(cert.String(), "")
//			if err != nil {
//				log.Fatal(err)
//				return err
//
//			}
//			defer creads.Clone()
//			conn, err := grpc.Dial(secureSever.String(), grpc.WithTransportCredentials(creads))
//			if err != nil {
//				logrus.Debugf("Dial: %v", err)
//				return err
//
//			}
//			defer func() { assert.Nil(conn.Close()) }()
//			client := openrtb.NewOrtbServiceClient(conn)
//			cl, err := client.OrtbStream(ctx)
//			if err != nil {
//				logrus.Debugf("client: %v", err)
//				return err
//
//			}
//			defer func() { assert.Nil(cl.CloseSend()) }()
//			logrus.Debugf("stream service started: %s")
//			safe.GoRoutine(ctx, func() {
//				for {
//					rs, err := cl.Recv()
//					if err == io.EOF {
//						continue
//					}
//					if err != nil {
//						cerr <- fmt.Errorf("stream recv: %v", err)
//					}
//					lock.RLock()
//					if or, ok := inprogress[rs.Id]; ok {
//						or.Response <- rs
//					}
//					lock.RUnlock()
//				}
//			})
//			safe.GoRoutine(ctx, func() {
//				for {
//					rq := <-RequestChannel
//					logrus.Debugf("stream request: %s", rq.BidRequest.Id)
//
//					lock.Lock()
//					inprogress[rq.BidRequest.Id] = rq
//					lock.Unlock()
//					go func() {
//						<-rq.Context.Done()
//						lock.Lock()
//						delete(inprogress, rq.BidRequest.Id)
//						lock.Unlock()
//					}()
//
//					err := cl.Send(rq.BidRequest)
//					if err != nil {
//						cerr <- fmt.Errorf("stream send: %v", err)
//					}
//				}
//			})
//
//			select {
//			case err := <-cerr:
//				logrus.Debug(err)
//				return err
//			case <-ctx.Done():
//				logrus.Info("Done")
//				return nil
//			}
//		}, time.Second)
//	}
//
//}

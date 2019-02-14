package client

import (
	"container/ring"
	"context"
	"net/http"
	"sync"
	"time"

	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
)

type initRestClient struct {
}

var (
	httpClient *http.Client
	maxIdle    = config.RegisterInt("crane.supplier.max_idle_connection", 30, "maximum idle concurrentConnections count")
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
	concurrentConnections = config.RegisterInt("crane.supplier.stream.concurrentConnections", 1, "")
	insecureSever         = config.RegisterString("crane.supplier.stream.address", "127.0.0.1:9801", "")
	token                 = config.RegisterString("crane.supplier.demand.token", "forbidden", "")
	timeout               = config.RegisterDuration("crane.supplier.timeout", time.Millisecond*150, "maximum timeout")
	// RequestChannel for stream
	RequestChannel = make(chan *StreamRequest)
	lock           = sync.RWMutex{}
	connections    *ring.Ring
)

// StreamRequest for stream
type StreamRequest struct {
	BidRequest *openrtb.BidRequest
	Context    context.Context
	Response   chan<- *openrtb.BidResponse
}

type initClient struct {
	server string
	cert   string
}

var inprogress = make(map[string]*StreamRequest)

type connection struct {
	connection *grpc.ClientConn
	client     openrtb.OrtbServiceClient
	sync.RWMutex
}

func (c connection) Get() openrtb.OrtbServiceClient {
	c.RLock()
	defer c.RUnlock()
	return c.client
}

func newConnection(ctx context.Context, server, cert string) (*connection, error) {
	c := &connection{
		connection: nil,
		client:     nil,
		RWMutex:    sync.RWMutex{},
	}
	var conn *grpc.ClientConn
	var err error
	var cread credentials.TransportCredentials
	if cert == "" {
		conn, err = grpc.Dial(server, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))

		if err != nil {
			return nil, err
		}
	} else {
		cread, err = credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			return nil, err
		}
		conn, err = grpc.Dial(server, grpc.WithTransportCredentials(cread), grpc.WithBalancerName(roundrobin.Name))
		if err != nil {
			return nil, err
		}
	}
	client := openrtb.NewOrtbServiceClient(conn)

	c.connection = conn
	c.client = client

	safe.GoRoutine(ctx, func() {
		for {
			if !conn.WaitForStateChange(ctx, conn.GetState()) {
				_ = conn.Close()
				return
			}
			c.Lock()
			st := conn.GetState().String()
			switch {
			case st == "TRANSIENT_FAILURE" || st == "SHUTDOWN" || st == "Invalid-State":
				logrus.Debugf("connection closed: %s", st)
				_ = conn.Close()
				safe.Try(func() error {
					if cert == "" {
						conn, err = grpc.Dial(server, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
						if err != nil {
							logrus.Debugf("grpc dial error: %s", err)
							return nil
						}
					} else {
						conn, err = grpc.Dial(server, grpc.WithTransportCredentials(cread), grpc.WithBalancerName(roundrobin.Name))
						if err != nil {
							logrus.Debugf("grpc dial error: %s", err)
						}
					}
					client := openrtb.NewOrtbServiceClient(conn)
					c.connection = conn
					c.client = client
					return nil
				}, time.Millisecond)
			case st == "IDLE" || st == "READY" || st == "CONNECTING":
				logrus.Debugf("connection closed: %s", st)
			}
			c.Unlock()
		}
	})
	return c, nil
}

func (ic *initClient) Initialize(ctx context.Context) {
	go unaryInit(ctx)
	if demand.String() != "managed" {
		return
	}
	connections = ring.New(concurrentConnections.Int())
	for i := 0; i < concurrentConnections.Int(); i++ {
		safe.Try(func() error {
			conn, err := newConnection(ctx, ic.server, ic.cert)
			if err != nil {
				logrus.Debug(err)
				return err
			}
			connections.Next().Value = conn
			return nil
		}, time.Second)
	}

	safe.GoRoutine(ctx, func() {
		for {
			rq := <-RequestChannel
			go func(s StreamRequest) {
				var conn *connection
				var ok bool
				for {
					conn, ok = connections.Next().Value.(*connection)
					if ok {
						break
					}
				}
				s.BidRequest.Token = token.String()
				res, err := conn.Get().Ortb(rq.Context, rq.BidRequest)
				if err != nil {
					s.Response <- nil
				}
				s.Response <- res
			}(*rq)

		}
	})

}

func init() {
	initializer.Register(&initClient{
		server: insecureSever.String(),
	}, 100)
}

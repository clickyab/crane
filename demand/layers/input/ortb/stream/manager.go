package stream

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"clickyab.com/crane/demand/layers/input/ortb"
	"clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	ocontext "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	cert         = config.RegisterString("crane.demand.stream.cert", "/tls/stream/tls.crt", "")
	key          = config.RegisterString("crane.demand.stream.key", "/tls/stream/tls.key", "")
	securePort   = config.RegisterInt("crane.demand.stream.port.secure", 9800, "port for stream")
	insecurePort = config.RegisterInt("crane.demand.stream.port.insecure", 9801, "port for stream")
)

type initClient struct {
}

func (*initClient) Initialize(ctx context.Context) {

	go safe.Try(func() error {
		logrus.Info("initiate secure stream server")

		ce := make(chan error)

		creads, err := credentials.NewServerTLSFromFile(cert.String(), key.String())
		if err != nil {
			return err
		}

		defer creads.Clone()
		lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", securePort.Int()))
		if err != nil {
			return err
		}
		defer func() { assert.Nil(lis.Close()) }()
		s := grpc.NewServer(grpc.Creds(creads))
		defer s.Stop()
		openrtb.RegisterOrtbServiceServer(s, &server{})
		logrus.Info("start secure stream server")
		go func() {
			ce <- fmt.Errorf("secure stream service is dead: %v", s.Serve(lis))

		}()
		select {
		case err := <-ce:
			logrus.Error(err)
			return err
		case <-ctx.Done():
			logrus.Println("secure stream server shutdown")
			return nil
		}

	}, time.Second)

	go safe.Try(func() error {
		logrus.Info("initiate insecure stream server")

		ce := make(chan error)

		lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", insecurePort.Int()))
		if err != nil {
			return err
		}
		defer func() { assert.Nil(lis.Close()) }()
		s := grpc.NewServer()
		defer s.Stop()
		openrtb.RegisterOrtbServiceServer(s, &server{})
		logrus.Info("start insecure stream server")
		go func() {
			ce <- fmt.Errorf("insecure stream service is dead: %v", s.Serve(lis))

		}()
		select {
		case err := <-ce:
			logrus.Error(err)
			return err
		case <-ctx.Done():
			logrus.Println("secure stream server shutdown")
			return nil
		}

	}, time.Second)

}
func init() {
	initializer.Register(&initClient{}, 100)
}

type server struct{}

func (*server) OrtbStream(b openrtb.OrtbService_OrtbStreamServer) error {
	for {
		r, err := b.Recv()
		if err == io.EOF {
			logrus.Debug("recv: EOF")
			continue
		}
		if err != nil {
			logrus.Debugf("recv err: %v", err)
			return err
		}

		res, err := ortb.GrpcHandler(b.Context(), r)
		logrus.Warn(res, err)
		if err != nil {
			logrus.Debugf("handller: %s", r.Id)
			res = &openrtb.BidResponse{
				Id: r.GetId(),
			}
		}
		err = b.Send(res)
		if err != nil {
			logrus.Warnf("stream send: ", err)
			return err
		}
	}
}

func (*server) Ortb(ctx ocontext.Context, req *openrtb.BidRequest) (*openrtb.BidResponse, error) {
	return ortb.GrpcHandler(ctx, req)
}

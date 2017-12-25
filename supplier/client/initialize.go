package client

import (
	"context"
	"net/http"

	"time"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
)

type initClient struct {
}

var (
	httpClient *http.Client
	maxIdle    = config.RegisterInt("crane.supplier.max_idle_connection", 30, "maximum idle connection count")
	timeout    = config.RegisterDuration("crane.supplier.timeout", time.Second, "maximum timeout")
)

func (*initClient) Initialize(context.Context) {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdle.Int(),
			MaxIdleConns:        maxIdle.Int() + 1,
		},
	}
}

func init() {
	initializer.Register(&initClient{}, 100)
}

package rabbitmq

import (
	"context"
	"testing"

	"gopkg.in/fzerorubigd/onion.v3"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInitialize(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()

	Convey("Initialize should init some connections and some publishers", t, func() {
		ctx, cnl := context.WithCancel(context.Background())

		o := onion.New()
		connection = o.RegisterInt("test.connections.count", 5)
		publisher = o.RegisterInt("test.publishers.count", 30)

		rb := initRabbit{}
		rb.amqp = &FakeAmqp{}
		rb.Initialize(ctx)

		So(
			rb.Statistics(),
			ShouldResemble,
			map[string]interface{}{
				"connections":                5,
				"publishers":                 30,
				"publishers_per_connections": map[int]int{0: 6, 1: 6, 2: 6, 3: 6, 4: 6},
			},
		)

		defer cnl()
	})
}

package selector

import (
	"services/broker"
	"services/broker/kafka"
	"services/broker/mock"
	"services/config"
	"services/initializer"
	"services/safe"

	"github.com/Sirupsen/logrus"
)

type cfg struct {
}

func (cfg) Initialize() config.DescriptiveLayer {
	layer := config.NewDescriptiveLayer()
	layer.Add("application is in test mode and broker is not active", "services.broker.test_mode", true)
	return layer
}

func (cfg) Loaded() {
	test := config.GetBool("services.broker.test_mode")
	devel := config.GetBool("core.devel_mode")

	if devel && test {
		p := mock.GetChannelBroker()
		broker.SetActiveBroker(p)
		safe.GoRoutine(
			func() {
				ch := mock.GetChannel(10)
				for j := range ch {
					data, err := j.Encode()
					logrus.WithField("key", j.Key()).
						WithField("topic", j.Topic()).
						WithField("encode_err", err).
						Debug(string(data))
				}
			},
		)
	} else {
		// this is the only place allowed to call new cluster
		// TODO : find a better way
		c := kafka.NewCluster()
		initializer.Register(c.(initializer.Interface), 0)
		broker.SetActiveBroker(c.(broker.Interface))
	}
}

func init() {
	config.Register(&cfg{})
}

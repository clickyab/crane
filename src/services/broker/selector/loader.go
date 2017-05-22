package selector

import (
	"services/broker"
	"services/broker/mock"
	"services/config"
	"services/safe"

	"services/broker/rabbitmq"

	"github.com/Sirupsen/logrus"
)

type cfg struct {
}

func (cfg) Initialize() config.DescriptiveLayer {
	layer := config.NewDescriptiveLayer()
	layer.Add("application is in test mode and broker is not active", "services.broker.provider", "mock")
	return layer
}

func (cfg) Loaded() {
	provider := config.GetString("services.broker.provider")
	devel := config.GetBool("core.devel_mode")

	switch provider {
	case "mock":
		if !devel {
			logrus.Panic("mock is not allowed when devel is not set")
		}
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
	case "rabbitmq":
		p := rabbitmq.NewRabbitBroker()
		broker.SetActiveBroker(p)
	default:
		logrus.Panicf("there is no valid broker configured , %s is not valid", provider)
	}
}

func init() {
	config.Register(&cfg{})
}

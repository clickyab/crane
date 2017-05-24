package materialize

import (
	"clickyab.com/exchange/octopus/exchange"
	"clickyab.com/exchange/octopus/exchange/materialize/jsonbackend"
	"clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/config"

	"github.com/Sirupsen/logrus"
)

const (
	jsonDriver  string = "json"
	emptyDriver        = "empty"
)

var (
	driver = config.RegisterString("octupos.exchange.materialize.driver", jsonDriver, "driver for materialize")
)

// DemandJob returns a job for demand
// TODO : add a duration to this. for better view this is important
func DemandJob(imp exchange.Impression, dmn exchange.Demand, ads map[string]exchange.Advertise) broker.Job {
	switch *driver {
	case jsonDriver:
		return jsonbackend.DemandJob(imp, dmn, ads)
	case emptyDriver:
		return job{
			data:  []byte("demand job"),
			key:   "key",
			topic: "demand_job",
		}
	default:
		logrus.Panic("invalid driver")
		return nil
	}
}

// WinnerJob return a broker job
func WinnerJob(imp exchange.Impression, ad exchange.Advertise, slotID string) broker.Job {
	switch *driver {
	case jsonDriver:
		return jsonbackend.WinnerJob(imp, ad, slotID)
	case emptyDriver:
		return job{
			data:  []byte("winner job"),
			key:   "key",
			topic: "winner_job",
		}
	default:
		logrus.Panic("invalid driver")
		return nil
	}
}

// ShowJob return a broker job
func ShowJob(trackID, demand, slotID, adID string, IP string, winner int64, t string, supplier string, publisher string) broker.Job {
	switch *driver {
	case jsonDriver:
		return jsonbackend.ShowJob(trackID, demand, slotID, adID, IP, winner, t, supplier, publisher)
	case emptyDriver:
		return job{
			data:  []byte("show job"),
			key:   "key",
			topic: "show_job",
		}
	default:
		logrus.Panic("invalid driver")
		return nil
	}
}

// ImpressionJob return a broker job
func ImpressionJob(imp exchange.Impression) broker.Job {
	switch *driver {
	case jsonDriver:
		return jsonbackend.ImpressionJob(imp)
	case emptyDriver:
		return job{
			data:  []byte("impression job"),
			key:   "key",
			topic: "impression_job",
		}
	default:
		logrus.Panic("invalid driver")
		return nil
	}
}

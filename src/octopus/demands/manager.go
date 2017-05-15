package demands

import (
	"octopus/core"
	"octopus/demands/internal/models"
	"octopus/demands/internal/restful"
	"os"
	"os/signal"
	"services/mysql"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
)

type demandManager struct {
	activeDemands []models.Demand
	lock          *sync.RWMutex
}

func (dm *demandManager) loadDemands() {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	dm.activeDemands = models.NewManager().ActiveDemands()
	core.ResetProviders()
	for _, demand := range dm.activeDemands {
		switch demand.Type {
		case models.DemandTypeRest:
			core.Register(restful.NewRestfulClient(demand, getRawImpression), demand.GetTimeout())
		default:
			logrus.Panicf("Not supported demand type : %s", demand.Type)
		}
	}
}

func (dm *demandManager) Initialize() {
	dm.loadDemands()
	reloadChan := make(chan os.Signal)
	signal.Notify(reloadChan, syscall.SIGHUP)
	go func() {
		for i := range reloadChan {
			logrus.Infof("Reloding demands config, due to signal %s", i)
			dm.loadDemands()
		}
	}()
}

func init() {
	mysql.Register(&demandManager{lock: &sync.RWMutex{}})
}

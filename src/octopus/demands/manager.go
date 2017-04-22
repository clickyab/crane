package demands

import (
	"context"
	"net/http"
	"octopus/core"
	"octopus/demands/internal/models"
	"octopus/demands/internal/restful"
	"os"
	"os/signal"
	"services/assert"
	"services/mysql"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/xhandler"
	"github.com/fzerorubigd/xmux"
)

type demandManager struct {
	activeDemands []models.Demand
	server        *http.Server
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
			core.Register(restful.NewRestfulClient(demand, getRawImpresssion), demand.GetTimeout())
		default:
			logrus.Panicf("Not supported demand type : %s", demand.Type)
		}

	}

	if dm.server != nil {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		assert.Nil(dm.server.Shutdown(ctx))
	}

	mux := xmux.New()
	core.Mount(mux)
	// TODO : move to config
	dm.server = &http.Server{Addr: ":8080", Handler: xhandler.New(context.Background(), mux)}
	go func() {
		<-time.After(time.Second)
		if err := dm.server.ListenAndServe(); err != nil {
			logrus.Warnf("listen: %s", err)
		}
	}()
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

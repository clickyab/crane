package demands

import (
	"assert"
	"context"
	"core"
	"demands/internal/models"
	"demands/internal/restful"
	"net/http"
	"os"
	"os/signal"
	"services/mysql"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
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
	core.Reset()
	for _, demand := range dm.activeDemands {
		assert.True(demand.Type == models.DemandTypeRest, "Not supported demand type")
		core.Register(restful.NewRestfulClient(demand), demand.GetTimeout())
	}

	if dm.server != nil {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		assert.Nil(dm.server.Shutdown(ctx))
	}

	mux := xmux.New()
	core.Mount(mux)
	dm.server = &http.Server{Addr: ":8080", Handler: xhandler.New(context.Background(), mux)}
	go func() {
		if err := dm.server.ListenAndServe(); err != nil {
			logrus.Errorf("listen: %s", err)
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

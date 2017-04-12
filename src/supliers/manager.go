package supliers

import (
	"entity"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"services/assert"
	"services/mysql"
	"supliers/internal/models"
	"supliers/internal/restful"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
)

var sm *supplierManager

type supplierManager struct {
	suppliers map[string]models.Supplier
	lock      *sync.RWMutex
}

func (sm *supplierManager) loadSuppliers() {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	m := models.NewManager()
	sm.suppliers = m.GetSuppliers()
}

func (sm *supplierManager) Initialize() {
	sm.loadSuppliers()
	reloadChan := make(chan os.Signal)
	signal.Notify(reloadChan, syscall.SIGHUP)
	go func() {
		for i := range reloadChan {
			logrus.Infof("Reloding supplier config, due to signal %s", i)
			sm.loadSuppliers()
		}
	}()
}

// getSupplier return a single supplier by its id
func getSupplier(key string) (*models.Supplier, error) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	if s, ok := sm.suppliers[key]; ok {
		return &s, nil
	}

	return nil, fmt.Errorf("supplier with key %s not found", key)
}

// GetImpression try to get an impression from a http request
func GetImpression(key string, r *http.Request) (entity.Impression, error) {
	sup, err := getSupplier(key)
	if err != nil {
		return nil, err
	}

	assert.True(sup.SType == "rest", "Not a supported type", sup.SType)
	return restful.GetImpression(sup, r)
}

func init() {
	sm = &supplierManager{lock: &sync.RWMutex{}}
	mysql.Register(sm)
}

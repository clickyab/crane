package apps

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var (
	app pool.Interface
)

// GetAppOrFake try to get app. do not use it in initializer
func GetAppOrFake(sup entity.Supplier, appPackage string) (entity.Publisher, error) {
	d := &entities.App{}
	res, err := app.Get(fmt.Sprintf("%s/%s", sup.Name(), appPackage), d)
	if err != nil {
		if sup.AllowCreate() {
			return entities.NewFakePublisher(sup, appPackage, entity.PublisherTypeApp), nil
		}
		return nil, err
	}
	if d.Status != 1 {
		return nil, fmt.Errorf("blocked app")
	}
	d = res.(*entities.App)
	d.Supp = sup
	return d, nil
}

// GetApp try to get app. do not use it in initializer
func GetApp(sup entity.Supplier, appPackage string, token string) (entity.Publisher, error) {
	d := &entities.App{}
	res, err := app.Get(fmt.Sprintf("%s/%s", sup.Name(), appPackage), d)
	if err == nil {
		d = res.(*entities.App)
		d.Supp = sup
		return d, nil
	}
	app, err := entities.FindOrAddApp(sup, appPackage, token)
	if err != nil {
		return nil, err
	}
	return app, nil
}

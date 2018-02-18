package clickyabapps

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var (
	app pool.Interface
)

// GetApp try to get app. do not use it in initializer
func GetApp(sup entity.Supplier, appToken string) (entity.Publisher, error) {
	d := &entities.App{}
	res, err := app.Get(appToken, d)
	if err != nil {
		return nil, err
	}
	if d.Status != 1 {
		return nil, fmt.Errorf("blocked app")
	}
	d = res.(*entities.App)
	d.Supp = sup
	return d, nil
}

package clickyabwebsite

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var (
	websites pool.Interface
)

// GetWebSite try to get website. do not use it in initializer
func GetWebSite(sup entity.Supplier, pubID string) (entity.Publisher, error) {
	d := &entities.Website{}
	res, err := websites.Get(pubID, d)
	if err != nil {
		return nil, err
	}
	if d.Status != 1 {
		return nil, fmt.Errorf("blocked site")
	}
	d = res.(*entities.Website)
	d.Supp = sup

	return d, nil
}

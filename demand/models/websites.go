package models

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var (
	websites pool.Interface
)

// GetWebSite try to get website. do not use it in initializer
func GetWebSite(sup entity.Supplier, domain string) (entity.Publisher, error) {
	d := &entities.Website{}
	res, err := websites.Get(fmt.Sprintf("%s/%s", sup.Name(), domain), d)
	if err != nil {
		if sup.AllowCreate() {
			return entities.NewFakePublisher(sup, domain), nil
		}
		return nil, err
	}
	if d.Status != 1 {
		return nil, fmt.Errorf("blocked site")
	}
	d = res.(*entities.Website)
	d.Supp = sup

	return d, nil
}

// GetWebSiteID try to get website. do not use it in initializer
func GetWebSiteID(sup entity.Supplier, domain string, pid int64) (int64, error) {
	d := &entities.Website{}
	res, err := websites.Get(fmt.Sprintf("%s/%s", sup.Name(), domain), d)
	if err == nil {
		d = res.(*entities.Website)
		d.Supp = sup
		return d.WID, nil
	}
	return entities.FindOrAddWebsite(sup, domain, pid)
}

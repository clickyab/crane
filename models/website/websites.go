package website

import (
	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var (
	websites pool.Interface
)

// GetWebSiteOrFake try to get website. do not use it in initializer
func GetWebSiteOrFake(sup entity.Supplier, domain string) (entity.Publisher, error) {
	d := &entities.Website{}
	res, err := websites.Get(fmt.Sprintf("%s/%s", sup.Name(), domain), d)
	if err != nil {
		if sup.AllowCreate() {
			return entities.NewFakePublisher(sup, domain, entity.PublisherTypeWeb), nil
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

// GetWebSite try to get website. do not use it in initializer
func GetWebSite(sup entity.Supplier, domain string, pid int64) (entity.Publisher, error) {
	d := &entities.Website{}
	res, err := websites.Get(fmt.Sprintf("%s/%s", sup.Name(), domain), d)
	if err == nil {
		d = res.(*entities.Website)
		d.Supp = sup
		return d, nil
	}
	return entities.FindOrAddWebsite(sup, domain, pid)
}

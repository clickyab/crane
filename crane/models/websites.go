package models

import (
	"fmt"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var (
	websites pool.Interface
)

// GetWebSite try to get website. do not use it in initializer
func GetWebSite(sup, domain string) (entity.Publisher, error) {
	d := &entities.Website{}
	res, err := websites.Get(fmt.Sprintf("%s/%s", sup, domain), d)
	if err != nil {
		return nil, err
	}

	return res.(entity.Publisher), nil
}

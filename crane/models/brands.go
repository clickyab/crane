package models

import (
	"clickyab.com/crane/crane/models"
	"github.com/clickyab/services/pool"
	"clickyab.com/crane/crane/models/internal/entities"
)

var brands pool.Interface

// GetBrandID return the brand object based on the brand name
func GetBrandID(b string) (int64, error) {
	bs := &entities.Brand{}
	data, err := brands.Get(b, bs)
	if err != nil {
		return 0, err
	}
	return data.(*entities.Brand).ID, nil
}

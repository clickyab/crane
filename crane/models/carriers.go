package models

import (
	"clickyab.com/crane/crane/models/internal/entities"
	"github.com/clickyab/services/pool"
)

var carriers pool.Interface

// GetCarriers return the carrier object based on the carrier name
func GetCarrierID(b string) (int64, error) {
	bs := &entities.Carrier{}
	data, err := carriers.Get(b, bs)
	if err != nil {
		return 0, err
	}
	return data.(*entities.Carrier).ID, nil
}
